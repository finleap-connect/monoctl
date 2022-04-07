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
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	mal "github.com/finleap-connect/monoctl/test/mock/domain"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("GetAuditLog_UsersOverview", func() {
	var (
		mockCtrl *gomock.Controller
		expectedServer = "m8.example.com"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var testData = []*audit.UserOverview{
		{
			Name: "Test User",
			Email: "test.user@monoskope.io",
			Roles: "user roles",
			Tenants: "the tenants to which the user has access",
			Clusters: "the clusters to which the user has access",
			Details: "UserOverview details",
		},
		{
			Name: "Another User",
			Email: "another.user@monoskope.io",
			Roles: "user roles",
			Tenants: "the tenants to which the user has access",
			Clusters: "the clusters to which the user has access",
			Details: "UserOverview details",
		},
	}

	It("should construct gRPC call to retrieve audit log of users overview", func() {
		ctx := context.Background()

		conf := config.NewConfig()
		conf.Server = expectedServer
		conf.AuthInformation = &config.AuthInformation{Token: "this-is-a-token"}

		galUc := NewGetAuditLogUsersOverviewUseCase(conf, &output.OutputOptions{ShowDeleted: true}).(*getAuditLogUsersOverviewUseCase)
		galUc.conn = grpc.CreateDummyGrpcConnection()

		getUsersOverviewClient := mal.NewMockAuditLog_GetUsersOverviewClient(mockCtrl)
		for _, overview := range testData {
			getUsersOverviewClient.EXPECT().Recv().Return(overview, nil)
		}
		getUsersOverviewClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient := mal.NewMockAuditLogClient(mockCtrl)
		mockClient.EXPECT().GetUsersOverview(ctx, &api.GetAllRequest{IncludeDeleted: true}).
			Return(getUsersOverviewClient, nil)

		galUc.auditLogClient = mockClient

		err := galUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := galUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(len(testData)))

		tbl.Render()
	})
})
