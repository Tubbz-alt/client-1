package avatars

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/lru"
	"github.com/keybase/client/go/protocol/keybase1"
)

type avatarLoadPair struct {
	username string
	format   keybase1.AvatarFormat
	path     string
}

type avatarLoadSpec struct {
	hits   []avatarLoadPair
	misses []avatarLoadPair
	stales []avatarLoadPair
}

func (a avatarLoadSpec) missUsernames() (res []string) {
	for _, u := range a.misses {
		res = append(res, u.username)
	}
	return res
}

func (a avatarLoadSpec) missFormats() (res []keybase1.AvatarFormat) {
	for _, f := range a.misses {
		res = append(res, f.format)
	}
	return res
}

type populateArg struct {
	username string
	format   keybase1.AvatarFormat
	url      string
}

type CachingSource struct {
	libkb.Contextified

	diskLRU        *lru.DiskLRU
	staleThreshold time.Duration
	simpleSource   *SimpleSource

	populateCacheCh chan populateArg
}

func NewCachingSource(g *libkb.GlobalContext, staleThreshold time.Duration, size int) *CachingSource {
	c := &CachingSource{
		Contextified:    libkb.NewContextified(g),
		diskLRU:         lru.NewDiskLRU("avatars", 1, size),
		staleThreshold:  staleThreshold,
		simpleSource:    NewSimpleSource(g),
		populateCacheCh: make(chan populateArg, 100),
	}
	for i := 0; i < 10; i++ {
		go c.populateCacheWorker()
	}
	return c
}

func (c *CachingSource) debug(ctx context.Context, msg string, args ...interface{}) {
	c.G().Log.CDebugf(ctx, "Avatars.CachingSource: %s", fmt.Sprintf(msg, args...))
}

func (c *CachingSource) avatarKey(username string, format keybase1.AvatarFormat) string {
	return fmt.Sprintf("%s:%s", username, format.String())
}

func (c *CachingSource) isStale(item lru.DiskLRUEntry) bool {
	return c.G().GetClock().Now().Sub(item.Ctime) > c.staleThreshold
}

func (c *CachingSource) specLoad(ctx context.Context, usernames []string, formats []keybase1.AvatarFormat) (res avatarLoadSpec, err error) {
	for _, username := range usernames {
		for _, format := range formats {
			found, entry, err := c.diskLRU.Get(ctx, c.G(), c.avatarKey(username, format))
			if err != nil {
				return res, err
			}
			lp := avatarLoadPair{
				username: username,
				format:   format,
			}
			if found {
				lp.path = entry.Value.(string)
				if c.isStale(entry) {
					res.stales = append(res.stales, lp)
				} else {
					res.hits = append(res.hits, lp)
				}
			} else {
				res.misses = append(res.misses, lp)
			}
		}
	}
	return res, nil
}

func (c *CachingSource) randomFileName() (string, error) {
	return libkb.RandHexString("avatar", 16)
}

func (c *CachingSource) commitAvatarToDisk(ctx context.Context, data io.ReadCloser) (path string, err error) {
	fileName, err := c.randomFileName()
	if err != nil {
		return path, err
	}
	path = filepath.Join(c.G().GetCacheDir(), fileName)
	file, err := os.Create(path)
	if err != nil {
		return path, err
	}
	_, err = io.Copy(file, data)
	if err != nil {
		return path, err
	}
	file.Close()
	return path, nil
}

func (c *CachingSource) populateCacheWorker() {
	for arg := range c.populateCacheCh {
		ctx := context.Background()
		resp, err := http.Get(arg.url)
		if err != nil {
			c.debug(ctx, "populateCacheWorker: failed to download avatar: %s", err)
			continue
		}
		path, err := c.commitAvatarToDisk(ctx, resp.Body)
		if err != nil {
			c.debug(ctx, "populateCacheWorker: failed to write to disk: %s", err)
			continue
		}
		evict, err := c.diskLRU.Put(ctx, c.G(), c.avatarKey(arg.username, arg.format), path)
		if err != nil {
			c.debug(ctx, "populateCacheWorker: failed to put into LRU: %s", err)
			continue
		}
		if evict != nil {
			// Remove the file
			if err := os.Remove(evict.Value.(string)); err != nil {
				c.debug(ctx, "populateCacheWorker: failed to remove evicted file: %s", err)
			}
		}
	}
}

func (c *CachingSource) makeURL(path string) keybase1.AvatarUrl {
	return keybase1.AvatarUrl("file://" + path)
}

func (c *CachingSource) mergeRes(res *keybase1.LoadAvatarsRes, m keybase1.LoadAvatarsRes) {
	for username, rec := range m.Picmap {
		for format, url := range rec {
			res.Picmap[username][format] = url
		}
	}
}

func (c *CachingSource) LoadUsers(ctx context.Context, usernames []string, formats []keybase1.AvatarFormat) (res keybase1.LoadAvatarsRes, err error) {
	loadSpec, err := c.specLoad(ctx, usernames, formats)
	if err != nil {
		return res, err
	}

	// Fill in the hits
	c.simpleSource.allocRes(&res, usernames)
	for _, hit := range loadSpec.hits {
		res.Picmap[hit.username][hit.format] = c.makeURL(hit.path)
	}
	// File in stale
	for _, stale := range loadSpec.stales {
		res.Picmap[stale.username][stale.format] = c.makeURL(stale.path)
	}

	// Go get the misses
	missRes, err := c.simpleSource.LoadUsers(ctx, loadSpec.missUsernames(), loadSpec.missFormats())
	if err != nil {
		return res, err
	}
	c.mergeRes(&res, missRes)
	return res, nil
}
