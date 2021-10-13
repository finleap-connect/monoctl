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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	mdom "gitlab.figo.systems/platform/monoskope/monoctl/test/mock/domain"
	mes "gitlab.figo.systems/platform/monoskope/monoctl/test/mock/eventsourcing"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// rendered output for certificate resource and issuer
//go:embed expected_m8_operator_bootstrap.yaml
var expectedResource string

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
		expectedDisplayName         = "the one cluster"
		expectedName                = "one-cluster"
		expectedApiServerAddress    = "one.example.com"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedUUID                = uuid.New()
		expectedJwt                 = "this-is-a-jwt.it-contains-chars-illegal-for-base64"
	)

	It("should construct a gRPC call", func() {
		var err error

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		ccUc := NewCreateClusterUseCase(conf, expectedDisplayName, expectedName,
			expectedApiServerAddress, expectedClusterCACertBundle).(*createClusterUseCase)

		Expect(ccUc.displayName).To(Equal(expectedDisplayName))
		Expect(ccUc.name).To(Equal(expectedName))
		Expect(ccUc.apiServerAddress).To(Equal(expectedApiServerAddress))
		Expect(ccUc.caCertBundle).To(Equal(expectedClusterCACertBundle))

		ctx := context.Background()

		// set up dummy connection
		ccUc.conn = grpc.CreateDummyGrpcConnection()

		// use mocked commandHandlerClient
		mockClient := mes.NewMockCommandHandlerClient(mockCtrl)
		commanddata := &cmdData.CreateCluster{
			DisplayName:      expectedDisplayName,
			Name:             expectedName,
			ApiServerAddress: expectedApiServerAddress,
			CaCertBundle:     expectedClusterCACertBundle,
		}

		command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster)
		_, err = cmd.AddCommandData(command, commanddata)
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

	It("should retrieve the jwt", func() {

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		ccUc := NewCreateClusterUseCase(conf, expectedDisplayName, expectedName,
			expectedApiServerAddress, expectedClusterCACertBundle).(*createClusterUseCase)

		ctx := context.Background()

		// use mocked commandHandlerClient
		mockClient := mdom.NewMockClusterClient(mockCtrl)

		mockClient.EXPECT().GetBootstrapToken(ctx, &wrapperspb.StringValue{Value: expectedUUID.String()}).Return(&wrapperspb.StringValue{Value: expectedJwt}, nil)
		ccUc.clusterClient = mockClient

		// SUT
		err := ccUc.queryJwt(ctx, expectedUUID.String())
		Expect(err).ToNot(HaveOccurred())

		Expect(ccUc.jwt).To(Equal(expectedJwt))
	})

	It("should render the certificate template correctly", func() {

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token: "this-is-a-token",
		}

		ccUc := NewCreateClusterUseCase(conf, expectedDisplayName, expectedName,
			expectedApiServerAddress, expectedClusterCACertBundle).(*createClusterUseCase)

		ccUc.jwt = expectedJwt

		actualResource, err := ccUc.renderOutput()
		Expect(err).ToNot(HaveOccurred())

		Expect(actualResource).To(Equal(expectedResource))
	})
})
