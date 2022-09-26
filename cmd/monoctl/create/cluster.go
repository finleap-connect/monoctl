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
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/spf13/cobra"
)

func NewCreateClusterCmd() *cobra.Command {
	var (
		name        string
		displayName string
	)

	cmd := &cobra.Command{
		Use:   "cluster <KUBE_API_SERVER_ADDRESS> <CA_CERT_FILE>",
		Short: "Create cluster.",
		Long:  `Creates a cluster. The name and display name are derived from the KubeAPIServer address given. They can be overridden by flags.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiServerAddress := args[0]
			caCertBundleFile := args[1]

			u, err := url.Parse(apiServerAddress)
			if err != nil {
				return err
			}
			if !u.IsAbs() || u.Hostname() == "" {
				return errors.New("invalid url format")
			}

			if name == "" {
				name = u.Hostname()
			}

			sanitizedName, err := k8s.GetK8sName(name)
			if err != nil {
				return err
			}
			name = sanitizedName

			if displayName == "" {
				displayName = name
			}

			caCertBundle, err := os.ReadFile(caCertBundleFile)
			if err != nil {
				return fmt.Errorf("failed to read CA certificates from '%s': %s", caCertBundleFile, err)
			}

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateClusterUseCase(configManager.GetConfig(), name, displayName, apiServerAddress, caCertBundle).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&name, "name", "n", "", "Name of the cluster")
	flags.StringVarP(&displayName, "display-name", "d", "", "Display name of the cluster")

	return cmd
}
