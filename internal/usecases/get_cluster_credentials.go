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
	"time"

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	mgrpc "gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	apiGateway "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclientauth "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

type getClusterCredentialsUseCase struct {
	useCaseBase
	conn                 *ggrpc.ClientConn
	configManager        *config.ClientConfigManager
	clusterServiceClient api.ClusterClient
	clusterAuthClient    apiGateway.ClusterAuthClient
	clusterName          string
	clusterRole          string
}

func NewGetClusterCredentialsUseCase(configManager *config.ClientConfigManager, name, role string) UseCase {
	useCase := &getClusterCredentialsUseCase{
		useCaseBase:   NewUseCaseBase("get-cluster-credentials", configManager.GetConfig()),
		configManager: configManager,
		clusterName:   name,
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
	var token string
	var expiry time.Time

	clusterAuthInfo := u.config.GetClusterAuthInformation(u.clusterName, u.config.AuthInformation.Username, u.clusterRole)
	if clusterAuthInfo != nil && clusterAuthInfo.IsValidExact() {
		token = clusterAuthInfo.Token
		expiry = clusterAuthInfo.Expiry
	} else {
		// Query cluster
		m8cluster, err := u.clusterServiceClient.GetByName(ctx, wrapperspb.String(u.clusterName))
		if err != nil {
			return err
		}

		// Get token from gateway
		response, err := u.clusterAuthClient.GetAuthToken(ctx, &apiGateway.ClusterAuthTokenRequest{
			ClusterId: m8cluster.Id,
			Role:      u.clusterRole,
		})
		if err != nil {
			return err
		}
		token = response.GetAccessToken()
		expiry = response.Expiry.AsTime()

		// Cache credentials
		u.config.SetClusterAuthInformation(u.clusterName, u.config.AuthInformation.Username, u.clusterRole, token, expiry)
		err = u.configManager.SaveConfig()
		if err != nil {
			return err
		}
	}

	// Convert to kubectl readable format
	execCredential := kclientauth.ExecCredential{
		TypeMeta: v1.TypeMeta{
			Kind:       "ExecCredential",
			APIVersion: "client.authentication.k8s.io/v1beta1",
		},
		Status: &kclientauth.ExecCredentialStatus{
			Token: token,
			ExpirationTimestamp: &v1.Time{
				Time: expiry,
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
