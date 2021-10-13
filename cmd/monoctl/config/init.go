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

package config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.figo.systems/platform/monoskope/monoctl/cmd/monoctl/flags"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
)

var (
	serverURL string
	force     bool
)

func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init monoctl config",
		Long:  `Init monoctl and create a new monoskope configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURL == "" {
				return errors.New("failed initializing monoconfig: server-url is required")
			}

			cfg := config.NewConfig()
			cfg.Server = serverURL

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)
			if err := configManager.InitConfig(cfg, force); err != nil {
				return fmt.Errorf("failed initializing monoconfig: %w", err)
			}

			return nil
		},
	}

	flags := initCmd.Flags()
	flags.StringVarP(&serverURL, "server-url", "u", "", "URL of the monoskope instance")
	flags.BoolVarP(&force, "force", "f", false, "Force overwrite configuration.")

	err := initCmd.MarkFlagRequired("server-url")
	if err != nil {
		panic(err)
	}

	return initCmd
}
