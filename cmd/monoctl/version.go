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

package main

import (
	"context"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/finleap-connect/monoctl/internal/version"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	var clientOnly bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  `Prints version information and the commit of the local client and of the server if authenticated`,
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintVersion()
			if clientOnly {
				return nil
			}

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetServerVersionUseCase(configManager.GetConfig()).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&clientOnly, "client", false, "Client version only (no server required).")
	return cmd
}
