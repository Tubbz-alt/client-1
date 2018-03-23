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

	StartBackgroundTasks()
	StopBackgroundTasks()
}

func CreateSourceFromEnv(g *libkb.GlobalContext) Source {
	c := NewCachingSource(g, 6*time.Hour, 1000)
	c.StartBackgroundTasks()
	return c
}
