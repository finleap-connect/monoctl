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

package main

import (
	"flag"
	"time"

	"github.com/spf13/cobra"
	"github.com/finleap-connect/monoctl/cmd/monoctl/auth"
	conf "github.com/finleap-connect/monoctl/cmd/monoctl/config"
	"github.com/finleap-connect/monoctl/cmd/monoctl/create"
	"github.com/finleap-connect/monoctl/cmd/monoctl/delete"
	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/cmd/monoctl/get"
	"github.com/finleap-connect/monoctl/cmd/monoctl/update"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "monoctl [command] [flags]",
		Short: "monoctl",
		Long: `
monoctl is the CLI of Monoskope. It allows the management and operation of tenants,
users and their roles in a Kubernetes multi-cluster environment.`,
		DisableAutoGenTag: true,
	}

	// Setup global flags
	fl := rootCmd.PersistentFlags()
	fl.AddGoFlagSet(flag.CommandLine)
	fl.StringVar(&flags.ExplicitFile, "monoconfig", "", "Path to explicit monoskope config file to use for CLI requests")
	fl.DurationVar(&flags.Timeout, "command-timeout", 10*time.Second, "Timeout for long running commands")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewCompletionCommand())

	rootCmd.AddCommand(conf.NewConfigCmd())
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(get.NewGetCmd())
	rootCmd.AddCommand(create.NewCreateCmd())
	rootCmd.AddCommand(update.NewUpdateCmd())
	rootCmd.AddCommand(delete.NewDeleteCmd())

	return rootCmd
}
