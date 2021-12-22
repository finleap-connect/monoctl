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
	"context"
	_ "embed"
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	mgw "github.com/finleap-connect/monoctl/test/mock/gateway"
	gw "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zalando/go-keyring"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("internal/usecases/create_api_token", func() {
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
		ctx              = context.Background()
		expectedUserId   = uuid.New()
		expectedExpiry   = time.Now().UTC().Add(1 * time.Hour)
		expectedToken    = "some-auth-token"
		expectedScopes   = []string{gw.AuthorizationScope_API.String()}
		expectedValidity = time.Hour * 1
		fakeConfigData   = `server: https://1.1.1.1`
	)

	It("should run", func() {
		var err error

		keyring.MockInit()

		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		confManager := config.NewLoaderFromExplicitFile(tempFile.Path)
		Expect(confManager.LoadConfig()).NotTo(HaveOccurred())

		confManager.GetConfig().AuthInformation = &config.AuthInformation{
			Username: "test-user",
			Expiry:   expectedExpiry,
		}

		apiTokenClient := mgw.NewMockAPITokenClient(mockCtrl)
		uc := NewCreateAPITokenUsecase(confManager, expectedUserId.String(), expectedScopes, expectedValidity).(*createAPITokenUsecase)
		uc.apiTokenClient = apiTokenClient
		uc.setInitialized()

		apiTokenClient.EXPECT().RequestAPIToken(ctx, gomock.Any()).Return(&gw.APITokenResponse{
			AccessToken: expectedToken,
			Expiry:      timestamppb.New(expectedExpiry),
		}, nil)

		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

})
