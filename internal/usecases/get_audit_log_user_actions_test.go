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
	mal "github.com/finleap-connect/monoctl/test/mock/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAuditLog_UserActions", func() {
	var (
		mockCtrl *gomock.Controller
		expectedServer = "m8.example.com"
		expectedIssuer = "admin@monoskope.io"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var testData = []*audit.HumanReadableEvent{
		{
			When: time.Now().Format(time.RFC822),
			Issuer: expectedIssuer,
			IssuerId: uuid.New().String(),
			EventType: events.UserCreated.String(),
			Details: "UserCreated details",
		},
		{
			When: time.Now().Format(time.RFC822),
			Issuer: expectedIssuer,
			IssuerId: uuid.New().String(),
			EventType: events.TenantCreated.String(),
			Details: "TenantCreated details",
		},
	}

	It("should construct gRPC call to retrieve audit log events caused by a user actions", func() {
		ctx := context.Background()

		conf := config.NewConfig()
		conf.Server = expectedServer
		conf.AuthInformation = &config.AuthInformation{Token: "this-is-a-token"}

		galUc := NewGetAuditLogUserActionsUseCase(conf, &output.OutputOptions{}, expectedIssuer).(*getAuditLogUserActionsUseCase)
		galUc.conn = grpc.CreateDummyGrpcConnection()

		getUserActionsClient := mal.NewMockAuditLog_GetUserActionsClient(mockCtrl)
		for _, event := range testData {
			getUserActionsClient.EXPECT().Recv().Return(event, nil)
		}
		getUserActionsClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient := mal.NewMockAuditLogClient(mockCtrl)
		mockClient.EXPECT().GetUserActions(ctx, wrapperspb.String(expectedIssuer)).Return(getUserActionsClient, nil)

		galUc.auditLogClient = mockClient

		err := galUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := galUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(len(testData)))

		tbl.Render()
	})
})
