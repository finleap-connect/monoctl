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

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	mdom "github.com/finleap-connect/monoctl/test/mock/domain"
	api_commandhandler "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("GetTenants", func() {
	var (
		mockCtrl       *gomock.Controller
		expectedName   = "The Power Team"
		expectedPrefix = "TPT"
		expectedUUID   = uuid.New()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var testData = []*projections.Tenant{
		{
			Id:       expectedUUID.String(),
			Name:     expectedName,
			Prefix:   expectedPrefix,
			Metadata: &projections.LifecycleMetadata{Created: timestamppb.Now()},
		},
		{
			Id:     expectedUUID.String(),
			Name:   "another tenant",
			Prefix: "TOT",
			Metadata: &projections.LifecycleMetadata{
				Created: timestamppb.New(time.Date(1975, time.April, 10, 11, 12, 13, 14, time.UTC)),
				Deleted: timestamppb.Now(),
			},
		},
	}

	It("should construct gRPC call to retrieve tenant data (included deleted)", func() {
		var err error

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		gtUc := NewGetTenantsUseCase(conf, &output.OutputOptions{ShowDeleted: true}).(*getTenantsUseCase)
		ctx := context.Background()

		// set up dummy connection
		gtUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked commandHandlerClient
		mockClient := mdom.NewMockTenantClient(mockCtrl)

		getAllClient := mdom.NewMockTenant_GetAllClient(mockCtrl)
		for _, tenant := range testData {
			getAllClient.EXPECT().Recv().Return(tenant, nil)
		}
		getAllClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient.EXPECT().GetAll(ctx, &api_commandhandler.GetAllRequest{
			IncludeDeleted: true,
		}).Return(getAllClient, nil)

		gtUc.client = mockClient

		// SUT
		err = gtUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := gtUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(2))

		tbl.Render()
	})
	It("should construct gRPC call to retrieve tenant data (included deleted)", func() {
		var err error

		includeDeleted := false

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		gtUc := NewGetTenantsUseCase(conf, &output.OutputOptions{ShowDeleted: includeDeleted}).(*getTenantsUseCase)
		ctx := context.Background()

		// set up dummy connection
		gtUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked commandHandlerClient
		mockClient := mdom.NewMockTenantClient(mockCtrl)

		getAllClient := mdom.NewMockTenant_GetAllClient(mockCtrl)
		for _, tenant := range testData {
			getAllClient.EXPECT().Recv().Return(tenant, nil)
		}
		getAllClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient.EXPECT().GetAll(ctx, &api_commandhandler.GetAllRequest{
			IncludeDeleted: includeDeleted,
		}).Return(getAllClient, nil)

		gtUc.client = mockClient

		// SUT
		err = gtUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := gtUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(2))

		tbl.Render()
	})
})
