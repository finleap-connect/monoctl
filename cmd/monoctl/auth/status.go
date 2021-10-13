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

package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoctl/cmd/monoctl/flags"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
)

func NewAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  `Shows if authenticated against any Monoskope instance and against which one.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			if err := configManager.LoadConfig(); err != nil {
				return fmt.Errorf("failed loading monoconfig: %w", err)
			}

			conf := configManager.GetConfig()
			authenticated := conf.HasAuthInformation() && conf.AuthInformation.HasToken()

			fmt.Printf("Authenticated: %v\n", authenticated)
			if authenticated {
				fmt.Printf("Server: %v\n", conf.Server)
				fmt.Printf("Token expiry: %v\n", conf.AuthInformation.Expiry)
				fmt.Printf("Token expired: %v\n", conf.AuthInformation.IsTokenExpired())
			}

			return nil
		},
	}
}
