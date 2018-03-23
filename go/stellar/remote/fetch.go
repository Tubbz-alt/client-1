package remote

import (
	"context"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
	"github.com/keybase/client/go/stellar/bundle"
)

type fetchRes struct {
	Status       libkb.AppStatus `json:"status"`
	EncryptedB64 string          `json:"encrypted"`
	VisibleB64   string          `json:"visible"`
}

func (r *fetchRes) GetAppStatus() *libkb.AppStatus {
	return &r.Status
}

// Fetch and unbox the latest bundle from the server.
func Fetch(ctx context.Context, g *libkb.GlobalContext) (res keybase1.StellarBundle, err error) {
	defer g.CTraceTimed(ctx, "Stellar.Fetch", func() error { return err })()
	arg := libkb.NewAPIArgWithNetContext(ctx, "stellar/bundle")
	arg.SessionType = libkb.APISessionTypeREQUIRED
	var apiRes fetchRes
	err = g.API.GetDecode(arg, &apiRes)
	if err != nil {
		return res, err
	}
	decodeRes, err := bundle.Decode(apiRes.EncryptedB64)
	if err != nil {
		return res, err
	}
	pukring, err := g.GetPerUserKeyring()
	if err != nil {
		return res, err
	}
	err = pukring.Sync(ctx)
	if err != nil {
		return res, err
	}
	puk, err := pukring.GetSeedByGeneration(ctx, decodeRes.Enc.Gen)
	if err != nil {
		return res, err
	}
	res, _, err = bundle.Unbox(decodeRes, apiRes.VisibleB64, puk)
	return res, err
}
