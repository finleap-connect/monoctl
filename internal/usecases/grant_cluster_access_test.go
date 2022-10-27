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
	mEventSourcing "github.com/finleap-connect/monoctl/test/mock/eventsourcing"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zalando/go-keyring"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("grant cluster-access", func() {
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
		ctx                 = context.Background()
		expectedClusterId   = uuid.New()
		expectedClusterName = "test-cluster"
		expectedTenantId    = uuid.New()
		expectedTenantName  = "test-tenant"
		expectedExpiry      = time.Now().UTC().Add(1 * time.Hour)
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

		conf := confManager.GetConfig()
		conf.AuthInformation = &config.AuthInformation{
			Username: "test-user",
			Expiry:   expectedExpiry,
		}

		mockClusterClient := mdomain.NewMockClusterClient(mockCtrl)
		mockTenantClient := mdomain.NewMockTenantClient(mockCtrl)
		mockCmdHandlerClient := mEventSourcing.NewMockCommandHandlerClient(mockCtrl)

		uc := NewGrantClusterAccessUseCase(conf, expectedTenantName, expectedClusterName).(*grantClusterAccessUseCase)
		uc.clusterClient = mockClusterClient
		uc.tenantClient = mockTenantClient
		uc.cmdHandlerClient = mockCmdHandlerClient
		uc.setInitialized()

		mockClusterClient.EXPECT().GetByName(ctx, wrapperspb.String(expectedClusterName)).Return(&projections.Cluster{
			Id:   expectedClusterId.String(),
			Name: expectedClusterName,
		}, nil)

		mockTenantClient.EXPECT().GetByName(ctx, wrapperspb.String(expectedTenantName)).Return(&projections.Tenant{
			Id:   expectedTenantId.String(),
			Name: expectedTenantName,
		}, nil)

		command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateTenantClusterBinding, &cmdData.CreateTenantClusterBindingCommandData{
			TenantId:  expectedTenantId.String(),
			ClusterId: expectedClusterId.String(),
		})
		Expect(err).NotTo(HaveOccurred())

		mockCmdHandlerClient.EXPECT().Execute(ctx, command)

		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

})
