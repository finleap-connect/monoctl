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
	"io"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/output"
	mdom "gitlab.figo.systems/platform/monoskope/monoctl/test/mock/domain"
	api_commandhandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("GetCluster", func() {
	var (
		mockCtrl                    *gomock.Controller
		expectedDisplayName         = "the one cluster"
		expectedName                = "one-cluster"
		expectedApiServerAddress    = "one.example.com"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedUUID                = uuid.New()
		expectedBootstrapToken      = "This should be a JWT"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should construct gRPC call to retrieve cluster data", func() {
		var err error

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		gcUc := NewGetClustersUseCase(conf, &output.OutputOptions{ShowDeleted: true}).(*getClustersUseCase)
		ctx := context.Background()

		// don't do setUp, as this would require a running control plane with working
		// credentials. Instead inject dependencies below
		// err = gcUc.setUp(ctx)
		// Expect(err).ToNot(HaveOccurred())

		// set up dummy connection
		gcUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked commandHandlerClient
		mockClient := mdom.NewMockClusterClient(mockCtrl)

		getAllClient := mdom.NewMockCluster_GetAllClient(mockCtrl)
		getAllClient.EXPECT().Recv().Return(&projections.Cluster{
			Id:               expectedUUID.String(),
			DisplayName:      expectedDisplayName,
			Name:             expectedName,
			ApiServerAddress: expectedApiServerAddress,
			CaCertBundle:     expectedClusterCACertBundle,
			BootstrapToken:   expectedBootstrapToken,
			Metadata: &projections.LifecycleMetadata{
				Created: timestamppb.Now(),
			},
		}, nil)
		getAllClient.EXPECT().Recv().Return(&projections.Cluster{
			Id:               expectedUUID.String(),
			DisplayName:      "another cluster",
			Name:             "another",
			ApiServerAddress: "two.exmaple.com",
			CaCertBundle:     expectedClusterCACertBundle,
			BootstrapToken:   "this-is-anoher-jwt",
			Metadata: &projections.LifecycleMetadata{
				Created: timestamppb.New(time.Date(1975, time.April, 10, 11, 12, 13, 14, time.UTC)),
				Deleted: timestamppb.Now(),
			},
		}, nil)
		getAllClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient.EXPECT().GetAll(ctx, &api_commandhandler.GetAllRequest{
			IncludeDeleted: true,
		}).Return(getAllClient, nil)

		gcUc.client = mockClient

		// SUT
		err = gcUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl := gcUc.tableFactory.ToTable()
		Expect(tbl.NumLines()).To(Equal(2))
		Expect(err).ToNot(HaveOccurred())

		tbl.Render()

	})
	It("should render the cluster data correctly", func() {

	})
})
