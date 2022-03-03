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
	"io/ioutil"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/k8s"
	mdomain "github.com/finleap-connect/monoctl/test/mock/domain"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateKubeconfig", func() {
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
		ctx                         = context.Background()
		expectedId                  = uuid.New()
		expectedDisplayName         = "testcluster"
		expectedName                = "test-cluster"
		expectedApiServerAddress    = "test.cluster.monokope.io"
		expectedClusterCACertBundle = []byte("This should be a certificate")
		expectedKubeClusterName     = expectedName
		expectedKubeContextName     = "test-cluster-default"
		expectedAuthInfoName        = "test-cluster-jane-doe-default"
		expectedNamespaceName       = "jane-doe"
	)

	It("should run", func() {
		var err error

		tmpfile, err := ioutil.TempFile("", "kubeconfig")
		Expect(err).ToNot(HaveOccurred())

		conf := config.NewConfig()
		conf.Server = "m8.example.com"
		conf.AuthInformation = &config.AuthInformation{
			Token:    "this-is-a-token",
			Username: "jane.doe",
		}

		mockClusterClient := mdomain.NewMockClusterAccessClient(mockCtrl)

		uc := NewCreateKubeConfigUseCase(conf).(*createKubeConfigUseCase)
		uc.clusterAccessClient = mockClusterClient
		uc.kubeConfig = k8s.NewKubeConfig()
		uc.kubeConfig.SetPath(tmpfile.Name())
		uc.setInitialized()

		getAllClient := mdomain.NewMockCluster_GetAllClient(mockCtrl)
		getAllClient.EXPECT().Recv().Return(&projections.Cluster{
			Id:               expectedId.String(),
			DisplayName:      expectedDisplayName,
			Name:             expectedName,
			ApiServerAddress: expectedApiServerAddress,
			CaCertBundle:     expectedClusterCACertBundle,
		}, nil)
		getAllClient.EXPECT().Recv().Return(nil, io.EOF)

		mockClusterClient.EXPECT().GetClusterAccessByUserId(ctx, &api.GetAllRequest{IncludeDeleted: false}).Return(getAllClient, nil)
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

})
