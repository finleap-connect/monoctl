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
	"io"
	"os"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/k8s"
	mdomain "github.com/finleap-connect/monoctl/test/mock/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateKubeconfig", func() {
	var (
		mockCtrl    *gomock.Controller
		m8TmpFile   *os.File
		kubeTmpFile *os.File
	)

	err := os.Setenv("KUBECONFIG", "")
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		var err error
		m8TmpFile, err = os.CreateTemp("", "monoskope")
		Expect(err).ToNot(HaveOccurred())
		kubeTmpFile, err = os.CreateTemp("", "kubeconfig")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		mockCtrl.Finish()

		os.Remove(m8TmpFile.Name())
		os.Remove(kubeTmpFile.Name())
	})

	var (
		ctx                         = context.Background()
		expectedId                  = uuid.New()
		expectedName                = "test-cluster"
		expectedApiServerAddress    = "test.cluster.monokope.io"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedKubeClusterName     = expectedName
		expectedKubeContextName     = "test-cluster-user"
		expectedAuthInfoName        = "test-cluster-jane-doe-user"
		expectedNamespaceName       = "jane-doe"
		expectedServer              = "m8.example.com"
	)

	newConfig := func() *config.Config {
		conf := config.NewConfig()
		conf.Server = expectedServer
		conf.AuthInformation = &config.AuthInformation{
			Token:    "this-is-a-token",
			Username: "jane.doe",
		}
		return conf
	}

	It("should run", func() {
		conf := newConfig()
		configManager := config.NewLoaderFromExplicitFile(m8TmpFile.Name())
		err = configManager.SaveToFile(conf, m8TmpFile.Name(), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = configManager.LoadConfig()
		Expect(err).ToNot(HaveOccurred())

		mockClusterAccessClient := mdomain.NewMockClusterAccessClient(mockCtrl)

		uc := NewUpdateKubeconfigUseCase(configManager, "", true).(*UpdateKubeconfigUseCase)
		uc.clusterAccessClient = mockClusterAccessClient
		uc.kubeConfig = k8s.NewKubeConfig()
		uc.kubeConfig.SetPath(kubeTmpFile.Name())
		uc.setInitialized()

		getClusterAccessClient := mdomain.NewMockClusterAccess_GetClusterAccessV2Client(mockCtrl)
		getClusterAccessClient.EXPECT().Recv().Return(&projections.ClusterAccessV2{
			Cluster: &projections.Cluster{
				Id:               expectedId.String(),
				Name:             expectedName,
				ApiServerAddress: expectedApiServerAddress,
				CaCertBundle:     expectedClusterCACertBundle,
			},
			ClusterRoles: []*projections.ClusterRole{{Scope: projections.ClusterRole_CLUSTER, Role: string(roles.User)}},
		}, nil)
		getClusterAccessClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClusterAccessClient.EXPECT().GetClusterAccessV2(ctx, &empty.Empty{}).Return(getClusterAccessClient, nil)
		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())

		kubeConfig, err := uc.kubeConfig.LoadConfig()
		Expect(err).ToNot(HaveOccurred())

		cluster, ok := kubeConfig.Clusters[expectedKubeClusterName]
		Expect(ok).To(BeTrue())
		Expect(cluster).NotTo(BeNil())
		Expect(cluster.Server).To(Equal(expectedApiServerAddress))
		Expect(cluster.CertificateAuthorityData).To(Equal(expectedClusterCACertBundle))

		kubeContext, ok := kubeConfig.Contexts[expectedKubeContextName]
		Expect(ok).To(BeTrue())
		Expect(kubeContext).NotTo(BeNil())
		Expect(kubeContext.Namespace).To(Equal(expectedNamespaceName))
		Expect(kubeContext.Cluster).To(Equal(expectedKubeClusterName))
		Expect(kubeContext.AuthInfo).To(Equal(expectedAuthInfoName))

		authInfo, ok := kubeConfig.AuthInfos[expectedAuthInfoName]
		Expect(ok).To(BeTrue())
		Expect(authInfo).NotTo(BeNil())
		Expect(authInfo.Exec).NotTo(BeNil())
		Expect(authInfo.Exec.Command).To(Equal("monoctl"))
	})
	It("should use kubeconfig file defined in m8 config", func() {
		conf := newConfig()
		conf.KubeConfigPath = kubeTmpFile.Name()
		configManager := config.NewLoaderFromExplicitFile(m8TmpFile.Name())
		err = configManager.SaveToFile(conf, m8TmpFile.Name(), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = configManager.LoadConfig()
		Expect(err).ToNot(HaveOccurred())

		mockClusterAccessClient := mdomain.NewMockClusterAccessClient(mockCtrl)

		uc := NewUpdateKubeconfigUseCase(configManager, "", true).(*UpdateKubeconfigUseCase)
		uc.clusterAccessClient = mockClusterAccessClient
		uc.kubeConfig = k8s.NewKubeConfig()
		uc.setInitialized()

		getClusterAccessClient := mdomain.NewMockClusterAccess_GetClusterAccessV2Client(mockCtrl)
		getClusterAccessClient.EXPECT().Recv().Return(nil, io.EOF)
		mockClusterAccessClient.EXPECT().GetClusterAccessV2(ctx, &empty.Empty{}).Return(getClusterAccessClient, nil)

		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(uc.kubeConfig.ConfigPath).To(Equal(kubeTmpFile.Name()))
	})
	It("should use kubeconfig file specified by the user", func() {
		conf := newConfig()
		conf.KubeConfigPath = "old/file/to/ignore"

		configManager := config.NewLoaderFromExplicitFile(m8TmpFile.Name())
		err = configManager.SaveToFile(conf, m8TmpFile.Name(), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = configManager.LoadConfig()
		Expect(err).ToNot(HaveOccurred())

		mockClusterAccessClient := mdomain.NewMockClusterAccessClient(mockCtrl)

		uc := NewUpdateKubeconfigUseCase(configManager, kubeTmpFile.Name(), true).(*UpdateKubeconfigUseCase)
		uc.clusterAccessClient = mockClusterAccessClient
		uc.kubeConfig = k8s.NewKubeConfig()
		uc.setInitialized()

		getClusterAccessClient := mdomain.NewMockClusterAccess_GetClusterAccessV2Client(mockCtrl)
		getClusterAccessClient.EXPECT().Recv().Return(nil, io.EOF)
		mockClusterAccessClient.EXPECT().GetClusterAccessV2(ctx, &empty.Empty{}).Return(getClusterAccessClient, nil)

		err = uc.Run(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = configManager.LoadConfig()
		Expect(err).ToNot(HaveOccurred())
		config := configManager.GetConfig()
		Expect(config.KubeConfigPath).To(Equal(kubeTmpFile.Name()))
	})
})
