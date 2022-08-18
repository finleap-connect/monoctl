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

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	"github.com/finleap-connect/monoctl/internal/util"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

func NewCreateRoleBindingCmd() *cobra.Command {
	var role string
	var scope string
	var resource string

	cmd := &cobra.Command{
		Use:   "rolebinding <EMAIL> [<EMAIL>,...]",
		Short: "Create rolebinding.",
		Long:  `Creates a rolebinding for the given user.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateRoleBindingUseCase(configManager.GetConfig(), args, role, scope, resource).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&role, "role", "r", "", "Role of the rolebinding.")
	flags.StringVarP(&scope, "scope", "s", "", "Scope of the rolebinding.")
	flags.StringVarP(&resource, "resource", "e", "", "Resource of the rolebinding.")

	util.PanicOnError(cmd.MarkFlagRequired("role"))
	util.PanicOnError(cmd.MarkFlagRequired("scope"))

	return cmd
}
