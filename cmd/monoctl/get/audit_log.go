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
	authutil "github.com/finleap-connect/monoctl/internal/util/auth"
	"github.com/spf13/cobra"
)

func NewGetAuditLogCmd() *cobra.Command {
	const dateLayoutISO8601 = "2006-01-02" // don't change. This corresponds to YYYY-MM-DD
	var (
		from string
		to string

		now = time.Now()
		firstOfMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth = firstOfMonth.AddDate(0, 1, -1)
		dateInputErr = func(input string) error { return fmt.Errorf("%s is invalid.\nPlease make sure to use the correct date layout. Example: %s", input, now.Format(dateLayoutISO8601))}
	)

	cmd := &cobra.Command{
		Use:     "audit-log",
		Aliases: []string{"audit"},
		Short:   "Get audit log",
		Long:    `Get audit log based on a date range. If no date range is specified the audit log of the current month will be returned. Please note that both input and output time is in UTC and should be considered when e.g. specifying the date range`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			configManager := config.NewLoaderFromExplicitFile(flags.ExplicitFile)

			return authutil.RetryOnAuthFail(cmd.Context(), configManager, func(ctx context.Context) error {
				return usecases.NewGetAuditLogUseCase(configManager.GetConfig(), getOutputOptions(), minTime, maxTime).Run(ctx)
			})
		},
	}

	cmdFlags := cmd.Flags()
	cmdFlags.StringVarP(&from, "from", "f", firstOfMonth.Format(dateLayoutISO8601),
		fmt.Sprintf("Specifys the starting point of the date range. If not specified the first day of the current month is used. Accepted layout: %s", now.Format(dateLayoutISO8601)))
	cmdFlags.StringVarP(&to, "to", "t", lastOfMonth.Format(dateLayoutISO8601),
		fmt.Sprintf("Specifys the ending point of the date range. If not specified the last day of the current month is used. Accepted layout: %s", now.Format(dateLayoutISO8601)))

	return cmd
}
