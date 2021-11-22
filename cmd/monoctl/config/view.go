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
	"fmt"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/spf13/cobra"
)

func NewViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View monoctl config",
		Long:  `View monoctl configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			if err := configManager.LoadConfig(); err != nil {
				return fmt.Errorf("failed loading monoconfig: %w", err)
			}

			fmt.Printf("%s:\n", configManager.GetConfigLocation())

			conf, err := configManager.GetConfig().String()
			if err != nil {
				return err
			}
			fmt.Println(conf)

			return nil
		},
	}

	return cmd
}
