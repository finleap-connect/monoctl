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

	"github.com/finleap-connect/monoctl/internal/output"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	authutil "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

const dateLayoutISO8601 = "2006-01-02" // don't change. This corresponds to YYYY-MM-DD
var (
	from string
	to   string

	now          = time.Now().UTC()
	firstOfMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth  = firstOfMonth.AddDate(0, 1, -1)
	dateInputErr = func(input string) error {
		return fmt.Errorf("%s is invalid.\nPlease make sure to use the correct date layout. Example: %s", input, now.Format(dateLayoutISO8601))
	}
)

func getAuditLogOptions() (*output.AuditLogOptions, error) {
	auditLogOptions := &output.AuditLogOptions{}
	err := parseDateRange(auditLogOptions)
	if err != nil {
		return nil, err
	}
	return auditLogOptions, nil
}

func NewGetAuditLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "audit-log",
		Aliases: []string{"audit"},
		Short:   "Get audit log",
		Long:    `Get audit log within specified date range. If no date range is specified the audit log of the current month will be returned.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			auditLogOptions, err := getAuditLogOptions()
			if err != nil {
				return err
			}
			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			return authutil.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetAuditLogUseCase(configManager.GetConfig(), getOutputOptions(), auditLogOptions).Run(ctx)
			})
		},
	}

	cmd.AddCommand(NewGetAuditLogByUserCmd())
	cmd.AddCommand(NewGetAuditLogUserActionsCmd())
	cmd.AddCommand(NewGetAuditLogUsersOverviewCmd())

	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&from, "from", "f", firstOfMonth.Format(dateLayoutISO8601),
		fmt.Sprintf("Specifies the starting point of the date range (UTC). If not specified the first day of the current month is used. Accepted layout: %s", now.Format(dateLayoutISO8601)))
	persistentFlags.StringVarP(&to, "to", "t", lastOfMonth.Format(dateLayoutISO8601),
		fmt.Sprintf("Specifies the ending point of the date range (UTC). If not specified the last day of the current month is used. Accepted layout: %s", now.Format(dateLayoutISO8601)))

	return cmd
}

func parseDateRange(auditLogOptions *output.AuditLogOptions) error {
	minTime, err := time.Parse(dateLayoutISO8601, from)
	if err != nil {
		if len(from) != 0 { // if not specified first day of current month is used
			return dateInputErr(from)
		}
		minTime = firstOfMonth
	}
	maxTime, err := time.Parse(dateLayoutISO8601, to)
	if err != nil {
		if len(to) != 0 { // if not specified last day of the current month is used
			return dateInputErr(to)
		}
		maxTime = lastOfMonth
	}

	auditLogOptions.MinTime = minTime
	auditLogOptions.MaxTime = maxTime
	return nil
}
