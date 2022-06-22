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
	mdomain "github.com/finleap-connect/monoctl/test/mock/domain"
	mgw "github.com/finleap-connect/monoctl/test/mock/gateway"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	gw "github.com/finleap-connect/monoskope/pkg/api/gateway"
	mk8s "github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zalando/go-keyring"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("GetClusterCredentials", func() {
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
		ctx                  = context.Background()
		expectedExpiry       = time.Now().UTC().Add(1 * time.Hour)
		expectedClusterToken = "some-auth-token"
		fakeConfigData       = `server: https://1.1.1.1`
		expectedRole         = string(mk8s.DefaultRole)
		// expectedAdminRole    = string(mk8s.AdminRole)
	)

	getClusters := func() []*projections.Cluster {
		return []*projections.Cluster{
			{
				Id:               uuid.New().String(),
				DisplayName:      "First Cluster",
				Name:             "first-cluster",
				ApiServerAddress: "first.cluster.monokope.io",
				CaCertBundle:     []byte("This should be a certificate"),
			},
			{
				Id:               uuid.New().String(),
				DisplayName:      "Second Cluster",
				Name:             "second-cluster",
				ApiServerAddress: "second.cluster.monokope.io",
				CaCertBundle:     []byte("This should be a certificate"),
			},
		}
	}

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

		mockClusterClient := mdomain.NewMockClusterClient(mockCtrl)
		mockClusterAuthClient := mgw.NewMockClusterAuthClient(mockCtrl)

		expectedClusters := getClusters()

		// mockClusterClient.EXPECT().GetByName(ctx, wrapperspb.String(expectedClusters[0].Name)).Return(expectedClusters[0], nil)

		// getAllClient := mdomain.NewMockCluster_GetAllClient(mockCtrl)
		// for _, expectedCluster := range expectedClusters {
		// 	getAllClient.EXPECT().Recv().Return(expectedCluster, nil)
		// }
		// getAllClient.EXPECT().Recv().Return(nil, io.EOF)
		// mockClusterClient.EXPECT().GetAll(ctx, &api.GetAllRequest{IncludeDeleted: false}).Return(getAllClient, nil)

		for _, expectedCluster := range expectedClusters {
			mockClusterAuthClient.EXPECT().GetAuthToken(ctx, &gw.ClusterAuthTokenRequest{
				ClusterId: expectedCluster.Id,
				Role:      expectedRole,
			}).Return(&gw.ClusterAuthTokenResponse{
				AccessToken: expectedClusterToken,
				Expiry:      timestamppb.New(expectedExpiry),
			}, nil)
		}

		uc := NewGetClusterCredentialsUseCase(confManager, expectedClusters[0].Id, expectedRole).(*getClusterCredentialsUseCase)
		uc.clusterServiceClient = mockClusterClient
		uc.clusterAuthClient = mockClusterAuthClient
		uc.setInitialized()
		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())

		uc = NewGetClusterCredentialsUseCase(confManager, expectedClusters[1].Id, expectedRole).(*getClusterCredentialsUseCase)
		uc.clusterServiceClient = mockClusterClient
		uc.clusterAuthClient = mockClusterAuthClient
		uc.setInitialized()
		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())

		c := confManager.GetConfig()
		Eventually(func(g Gomega) {
			for _, expectedCluster := range expectedClusters {
				g.Expect(c.GetClusterAuthInformation(expectedCluster.Id, c.AuthInformation.Username, expectedRole)).ToNot(BeNil())
			}
		}).Should(Succeed())
	})

	// It("should get all clusters credentials for default role only", func() {
	// 	var err error

	// 	keyring.MockInit()

	// 	tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
	// 	Expect(err).NotTo(HaveOccurred())
	// 	defer tempFile.Close()

	// 	confManager := config.NewLoaderFromExplicitFile(tempFile.Path)
	// 	Expect(confManager.LoadConfig()).NotTo(HaveOccurred())

	// 	confManager.GetConfig().AuthInformation = &config.AuthInformation{
	// 		Username: "test-user",
	// 		Expiry:   expectedExpiry,
	// 	}

	// 	mockClusterClient := mdomain.NewMockClusterClient(mockCtrl)
	// 	mockClusterAuthClient := mgw.NewMockClusterAuthClient(mockCtrl)

	// 	expectedClusters := getClusters()
	// 	expectedClusterAdmin := expectedClusters[0]

	// 	uc := NewGetClusterCredentialsUseCase(confManager, expectedClusterAdmin.Id, expectedAdminRole).(*getClusterCredentialsUseCase)
	// 	uc.clusterServiceClient = mockClusterClient
	// 	uc.clusterAuthClient = mockClusterAuthClient
	// 	uc.setInitialized()

	// 	mockClusterClient.EXPECT().GetByName(ctx, wrapperspb.String(expectedClusterAdmin.Name)).Return(expectedClusterAdmin, nil)

	// 	mockClusterAuthClient.EXPECT().GetAuthToken(ctx, &gw.ClusterAuthTokenRequest{
	// 		ClusterId: expectedClusterAdmin.Id,
	// 		Role:      expectedAdminRole,
	// 	}).Return(&gw.ClusterAuthTokenResponse{
	// 		AccessToken: expectedClusterToken,
	// 		Expiry:      timestamppb.New(expectedExpiry),
	// 	}, nil)

	// 	err = uc.Run(ctx)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	c := confManager.GetConfig()
	// 	for _, expectedCluster := range expectedClusters {
	// 		authInfo := c.GetClusterAuthInformation(expectedCluster.Id, c.AuthInformation.Username, expectedAdminRole)
	// 		if expectedClusterAdmin == expectedCluster {
	// 			Expect(authInfo).ToNot(BeNil())
	// 		} else {
	// 			Expect(authInfo).To(BeNil())
	// 		}
	// 	}
	// })
})
