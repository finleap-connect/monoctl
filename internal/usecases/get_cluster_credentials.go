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
	"encoding/json"
	"fmt"

	"github.com/finleap-connect/monoctl/internal/config"
	mgrpc "github.com/finleap-connect/monoctl/internal/grpc"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	apiGateway "github.com/finleap-connect/monoskope/pkg/api/gateway"
	ggrpc "google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclientauth "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

type getClusterCredentialsUseCase struct {
	useCaseBase
	conn                 *ggrpc.ClientConn
	configManager        *config.ClientConfigManager
	clusterServiceClient api.ClusterClient
	clusterAuthClient    apiGateway.ClusterAuthClient
	clusterId            string
	clusterRole          string
}

func NewGetClusterCredentialsUseCase(configManager *config.ClientConfigManager, clusterId, role string) UseCase {
	useCase := &getClusterCredentialsUseCase{
		useCaseBase:   NewUseCaseBase("get-cluster-credentials", configManager.GetConfig()),
		configManager: configManager,
		clusterId:     clusterId,
		clusterRole:   role,
	}
	return useCase
}

func (u *getClusterCredentialsUseCase) init(ctx context.Context) error {
	if u.initialized {
		return nil
	}

	conn, err := mgrpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.clusterServiceClient = api.NewClusterClient(u.conn)
	u.clusterAuthClient = apiGateway.NewClusterAuthClient(u.conn)

	u.setInitialized()

	return nil
}

func (u *getClusterCredentialsUseCase) run(ctx context.Context) error {
	clusterAuthInfo := u.config.GetClusterAuthInformation(u.clusterId, u.config.AuthInformation.Username, u.clusterRole)
	if clusterAuthInfo == nil || !clusterAuthInfo.IsValidExact() {
		// Get cluster credentials
		_, err := u.requestClusterAuthInformation(ctx, u.clusterId)
		if err != nil {
			return err
		}
		clusterAuthInfo = u.config.GetClusterAuthInformation(u.clusterId, u.config.AuthInformation.Username, u.clusterRole)

	}

	// Convert to kubectl readable format
	execCredential := kclientauth.ExecCredential{
		TypeMeta: v1.TypeMeta{
			Kind:       "ExecCredential",
			APIVersion: "client.authentication.k8s.io/v1beta1",
		},
		Status: &kclientauth.ExecCredentialStatus{
			Token: clusterAuthInfo.Token,
			ExpirationTimestamp: &v1.Time{
				Time: clusterAuthInfo.Expiry,
			},
		},
	}

	// Marshal
	bytes, err := json.Marshal(execCredential)
	if err != nil {
		return err
	}

	// Write marshalled yaml to stdout
	fmt.Println(string(bytes))

	return nil
}

func (u *getClusterCredentialsUseCase) requestClusterAuthInformation(ctx context.Context, clusterId string) (response *apiGateway.ClusterAuthTokenResponse, err error) {
	// Get token from gateway
	response, err = u.clusterAuthClient.GetAuthToken(ctx, &apiGateway.ClusterAuthTokenRequest{
		ClusterId: clusterId,
		Role:      u.clusterRole,
	})
	if err != nil {
		return
	}

	// Get credentials
	u.config.SetClusterAuthInformation(clusterId, u.config.AuthInformation.Username, u.clusterRole, response.AccessToken, response.Expiry.AsTime())

	// Save credentials
	err = u.configManager.SaveConfig()
	if err != nil {
		return
	}

	return
}

func (u *getClusterCredentialsUseCase) Run(ctx context.Context) error {
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

	return nil
}
