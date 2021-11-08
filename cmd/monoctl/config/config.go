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
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "View and modify monoctl config files",
		Long:                  `View and modify monoctl config files`,
	}
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewViewCmd())
	return cmd
}
