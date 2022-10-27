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
	"strings"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	"github.com/finleap-connect/monoctl/internal/util"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/spf13/cobra"
)

func NewCreateClusterCmd() *cobra.Command {
	var (
		apiServerAddress string
		caCertBundleFile string
	)

	cmd := &cobra.Command{
		Use:   "cluster <NAME>",
		Short: "Create cluster.",
		Long:  `Creates a Kubernetes cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			name := args[0]

			u, err := url.Parse(apiServerAddress)
			if err != nil {
				return err
			}
			if !u.IsAbs() || u.Hostname() == "" {
				return errors.New("invalid url format")
			}

			sanitizedName, err := k8s.GetK8sName(name)
			if err != nil {
				return err
			}
			name = sanitizedName

			caCertBundle, err := os.ReadFile(caCertBundleFile)
			if err != nil {
				return fmt.Errorf("failed to read CA certificates from '%s': %s", caCertBundleFile, err)
			}
			caCertBundle = []byte(strings.TrimSpace(string(caCertBundle)))

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewCreateClusterUseCase(configManager.GetConfig(), name, apiServerAddress, caCertBundle).Run(ctx)
			})
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&caCertBundleFile, "ca-filepath", "c", "", "Path to the file containing the CA certificate bundle of the cluster in PEM format.")
	flags.StringVarP(&caCertBundleFile, "api-server-address", "a", "", "Address of the KubeAPIServer of the cluster.")
	util.PanicOnError(cmd.MarkFlagRequired("ca-filepath"))
	util.PanicOnError(cmd.MarkFlagRequired("api-server-address"))

	return cmd
}
