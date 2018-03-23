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
	name   string
	format keybase1.AvatarFormat
	path   string
}

type avatarLoadSpec struct {
	hits   []avatarLoadPair
	misses []avatarLoadPair
	stales []avatarLoadPair
}

func (a avatarLoadSpec) details(l []avatarLoadPair) (names []string, formats []keybase1.AvatarFormat) {
	fmap := make(map[keybase1.AvatarFormat]bool)
	umap := make(map[string]bool)
	for _, m := range l {
		umap[m.name] = true
		fmap[m.format] = true
	}
	for u := range umap {
		names = append(names, u)
	}
	for f := range fmap {
		formats = append(formats, f)
	}
	return names, formats
}

func (a avatarLoadSpec) missDetails() ([]string, []keybase1.AvatarFormat) {
	return a.details(a.misses)
}

func (a avatarLoadSpec) staleDetails() ([]string, []keybase1.AvatarFormat) {
	return a.details(a.stales)
}

type populateArg struct {
	name   string
	format keybase1.AvatarFormat
	url    keybase1.AvatarUrl
}

type CachingSource struct {
	libkb.Contextified

	diskLRU        *lru.DiskLRU
	staleThreshold time.Duration
	simpleSource   *SimpleSource
	httpSrv        *libkb.RandomPortHTTPSrv
	simpleMode     bool

	populateCacheCh chan populateArg
}

var _ Source = (*CachingSource)(nil)

func NewCachingSource(g *libkb.GlobalContext, staleThreshold time.Duration, size int) *CachingSource {
	return &CachingSource{
		Contextified:   libkb.NewContextified(g),
		diskLRU:        lru.NewDiskLRU("avatars", 1, size),
		staleThreshold: staleThreshold,
		simpleSource:   NewSimpleSource(g),
		httpSrv:        libkb.NewRandomPortHTTPSrv(),
	}
}

func (c *CachingSource) StartBackgroundTasks() {
	c.populateCacheCh = make(chan populateArg, 100)
	for i := 0; i < 10; i++ {
		go c.populateCacheWorker()
	}
	if err := c.httpSrv.Start(); err != nil {
		c.debug(context.Background(), "failed to start local http server, defaulting to simple mode")
		c.simpleMode = true
	} else {
		c.httpSrv.HandleFunc("/a", c.serveHTTPAvatar)
	}
}

func (c *CachingSource) StopBackgroundTasks() {
	close(c.populateCacheCh)
	c.simpleMode = false
	c.httpSrv.Stop()
}

func (c *CachingSource) debug(ctx context.Context, msg string, args ...interface{}) {
	c.G().Log.CDebugf(ctx, "Avatars.CachingSource: %s", fmt.Sprintf(msg, args...))
}

func (c *CachingSource) avatarKey(name string, format keybase1.AvatarFormat) string {
	return fmt.Sprintf("%s:%s", name, format.String())
}

func (c *CachingSource) isStale(item lru.DiskLRUEntry) bool {
	return c.G().GetClock().Now().Sub(item.Ctime) > c.staleThreshold
}

