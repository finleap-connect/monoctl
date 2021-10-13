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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zalando/go-keyring"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	mdomain "gitlab.figo.systems/platform/monoskope/monoctl/test/mock/domain"
	mgw "gitlab.figo.systems/platform/monoskope/monoctl/test/mock/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	mk8s "gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("CreateKubeconfig", func() {
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
		ctx                         = context.Background()
		expectedId                  = uuid.New()
		expectedDisplayName         = "Test Cluster"
		expectedName                = "test-cluster"
		expectedApiServerAddress    = "test.cluster.monokope.io"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedExpiry              = time.Now().UTC().Add(1 * time.Hour)
		expectedClusterToken        = "some-auth-token"
	)

	It("should run", func() {
		var err error

		keyring.MockInit()
		fakeConfigData := `server: https://1.1.1.1`

		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		confManager := config.NewLoaderFromExplicitFile(tempFile.Path)
		Expect(confManager.LoadConfig()).NotTo(HaveOccurred())

		confManager.GetConfig().AuthInformation = &config.AuthInformation{
			Username: "test-user",
			Expiry:   expectedExpiry,
		}

		mockClusterClient := mdomain.NewMockClusterClient(mockCtrl)
		mockClusterAuthClient := mgw.NewMockClusterAuthClient(mockCtrl)

		uc := NewGetClusterCredentialsUseCase(confManager, expectedDisplayName, string(mk8s.DefaultRole)).(*getClusterCredentialsUseCase)
		uc.clusterServiceClient = mockClusterClient
		uc.clusterAuthClient = mockClusterAuthClient
		uc.setInitialized()

		mockClusterClient.EXPECT().GetByName(ctx, wrapperspb.String(expectedDisplayName)).Return(&projections.Cluster{
			Id:               expectedId.String(),
			DisplayName:      expectedDisplayName,
			Name:             expectedName,
			ApiServerAddress: expectedApiServerAddress,
			CaCertBundle:     expectedClusterCACertBundle,
		}, nil)

		mockClusterAuthClient.EXPECT().GetAuthToken(ctx, &gw.ClusterAuthTokenRequest{
			ClusterId: expectedId.String(),
			Role:      string(mk8s.DefaultRole),
		}).Return(&gw.ClusterAuthTokenResponse{
			AccessToken: expectedClusterToken,
			Expiry:      timestamppb.New(expectedExpiry),
		}, nil)

		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

})
