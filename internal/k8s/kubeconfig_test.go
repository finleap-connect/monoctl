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

package k8s

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/homedir"
	"os"
	"path"
)

var _ = Describe("Internal/K8s/KubeConfig", func() {
	var kc *KubeConfig
	var tmpFile *os.File

	JustBeforeEach(func() {
		kc = NewKubeConfig()
		Expect(kc).ToNot(BeNil())

		err := os.Setenv("KUBECONFIG", "")
		Expect(err).ToNot(HaveOccurred())

		tmpFile, err = os.CreateTemp("", "kubeconfig")
		Expect(err).ToNot(HaveOccurred())
	})

	JustAfterEach(func() {
		os.Remove(tmpFile.Name())
	})

	CheckConfig := func() {
		conf, err := kc.LoadConfig()
		Expect(err).ToNot(HaveOccurred())
		Expect(conf).ToNot(BeNil())
	}

	Context("loading config", func() {
		It("can the default file", func() {
			CheckConfig()
		})

		It("can use a custom file", func() {
			kc.ConfigPath = tmpFile.Name()
			CheckConfig()
		})

		It("can use a custom file containing env variables and/or home directory shorthand (~)", func() {
			envVar := "M8_TMP_SHORT_PATH"
			err := os.Setenv(envVar, tmpFile.Name())
			defer os.Unsetenv(envVar)
			Expect(err).NotTo(HaveOccurred())
			tmpFileNameShort := "~" + "$" + envVar
			homeDir := homedir.HomeDir()
			Expect(homeDir).ToNot(BeEmpty())

			kc.ConfigPath = tmpFileNameShort
			CheckConfig()
			Expect(kc.ConfigPath).To(Equal(path.Join(homeDir, tmpFile.Name())))
		})

		It("can use the file specified by $"+kubeConfigEnvVar, func() {
			err := os.Setenv(kubeConfigEnvVar, tmpFile.Name())
			Expect(err).ToNot(HaveOccurred())
			CheckConfig()
			Expect(kc.ConfigPath).To(Equal(tmpFile.Name()))
		})

		It("will ask when multiple files are specified in $"+kubeConfigEnvVar, func() {
			err := os.Setenv(kubeConfigEnvVar, tmpFile.Name()+":another/file")
			Expect(err).ToNot(HaveOccurred())
			_, err = kc.LoadConfig()
			Expect(err.Error()).To(Equal("^D")) // asked but no user input
		})
	})
	It("can store config", func() {
		conf, err := kc.LoadConfig()
		Expect(err).ToNot(HaveOccurred())
		Expect(conf).ToNot(BeNil())

		err = kc.StoreConfig(conf)
		Expect(err).ToNot(HaveOccurred())
	})
})
