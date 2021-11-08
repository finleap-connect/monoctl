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

package get

import (
	"context"
	"errors"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

func NewGetClusterCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster-credentials [CLUSTER] [ROLE]",
		Short: "Get cluster credentials.",
		Long:  `Get credentials for a specific cluster known to the m8 control plane.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterName := args[0]
			if clusterName == "" {
				return errors.New("cluster must be specified")
			}

			clusterRole := args[1]
			if clusterRole == "" {
				return errors.New("role must be specified")
			}

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetClusterCredentialsUseCase(configManager, clusterName, clusterRole).Run(ctx)
			})
		},
	}
	return cmd
}
