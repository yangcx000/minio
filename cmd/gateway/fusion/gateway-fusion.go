/*
 * MinIO Object Storage (c) 2021 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 *
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
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS]{{end}} [ENDPOINT]
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
ENDPOINT:
  s3 server endpoint. Default ENDPOINT is https://s3.amazonaws.com

EXAMPLES:
  1. Start minio gateway server for AWS S3 backend
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_USER{{.AssignmentOperator}}accesskey
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
     {{.Prompt}} {{.HelpName}}

  2. Start minio gateway server for AWS S3 backend with edge caching enabled
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_USER{{.AssignmentOperator}}accesskey
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_DRIVES{{.AssignmentOperator}}"/mnt/drive1,/mnt/drive2,/mnt/drive3,/mnt/drive4"
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_EXCLUDE{{.AssignmentOperator}}"bucket1/*,*.png"
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_QUOTA{{.AssignmentOperator}}90
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_AFTER{{.AssignmentOperator}}3
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_LOW{{.AssignmentOperator}}75
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_HIGH{{.AssignmentOperator}}85
     {{.Prompt}} {{.HelpName}}
`
	minio.RegisterGatewayCommand(cli.Command{
		Name:               minio.FusionBackendGateway,
		Usage:              "FusionStore storage Service (Fusion)",
		Action:             fusionGatewayMain,
		CustomHelpTemplate: fusionGatewayTemplate,
		HideHelpCommand:    true,
	})
}

// Handler for 'minio gateway fusion' command line.
func fusionGatewayMain(ctx *cli.Context) {
	mgsAddr := ctx.String("mgs")
	if len(mgsAddr) == 0 {
		logger.FatalIf(errors.New("mgs addr empty"), "", nil)
	}
	// Start the gateway.
	minio.StartGateway(ctx, &Fusion{
		mgsAddr: mgsAddr,
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

// NewGatewayLayer returns s3 ObjectLayer.
func (g *Fusion) NewGatewayLayer(creds madmin.Credentials) (minio.ObjectLayer, error) {
	_ = creds
	return fusionstore.New(g.mgsAddr)
}
