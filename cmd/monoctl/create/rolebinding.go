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

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoctl/cmd/monoctl/flags"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/usecases"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/util"
	auth_util "gitlab.figo.systems/platform/monoskope/monoctl/internal/util/auth"
)

func NewCreateRoleBindingCmd() *cobra.Command {
	var role string
	var scope string
	var resource string

	createRoleBindingCmd := &cobra.Command{
		Use:   "rolebinding <EMAIL>",
		Short: "Create rolebinding.",
		Long:  `Creates a rolebinding for the given user.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateRoleBindingUseCase(configManager.GetConfig(), args[0], role, scope, resource).Run(ctx)
			})
		},
	}

	flags := createRoleBindingCmd.Flags()
	flags.StringVarP(&role, "role", "r", "", "Role of the rolebinding.")
	flags.StringVarP(&scope, "scope", "s", "", "Scope of the rolebinding.")
	flags.StringVarP(&resource, "resource", "e", "", "Resource of the rolebinding.")

	util.PanicOnError(createRoleBindingCmd.MarkFlagRequired("role"))
	util.PanicOnError(createRoleBindingCmd.MarkFlagRequired("scope"))

	return createRoleBindingCmd
}
