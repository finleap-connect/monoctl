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
	"fmt"
	"io"
	"strings"

	"github.com/finleap-connect/monoctl/internal/config"
	mgrpc "github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/k8s"
	"github.com/finleap-connect/monoctl/internal/spinner"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	mk8s "github.com/finleap-connect/monoskope/pkg/k8s"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	kapi "k8s.io/client-go/tools/clientcmd/api"
)

type createKubeConfigUseCase struct {
	useCaseBase
	conn                *ggrpc.ClientConn
	clusterAccessClient api.ClusterAccessClient
	userClient          api.UserClient
	kubeConfig          *k8s.KubeConfig
}

func NewCreateKubeConfigUseCase(config *config.Config) UseCase {
	useCase := &createKubeConfigUseCase{
		useCaseBase: NewUseCaseBase("create-kubeconfig", config),
	}
	return useCase
}

func (u *createKubeConfigUseCase) init(ctx context.Context) error {
	if u.initialized {
		return nil
	}

	conn, err := mgrpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.clusterAccessClient = api.NewClusterAccessClient(u.conn)
	u.userClient = api.NewUserClient(u.conn)

	u.kubeConfig = k8s.NewKubeConfig()
	u.setInitialized()

	return nil
}

func (u *createKubeConfigUseCase) getNaming(m8ClusterName string, clusterRole mk8s.K8sRole) (clusterName, contextName, nsName, authInfoName string, err error) {
	nsName, err = mk8s.GetNamespaceName(strings.Replace(u.config.AuthInformation.Username, " ", "-", -1))
	if err != nil {
		return
	}

	clusterName, err = mk8s.GetNamespaceName(m8ClusterName)
	if err != nil {
		return
	}

	authInfoName = fmt.Sprintf("%s-%s-%s", clusterName, nsName, clusterRole)
	contextName = fmt.Sprintf("%s-%s", clusterName, clusterRole)

	u.log.Info("Naming configured for cluster.", "cluster", clusterName, "context", contextName, "ns", nsName, "authinfo", authInfoName)

	return
}

// setContext sets the context the given on kubeconfig
func (u *createKubeConfigUseCase) setContext(kubeConfig *kapi.Config, clusterName, contextName, nsName, authInfoName string) {
	var ok bool
	var kubeContext *kapi.Context
	if kubeContext, ok = kubeConfig.Contexts[contextName]; !ok {
		kubeContext = kapi.NewContext()
		kubeConfig.Contexts[contextName] = kubeContext
	}
	kubeContext.Namespace = nsName
	kubeContext.Cluster = clusterName
	kubeContext.AuthInfo = authInfoName

	u.log.Info("Context created/updated.", "context", contextName)
}

// setCluster sets the cluster configuration the given on kubeconfig
func (u *createKubeConfigUseCase) setCluster(kubeConfig *kapi.Config, m8Cluster *projections.Cluster, clusterName string) {
	var ok bool
	var cluster *kapi.Cluster
	if cluster, ok = kubeConfig.Clusters[clusterName]; !ok {
		cluster = kapi.NewCluster()
		kubeConfig.Clusters[clusterName] = cluster
	}

	cluster.CertificateAuthorityData = m8Cluster.CaCertBundle
	cluster.CertificateAuthority = "" // clear other authority data which clashes
	cluster.Server = m8Cluster.ApiServerAddress

	u.log.Info("Cluster created/updated.", "cluster", clusterName)
}

// setAuthInfo sets the auth information on kubeconfig
func (u *createKubeConfigUseCase) setAuthInfo(kubeConfig *kapi.Config, authInfoName, clusterName string, clusterRole mk8s.K8sRole) {
	var ok bool
	var kubeAuthInfo *kapi.AuthInfo
	if kubeAuthInfo, ok = kubeConfig.AuthInfos[authInfoName]; !ok {
		kubeAuthInfo = kapi.NewAuthInfo()
		kubeConfig.AuthInfos[authInfoName] = kubeAuthInfo
	}
	kubeAuthInfo.Exec = &kapi.ExecConfig{
		APIVersion:  "client.authentication.k8s.io/v1beta1",
		InstallHint: "Monoskope's commandline tool `monoctl` is required to authenticate to the current cluster.",
		Command:     "monoctl",
		Args: []string{
			"get", "cluster-credentials", clusterName, string(clusterRole),
		},
		Env: make([]kapi.ExecEnvVar, 0),
	}
	u.log.Info("AuthInfo created/updated.", "authinfo", authInfoName)
}

func (u *createKubeConfigUseCase) run(ctx context.Context) error {
	var err error

	// Load kubeconfig of current user
	var kubeConfig *kapi.Config
	if kubeConfig, err = u.kubeConfig.LoadConfig(); err != nil {
		return err
	}

	// Get user information of the current user
	user, err := u.userClient.GetByEmail(ctx, wrapperspb.String(u.config.AuthInformation.Username))
	if err != nil {
		return err
	}

	// Get cluster information from control plane
	clusterAccesses, err := u.clusterAccessClient.GetClusterAccessByUserId(ctx, wrapperspb.String(user.Id))
	if err != nil {
		return err
	}

	for {
		// Read next
		clusterAccess, err := clusterAccesses.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		for _, clusterRole := range mk8s.AvailableRoles {
			// Get naming
			clusterName, contextName, nsName, authInfoName, err := u.getNaming(clusterAccess.Name, clusterRole)
			if err != nil {
				return err
			}

			// Set cluster on kubeconfig
			u.setCluster(kubeConfig, clusterAccess, clusterName)

			// Set context on kubeconfig
			u.setContext(kubeConfig, clusterName, contextName, nsName, authInfoName)

			// Set credentials on kubeconfig
			u.setAuthInfo(kubeConfig, authInfoName, clusterName, clusterRole)
		}
	}

	return u.kubeConfig.StoreConfig(kubeConfig)
}

func (u *createKubeConfigUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	err := u.init(ctx)
	if err != nil {
		return err
	}
	if u.conn != nil {
		defer u.conn.Close()
	}

	err = u.run(ctx)
	if err != nil {
		return err
	}
	s.Stop()

	fmt.Println("Your kubeconfig has been generated/updated.")
	fmt.Println("Use `kubectl config get-contexts` to see available contexts.")
	fmt.Println("Use `kubectl config use-context <CONTEXTNAME>` to switch between clusters.")

	return nil
}
