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
	"fmt"
	"time"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	auth_util "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

func NewGetAuditLogUsersOverviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users-overview [TIMESTAMP]",
		Short: "Get audit log overview of all users at a given point in time.",
		Long:  fmt.Sprintf(`Get audit log overview of all users at a given point in time, tenants/clusters they belong to, and their roles within the system or tenant/cluster. Accepted layout: %s`, now.Format(time.RFC3339)),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			timestamp, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return fmt.Errorf("%s is invalid.\nPlease make sure to use the correct timestamp layout. Example: %s", args[0], now.Format(time.RFC3339))
			}
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			return auth_util.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetAuditLogUsersOverviewUseCase(configManager.GetConfig(), getOutputOptions(), timestamp).Run(ctx)
			})
		},
	}

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		_ = command.Flags().MarkHidden("from")
		_ = command.Flags().MarkHidden("to")
		command.Parent().HelpFunc()(command, strings)
	})
	cmd.SetUsageFunc(func(command *cobra.Command) error {
		_ = command.Flags().MarkHidden("from")
		_ = command.Flags().MarkHidden("to")
		return command.Parent().UsageFunc()(command)
	})

	return cmd
}
