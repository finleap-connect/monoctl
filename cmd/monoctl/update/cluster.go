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

package update

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

func NewUpdateClusterCmd() *cobra.Command {
	var (
		displayName      string
		apiServerAddress string
		caCertBundleFile string
	)

	cmd := &cobra.Command{
		Use:   "cluster <NAME>",
		Short: "Update cluster.",
		Long:  `Updates a cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var caCertBundle []byte

			if caCertBundleFile == "" && displayName == "" && apiServerAddress == "" && caCertBundle == nil {
				return errors.New("nothing to update")
			}

			if caCertBundleFile != "" {
				caCertBundle, err = os.ReadFile(caCertBundleFile)
				if err != nil {
					return fmt.Errorf("failed to read CA certificates from '%s': %s", caCertBundleFile, err)
				}
			}

			if apiServerAddress != "" {
				u, err := url.Parse(apiServerAddress)
				if err != nil {
					return err
				}
				if !u.IsAbs() || u.Hostname() == "" {
					return errors.New("invalid url format")
				}
			}

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewUpdateClusterUseCase(configManager.GetConfig(), args[0], displayName, apiServerAddress, caCertBundle).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&displayName, "display-name", "d", "", "New display name of the cluster")
	flags.StringVarP(&apiServerAddress, "api-server-address", "a", "", "New KubeAPIServer address")
	flags.StringVarP(&caCertBundleFile, "ca-cert-path", "c", "", "New CA certificate bundle file")

	return cmd
}
