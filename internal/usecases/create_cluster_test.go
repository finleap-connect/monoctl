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

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	mes "github.com/finleap-connect/monoctl/test/mock/eventsourcing"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	es "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateCluster", func() {
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
		expectedName                = "one-cluster"
		expectedApiServerAddress    = "one.example.com"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedServer              = "m8.example.com"
	)

	It("should construct a gRPC call", func() {
		var err error

		conf := config.NewConfig()
		conf.Server = expectedServer
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		ccUc := NewCreateClusterUseCase(conf, expectedName,
			expectedApiServerAddress, expectedClusterCACertBundle).(*createClusterUseCase)

		Expect(ccUc.name).To(Equal(expectedName))
		Expect(ccUc.apiServerAddress).To(Equal(expectedApiServerAddress))
		Expect(ccUc.caCertBundle).To(Equal(expectedClusterCACertBundle))

		ctx := context.Background()

		// set up dummy connection
		ccUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked commandHandlerClient
		mockClient := mes.NewMockCommandHandlerClient(mockCtrl)
		commanddata := &cmdData.CreateCluster{
			Name:             expectedName,
			ApiServerAddress: expectedApiServerAddress,
			CaCertBundle:     expectedClusterCACertBundle,
		}

		command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateCluster, commanddata)
		Expect(err).ToNot(HaveOccurred())

		generatedId := uuid.New().String()
		reply := &es.CommandReply{
			AggregateId: generatedId,
			Version:     0,
		}
		mockClient.EXPECT().Execute(ctx, command).Return(reply, nil)
		ccUc.cHandlerClient = mockClient

		// SUT
		newId, err := ccUc.doCreate(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(newId).ToNot(Equal(uuid.Nil))
	})

})
