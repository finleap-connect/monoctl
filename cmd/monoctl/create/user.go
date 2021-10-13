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

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoctl/cmd/monoctl/flags"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/usecases"
	auth_util "gitlab.figo.systems/platform/monoskope/monoctl/internal/util/auth"
)

func NewCreateUserCmd() *cobra.Command {
	var username string

	createUserCmd := &cobra.Command{
		Use:   "user <EMAIL>",
		Short: "Create user.",
		Long:  `Creates a user.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if username == "" {
				username = strings.Split(args[0], "@")[0]
			}
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateUserUseCase(configManager.GetConfig(), username, args[0]).Run(ctx)
			})
		},
	}

	flags := createUserCmd.Flags()

	flags.StringVarP(&username, "username", "u", "", "Name of the user. By default the local part of the email address is used.")

	return createUserCmd
}
