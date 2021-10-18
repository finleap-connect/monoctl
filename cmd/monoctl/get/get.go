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
	"github.com/spf13/cobra"
	"github.com/finleap-connect/monoctl/internal/output"
)

var showDeleted bool
var sortBy string
var sortDescending bool

func getOutputOptions() *output.OutputOptions {
	sortOpt := output.SortOptions{SortByColumn: sortBy}
	if sortDescending {
		sortOpt.Order = output.Descending
	}
	return &output.OutputOptions{ShowDeleted: showDeleted, SortOptions: sortOpt}
}

func NewGetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:                   "get",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Short:                 "Get information about entities in Monoskope",
		Long:                  `Get information about entities in Monoskope`,
	}

	getCmd.AddCommand(NewGetRolesCmd())
	getCmd.AddCommand(NewGetScopesCmd())
	getCmd.AddCommand(NewGetPoliciesCmd())
	getCmd.AddCommand(NewGetUsersCmd())
	getCmd.AddCommand(NewGetClustersCmd())
	getCmd.AddCommand(NewGetTenantsCmd())
	getCmd.AddCommand(NewGetRoleBindingsCmd())
	getCmd.AddCommand(NewGetTenantUsersCmd())
	getCmd.AddCommand(NewGetClusterCredentials())

	flags := getCmd.PersistentFlags()
	flags.BoolVarP(&showDeleted, "deleted", "d", false, "Show deleted resources.")
	flags.StringVar(&sortBy, "sort-by", "", "Column to sort result by. Uses the first column by default.")
	flags.BoolVar(&sortDescending, "descending", false, "Sort result in descending order.")

	return getCmd
}