func (c *CachingSource) specLoad(ctx context.Context, names []string, formats []keybase1.AvatarFormat) (res avatarLoadSpec, err error) {
	for _, name := range names {
		for _, format := range formats {
			found, entry, err := c.diskLRU.Get(ctx, c.G(), c.avatarKey(name, format))
			if err != nil {
				return res, err
			}
			lp := avatarLoadPair{
				name:   name,
				format: format,
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
		c.debug(ctx, "populateCacheWorker: fetching: name: %s format: %s url: %s", arg.name,
			arg.format, arg.url)
		resp, err := http.Get(arg.url.String())
		if err != nil {
			c.debug(ctx, "populateCacheWorker: failed to download avatar: %s", err)
			continue
		}
		path, err := c.commitAvatarToDisk(ctx, resp.Body)
		if err != nil {
			c.debug(ctx, "populateCacheWorker: failed to write to disk: %s", err)
			continue
		}
		evict, err := c.diskLRU.Put(ctx, c.G(), c.avatarKey(arg.name, arg.format), path)
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

func (c *CachingSource) dispatchPopulateFromRes(ctx context.Context, res keybase1.LoadAvatarsRes) {
	for name, rec := range res.Picmap {
		for format, url := range rec {
			if url != "" {
				c.populateCacheCh <- populateArg{
					name:   name,
					format: format,
					url:    url,
				}
			}
		}
	}
}

func (c *CachingSource) serveHTTPAvatar(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("p")
	file, err := os.Open(path)
	if err != nil {
		c.debug(context.Background(), "serveHTTPAvatar: failed to read file: %s", err)
		return
	}
	io.Copy(w, file)
}

func (c *CachingSource) makeURL(path string) keybase1.AvatarUrl {
	addr, _ := c.httpSrv.Addr()
	return keybase1.MakeAvatarURL(fmt.Sprintf("http://%s/a?p=%s", addr, path))
}

func (c *CachingSource) mergeRes(res *keybase1.LoadAvatarsRes, m keybase1.LoadAvatarsRes) {
	for username, rec := range m.Picmap {
		for format, url := range rec {
			res.Picmap[username][format] = url
		}
	}
}

func (c *CachingSource) loadNames(ctx context.Context, names []string, formats []keybase1.AvatarFormat,
	remoteFetch func(context.Context, []string, []keybase1.AvatarFormat) (keybase1.LoadAvatarsRes, error)) (res keybase1.LoadAvatarsRes, err error) {
	loadSpec, err := c.specLoad(ctx, names, formats)
	if err != nil {
		return res, err
	}
	c.debug(ctx, "loadSpec: hits: %d stales: %d misses: %d", len(loadSpec.hits), len(loadSpec.stales),
		len(loadSpec.misses))

	// Fill in the hits
	c.simpleSource.allocRes(&res, names)
	for _, hit := range loadSpec.hits {
		res.Picmap[hit.name][hit.format] = c.makeURL(hit.path)
	}
	// Fill in stales
	for _, stale := range loadSpec.stales {
		res.Picmap[stale.name][stale.format] = c.makeURL(stale.path)
	}

	// Go get the misses
	missNames, missFormats := loadSpec.missDetails()
	if len(missNames) > 0 {
		loadRes, err := remoteFetch(ctx, missNames, missFormats)
		if err == nil {
			c.mergeRes(&res, loadRes)
			c.dispatchPopulateFromRes(ctx, loadRes)
		} else {
			c.debug(ctx, "loadSpec: failed to load server miss reqs: %s", err)
		}
	}
	// Spawn off a goroutine to reload stales
	staleNames, staleFormats := loadSpec.staleDetails()
	if len(staleNames) > 0 {
		go func() {
			c.debug(context.Background(), "loadSpec: spawning stale background load: names: %d",
				len(staleNames))
			loadRes, err := remoteFetch(context.Background(), staleNames, staleFormats)
			if err == nil {
				c.dispatchPopulateFromRes(ctx, loadRes)
			} else {
				c.debug(ctx, "loadSpec: failed to load server stale reqs: %s", err)
			}
		}()
	}
	return res, nil
}

func (c *CachingSource) LoadUsers(ctx context.Context, usernames []string, formats []keybase1.AvatarFormat) (res keybase1.LoadAvatarsRes, err error) {
	defer c.G().Trace("CachingSource.LoadUsers", func() error { return err })()
	// If we failed to startup our HTTP server, we just don't do any caching
	if c.simpleMode {
		return c.simpleSource.LoadUsers(ctx, usernames, formats)
	}
	return c.loadNames(ctx, usernames, formats, c.simpleSource.LoadUsers)
}

func (c *CachingSource) LoadTeams(ctx context.Context, teams []string, formats []keybase1.AvatarFormat) (res keybase1.LoadAvatarsRes, err error) {
	defer c.G().Trace("CachingSource.LoadTeams", func() error { return err })()
	if c.simpleMode {
		return c.simpleSource.LoadTeams(ctx, teams, formats)
	}
	return c.loadNames(ctx, teams, formats, c.simpleSource.LoadTeams)
}
