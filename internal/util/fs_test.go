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

package util

import (
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("util.fs", func() {
	It("can check if file exists", func() {
		tempDir, err := testutil_fs.NewTempDir()
		Expect(err).NotTo(HaveOccurred())
		defer tempDir.Close()

		exists, err := FileExists(tempDir.Path)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		tempDir.Close()
		exists, err = FileExists(tempDir.Path)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})
	It("can determine homedir", func() {
		Expect(HomeDir()).NotTo(BeEmpty())
	})
	It("can create a file only if doesn't exist", func() {
		tempFile, err := ioutil.TempFile(os.TempDir(), "m8-")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tempFile.Name())

		exists, err := FileExists(tempFile.Name())
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		file, err := NewFileSafe(tempFile.Name())
		Expect(err).To(HaveOccurred())
		Expect(file).To(BeNil())

		err = os.Remove(tempFile.Name())
		Expect(err).ToNot(HaveOccurred())

		file, err = NewFileSafe(tempFile.Name())
		Expect(err).ToNot(HaveOccurred())

		exists, err = FileExists(tempFile.Name())
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())
	})
})
