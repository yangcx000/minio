/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package fusion

import (
	"errors"

	"github.com/minio/cli"
	"github.com/minio/madmin-go"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore"
	"github.com/minio/minio/internal/logger"
)

func init() {
	const fusionGatewayTemplate = `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[FLAGS]{{end}}
	{{if .VisibleFlags}}
FLAGS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
	`
	minio.RegisterGatewayCommand(cli.Command{
		Name:               minio.FusionBackendGateway,
		Usage:              "FusionStore Object Store Gateway",
		Action:             fusionGatewayMain,
		CustomHelpTemplate: fusionGatewayTemplate,
		HideHelpCommand:    true,
	})
}

func fusionGatewayMain(ctx *cli.Context) {
	if len(ctx.String("address")) == 0 {
		logger.FatalIf(errors.New("params error"), "fusion gateway address empty")
	}
	if len(ctx.String("mgs")) == 0 {
		logger.FatalIf(errors.New("params error"), "mgs addr empty")
	}
	// Start the gateway.
	minio.StartGateway(ctx, &Fusion{
		mgsAddr: ctx.String("mgs"),
	})
}

// Fusion implements Gateway.
type Fusion struct {
	mgsAddr string
}

// Name implements Gateway interface.
func (g *Fusion) Name() string {
	return minio.FusionBackendGateway
}

// NewGatewayLayer returns fusionstore ObjectLayer.
func (g *Fusion) NewGatewayLayer(creds madmin.Credentials) (minio.ObjectLayer, error) {
	// env cred no use
	_ = creds
	return fusionstore.New(g.mgsAddr)
}
