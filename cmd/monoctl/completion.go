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
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCommand() *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Output shell completion code for the specified shell (bash, zsh or fish)",
		Long: `To load completions:
	
	Bash:
	
	  $ source <(monoctl completion bash)
	
	  # To load completions for each session, execute once:
	  # Linux:
	  $ monoctl completion bash > /etc/bash_completion.d/monoctl
	  # macOS:
	  $ monoctl completion bash > /usr/local/etc/bash_completion.d/monoctl
	
	Zsh:
	
	  # If shell completion is not already enabled in your environment,
	  # you will need to enable it.  You can execute the following once:
	
	  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	  # To load completions for each session, execute once:
	  $ monoctl completion zsh > "${fpath[1]}/_monoctl"
	
	  # You will need to start a new shell for this setup to take effect.
	
	fish:
	
	  $ monoctl completion fish | source
	
	  # To load completions for each session, execute once:
	  $ monoctl completion fish > ~/.config/fish/completions/monoctl.fish
	
	PowerShell:
	
	  PS> monoctl completion powershell | Out-String | Invoke-Expression
	
	  # To load completions for every new session, run:
	  PS> monoctl completion powershell > monoctl.ps1
	  # and source this file from your PowerShell profile.
	`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	return completionCmd
}
