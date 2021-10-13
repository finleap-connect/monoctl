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

package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	Charset  = spinner.CharSets[14]
	Duration = 100 * time.Millisecond
)

// NewSpinner creates and starts the new default cli spinner. Usage:
func NewSpinner() *spinner.Spinner {
	s := spinner.New(Charset, Duration) // Build our new spinner
	s.Start()                           // Start the spinner
	return s
}
