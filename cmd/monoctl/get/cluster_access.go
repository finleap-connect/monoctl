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

func NewGetClusterAccess() *cobra.Command {
	var tenantName string
	var clusterName string

	cmd := &cobra.Command{
		Use:     "cluster-access",
		Aliases: []string{"cluster"},
		Short:   "Get cluster-access.",
		Long:    `Get cluster-access.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tenantName) > 0 && len(clusterName) > 0 {
				return errors.New("only tenant OR cluster has can be specified")
			} else if len(tenantName) == 0 && len(clusterName) == 0 {
				return errors.New("neither tenant nor cluster has been specified")
			}

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetClusterAccessUseCase(configManager.GetConfig(), getOutputOptions(), tenantName, clusterName).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tenantName, "tenant-name", "t", "", "Specify to see clusters the given tenant has access to.")
	flags.StringVarP(&clusterName, "cluster-name", "c", "", "Specify to see tenants which have access to the given cluster.")

	return cmd
}
