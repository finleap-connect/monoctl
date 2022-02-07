// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package create

import (
	"context"
	"strings"
	"time"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	"github.com/finleap-connect/monoctl/internal/util"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	apiGateway "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/spf13/cobra"
)

func NewCreateAPITokenCmd() *cobra.Command {
	var userId string
	var scopes []string
	validity := time.Hour * 24

	cmd := &cobra.Command{
		Use:   "api-token",
		Short: "Let m8 issue an API token.",
		Long:  `Retrieve an API token issued by the m8 control plane.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateAPITokenUsecase(configManager, userId, scopes, validity).Run(ctx)
			})
		},
	}
	flags := cmd.Flags()

	flags.StringVarP(&userId, "user", "u", "", "Specify the name or UUID of the user for whom the token should be issued. If not a UUID it will be treated as username.")
	util.PanicOnError(cmd.MarkFlagRequired("user"))

	avialbleScopes := make([]string, len(apiGateway.AuthorizationScope_name))
	for _, value := range apiGateway.AuthorizationScope_name {
		avialbleScopes[apiGateway.AuthorizationScope_value[value]] = value
	}
	scopesUsage := "Specify the scopes for which the token should be valid.\nAvailable scopes: " + strings.Join(avialbleScopes, ", ")
	flags.StringSliceVarP(&scopes, "scopes", "s", scopes, scopesUsage)
	util.PanicOnError(cmd.MarkFlagRequired("scopes"))

	flags.DurationVarP(&validity, "validity", "v", validity, "Specify the validity period of the token.")

	return cmd
}
