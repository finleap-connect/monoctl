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

package auth

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "embed"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Auth", func() {
	Context("should not run multiple authentication flows simultaneously", func() {
		Specify("FirstProcess", func() {
			lock, err := acquireLock(context.Background(), false)
			Expect(err).ToNot(HaveOccurred())

			ex, err := os.Executable()
			Expect(err).ToNot(HaveOccurred())
			err = os.Setenv("GINKGO_EDITOR_INTEGRATION", "true") // exit status 197 - Programmatic Focus
			Expect(err).ToNot(HaveOccurred())
			cmd := exec.Command(os.Args[0], "-ginkgo.focus", "SecondProcess", filepath.Dir(ex))
			cmd.Env = append(
				os.Environ(), // os.Environ() must be preserved on Windows otherwise "it will fail in weird and wonderful ways"
				"RUN_NAMED_MUTEX_TEST_HELPER=true",
			)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			timeOut := 5 * time.Second // the default 1s is too fast for gexec.Start
			Eventually(session, timeOut).Should(gbytes.Say("Another monoctl instance is already running the authentication flow..."))

			lock.Release()
			Eventually(session, timeOut).Should(gexec.Exit(0))
		})

		Specify("SecondProcess", func() {
			skip := os.Getenv("RUN_NAMED_MUTEX_TEST_HELPER") != "true"
			if skip {
				Skip("This is to be run in a separate process")
			}

			_, err := acquireLock(context.Background(), false)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
