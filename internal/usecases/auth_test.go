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

package usecases

import (
	_ "embed"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
)

// rendered output for certificate resource and issuer
//go:embed expected_CallbackServerSuccessPage.html
var expectedStatusPage string

var _ = Describe("render auth page", func() {
	var (
		mockCtrl *gomock.Controller
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var (
		expectedApiServerAddress = "m8.example.com:443"
	)

	It("should render the index page correctly", func() {

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		confManager := config.NewLoaderFromConfig(conf)

		aUc := NewAuthUsecase(confManager, false).(*authUseCase)

		version := "0.0.1-local"
		commit := "1a2b3c"

		actualStatusPage, err := aUc.renderLocalServerSuccessHTML(expectedApiServerAddress, version, commit)
		Expect(err).ToNot(HaveOccurred())

		Expect(actualStatusPage).To(Equal(expectedStatusPage))
	})
})
