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

package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	keyring "github.com/zalando/go-keyring"
)

var _ = Describe("client config loader", func() {
	keyring.MockInit()
	fakeConfigData := `server: https://1.1.1.1`

	It("can load config from bytes", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).ToNot(BeNil())
	})
	It("errors for empty config", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(""))
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ErrEmptyServer))
		Expect(conf).To(BeNil())
	})
	It("loads config from env var path", func() {
		loader := NewLoader()

		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		os.Setenv(RecommendedConfigPathEnvVar, tempFile.Path)
		err = loader.LoadConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("loads config from explicit file path", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can init config for explicit file path", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		conf := NewConfig()
		conf.Server = "localhost"

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		os.Remove(tempFile.Path)
		err = loader.InitConfig(conf, false)

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can init config for env var path", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoader()
		conf := NewConfig()
		conf.Server = "localhost"

		confPath := path.Join(tempFile.Path, ".monoskope", "config")
		os.Setenv(RecommendedConfigPathEnvVar, confPath)
		os.Remove(tempFile.Path)
		err = loader.InitConfig(conf, false)

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can save config", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())

		conf := loader.GetConfig()
		conf.Server = "monoskope.io"

		conf.AuthInformation = &AuthInformation{}
		conf.AuthInformation.Username = "user"
		conf.AuthInformation.Token = "token"
		conf.AuthInformation.Expiry = time.Now().UTC().Add(1 * time.Hour)

		expectedClusterId := uuid.New()
		expectedClusterRole := "default"
		expectedClusterAuthInfo := &AuthInformation{
			Username: "clusteruser",
			Token:    "clustertoken",
			Expiry:   time.Now().UTC().Add(1 * time.Hour),
		}
		conf.ClusterAuthInformation[fmt.Sprintf("%s/%s/%s", expectedClusterId.String(), expectedClusterAuthInfo.Username, expectedClusterRole)] = expectedClusterAuthInfo

		err = loader.SaveConfig()
		Expect(err).NotTo(HaveOccurred())

		loader = NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
		Expect(loader.config.Server).To(Equal(conf.Server))
		Expect(loader.config.AuthInformation.Username).To(Equal(conf.AuthInformation.Username))
		Expect(loader.config.AuthInformation.Token).To(Equal(conf.AuthInformation.Token))

		clusterAuthInfoFromFile := loader.config.GetClusterAuthInformation(expectedClusterId.String(), expectedClusterAuthInfo.Username, expectedClusterRole)
		Expect(clusterAuthInfoFromFile).ToNot(BeNil())
		Expect(expectedClusterAuthInfo.Username).To(Equal(clusterAuthInfoFromFile.Username))
		Expect(expectedClusterAuthInfo.Token).To(Equal(clusterAuthInfoFromFile.Token))

	})
	It("can validate", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).ToNot(BeNil())

		Expect(conf.HasAuthInformation()).To(BeFalse())

		conf.AuthInformation = &AuthInformation{}
		Expect(conf.HasAuthInformation()).To(BeTrue())
		Expect(conf.AuthInformation.HasToken()).To(BeFalse())
		Expect(conf.AuthInformation.IsValid()).To(BeFalse())

		conf.AuthInformation.Token = "test"
		conf.AuthInformation.Expiry = time.Now().Add(1 * time.Hour)
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.IsTokenExpired()).To(BeFalse())
		Expect(conf.AuthInformation.IsValid()).To(BeTrue())

		conf.AuthInformation.Expiry = time.Now().Add(-1 * time.Hour)
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.IsTokenExpired()).To(BeTrue())
		Expect(conf.AuthInformation.IsValid()).To(BeFalse())
	})
})
