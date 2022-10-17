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
	"fmt"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

var (
	kubeConfigPath string
	overwrite      bool
)

func NewUpdateKubeconfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Updates the users kubeconfig file with endpoint information.",
		Long: `Updates the users kubeconfig file with endpoint information to point kubectl at any cluster available to the monskope user.

By default the default config file of kubectl will be used ($HOME/.kube/config).

If the KUBECONFIG environment variable is set the file specified will be used. if a list of files is specified you will be asked to choose one.

You can also specify a custom file by utilising the file option (--file). In this case please make sure to update the KUBECONFIG environment variable.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			if err := configManager.LoadConfig(); err != nil {
				return fmt.Errorf("failed loading monoconfig: %w", err)
			}
			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewUpdateKubeconfigUseCase(configManager, kubeConfigPath, overwrite).Run(ctx)
			})
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&kubeConfigPath, "file", "f", "", "the file, in which kubeconfig will be written")
	flags.BoolVarP(&overwrite, "overwrite", "o", false, "Overwrites the existing kubeconfig.")

	return cmd
}
