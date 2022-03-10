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
	"github.com/finleap-connect/monoctl/internal/output"
	"github.com/spf13/cobra"
)

var showDeleted bool
var sortBy string
var sortDescending bool
var exportFile string

func getOutputOptions() *output.OutputOptions {
	sortOpt := output.SortOptions{SortByColumn: sortBy}
	if sortDescending {
		sortOpt.Order = output.Descending
	}
	exportOpt := output.ExportOptions{
		Format: output.CSV,
		File: exportFile,
	}
	return &output.OutputOptions{ShowDeleted: showDeleted, SortOptions: sortOpt, ExportOptions: exportOpt}
}

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "Get information about entities in Monoskope",
		Long:                  `Get information about entities in Monoskope`,
	}

	cmd.AddCommand(NewGetRolesCmd())
	cmd.AddCommand(NewGetScopesCmd())
	cmd.AddCommand(NewGetPoliciesCmd())
	cmd.AddCommand(NewGetUsersCmd())
	cmd.AddCommand(NewGetClustersCmd())
	cmd.AddCommand(NewGetTenantsCmd())
	cmd.AddCommand(NewGetRoleBindingsCmd())
	cmd.AddCommand(NewGetTenantUsersCmd())
	cmd.AddCommand(NewGetClusterCredentials())
	cmd.AddCommand(NewGetClusterAccess())
	cmd.AddCommand(NewGetAuditLogCmd())

	flags := cmd.PersistentFlags()
	flags.BoolVarP(&showDeleted, "deleted", "d", false, "Show deleted resources.")
	flags.StringVar(&sortBy, "sort-by", "", "Column to sort result by. Uses the first column by default.")
	flags.BoolVar(&sortDescending, "descending", false, "Sort result in descending order.")
	flags.StringVar(&exportFile, "export", "", "exports the result to a file in csv format. If no file is specified m8-output.csv will be written in the current directory if it doesn't exists")
	flags.Lookup("export").NoOptDefVal = "m8-output.csv"

	return cmd
}
