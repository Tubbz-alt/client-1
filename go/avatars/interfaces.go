package avatars

import (
	"context"
	"time"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
)

type Source interface {
	LoadUsers(context.Context, []string, []keybase1.AvatarFormat) (keybase1.LoadAvatarsRes, error)
	LoadTeams(context.Context, []string, []keybase1.AvatarFormat) (keybase1.LoadAvatarsRes, error)
}

func CreateSourceFromEnv(g *libkb.GlobalContext) Source {
	s, err := NewCachingSource(g, 6*time.Hour, 1000)
	if err != nil {
		return NewSimpleSource(g)
	}
	return s
}
