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
	"encoding/base64"
	"encoding/json"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	mes "github.com/finleap-connect/monoctl/test/mock/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: remove focus test
var _ = Describe("Get Audit Log", func() {
	var (
		mockCtrl               *gomock.Controller
		excepectedCreatorEmail = "admin@monoskope.io"
		excepectedCreatorID    = uuid.New().String()
		excepectedCreatedEmail = "user@monoskope.io"
		expectedCreatedID      = uuid.New().String()
		data = es.ToEventDataFromProto(&eventdata.UserCreated{
			Name: "user.monoskope",
			Email: excepectedCreatedEmail,
		})
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var testData = []*eventsourcing.Event{
		{
			Type:             events.UserCreated.String(),
			Timestamp:        timestamppb.New(time.Date(2021, time.December, 10, 23, 14, 13, 14, time.UTC)),
			AggregateId:      expectedCreatedID,
			AggregateType:    aggregates.User.String(),
			AggregateVersion: &wrapperspb.UInt64Value{Value: 1},
			Data: data,
			Metadata: map[string]string{
				"x-auth-email": excepectedCreatorEmail,
				"x-auth-id": excepectedCreatorID,
			},
		},
		{
			Type:             events.UserRoleBindingCreated.String(),
			Timestamp:        timestamppb.New(time.Date(2021, time.December, 10, 23, 18, 13, 14, time.UTC)),
			AggregateId:      uuid.New().String(),
			AggregateType:    aggregates.UserRoleBinding.String(),
			AggregateVersion: &wrapperspb.UInt64Value{Value: 1},
			Data: es.ToEventDataFromProto(&eventdata.UserRoleAdded{
				Role: roles.Admin.String(),
				Scope: scopes.System.String(),
				UserId: expectedCreatedID,
			}),
			Metadata: map[string]string{
				"x-auth-email": "system@monoskope.local",
				"x-auth-id": uuid.New().String(),
			},
		},
	}

	// TODO: delete this
	testData = readFromFile()

	It("should construct gRPC call to retrieve tenant data (included deleted)", func() {
		var err error

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		galUc := NewGetAuditLogUseCase(conf, &output.OutputOptions{}).(*getAuditLogUseCase)
		ctx := context.Background()

		// set up dummy connection
		galUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked eventStoreClient
		mockClient := mes.NewMockEventStoreClient(mockCtrl)

		retrieveClient := mes.NewMockEventStore_RetrieveClient(mockCtrl)
		for _, event := range testData {
			retrieveClient.EXPECT().Recv().Return(event, nil)
		}
		retrieveClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClient.EXPECT().Retrieve(ctx, &esApi.EventFilter{}).Return(retrieveClient, nil)

		galUc.client = mockClient

		// SUT
		err = galUc.doRun(ctx)
		Expect(err).ToNot(HaveOccurred())

		tbl, err := galUc.tableFactory.ToTable()
		Expect(err).ToNot(HaveOccurred())
		Expect(tbl.NumLines()).To(Equal(len(testData)))

		tbl.Render()
	})
})

// TODO: delete this
func readFromFile() []*eventsourcing.Event {
	jsonFile, err := os.Open("/Users/hanialshikh/Desktop/FCloud/m8/monoskope/event_store_dump.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var events_ []map[string]interface{}
	err = json.Unmarshal(byteValue, &events_)
	if err != nil {
		return nil
	}

	var testData []*eventsourcing.Event
	for _, e := range events_ {

		metadata := make(map[string]string)
		for k, v := range e["metadata"].(map[string]interface{}) {
			metadata[k] = v.(string)
		}

		var data []byte
		if d, ok := e["data"].(string); ok {
			data, _ = base64.StdEncoding.DecodeString(d)
		}

		testData = append(testData, &eventsourcing.Event{
			Type: e["type"].(string),
			AggregateId: e["aggregateId"].(string),
			AggregateType: e["aggregateType"].(string),
			Data: data,
			Metadata: metadata,
		})
	}

	return testData
}
