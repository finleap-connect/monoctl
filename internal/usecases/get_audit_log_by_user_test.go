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

	mal "github.com/finleap-connect/monoctl/test/mock/domain"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAuditLog_ByUser", func() {
	var (
		mockCtrl        *gomock.Controller
		expectedServer  = "m8.example.com"
		expectedIssuer  = "admin@monoskope.io"
		expectedUser    = "jane.dou@monoskope.io"
		auditLogOptions = &output.AuditLogOptions{
			MinTime: time.Now(),
			MaxTime: time.Now(),
		}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var testData = []*audit.HumanReadableEvent{
		{
			Timestamp: timestamppb.New(auditLogOptions.MinTime),
			Issuer:    expectedIssuer,
			IssuerId:  uuid.New().String(),
			EventType: events.UserCreated.String(),
			Details:   "UserCreated details for " + expectedUser,
		},
		{
			Timestamp: timestamppb.New(auditLogOptions.MaxTime),
			Issuer:    expectedIssuer,
			IssuerId:  uuid.New().String(),
			EventType: events.UserUpdated.String(),
			Details:   "UserUpdated details for " + expectedUser,
		},
	}

	It("should construct gRPC call to retrieve audit log events affecting a user", func() {
		ctx := context.Background()

		conf := config.NewConfig()
		conf.Server = expectedServer
		conf.AuthInformation = &config.AuthInformation{Token: "this-is-a-token"}

		galUc := NewGetAuditLogByUserUseCase(conf, &output.OutputOptions{}, auditLogOptions, expectedUser).(*getAuditLogByUserUseCase)
		galUc.conn = grpc.CreateDummyGrpcConnection()

		getByUserClient := mal.NewMockAuditLog_GetByUserClient(mockCtrl)
		for _, event := range testData {
			getByUserClient.EXPECT().Recv().Return(event, nil)
		}
		getByUserClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient := mal.NewMockAuditLogClient(mockCtrl)
		mockClient.EXPECT().GetByUser(ctx, &api.GetByUserRequest{
			Email: wrapperspb.String(expectedUser),
			DateRange: &api.GetAuditLogByDateRangeRequest{
				MinTimestamp: timestamppb.New(auditLogOptions.MinTime),
				MaxTimestamp: timestamppb.New(auditLogOptions.MaxTime),
			},
		}).Return(getByUserClient, nil)

		galUc.auditLogClient = mockClient

		err := galUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := galUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(len(testData)))

		tbl.Render()
	})
})
