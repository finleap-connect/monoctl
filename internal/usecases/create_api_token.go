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
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	mgrpc "github.com/finleap-connect/monoctl/internal/grpc"
	apiGateway "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/google/uuid"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

type createAPITokenUsecase struct {
	useCaseBase
	conn           *ggrpc.ClientConn
	configManager  *config.ClientConfigManager
	apiTokenClient apiGateway.APITokenClient
	userId         string
	scopes         []string
	validity       time.Duration
}

func NewCreateAPITokenUsecase(configManager *config.ClientConfigManager, userId string, scopes []string, validity time.Duration) UseCase {
	useCase := &createAPITokenUsecase{
		useCaseBase:   NewUseCaseBase("create-api-token", configManager.GetConfig()),
		configManager: configManager,
		userId:        userId,
		scopes:        scopes,
		validity:      validity,
	}
	return useCase
}

func (u *createAPITokenUsecase) init(ctx context.Context) error {
	if u.initialized {
		return nil
	}

	conn, err := mgrpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.apiTokenClient = apiGateway.NewAPITokenClient(u.conn)

	u.setInitialized()

	return nil
}

func (u *createAPITokenUsecase) run(ctx context.Context) error {
	var authScopes []apiGateway.AuthorizationScope
	for _, scope := range u.scopes {
		authScopes = append(authScopes, apiGateway.AuthorizationScope(apiGateway.AuthorizationScope_value[scope]))
	}

	request := &apiGateway.APITokenRequest{
		AuthorizationScopes: authScopes,
		Validity:            durationpb.New(u.validity),
	}

	if _, err := uuid.Parse(u.userId); err != nil {
		request.User = &apiGateway.APITokenRequest_UserId{
			UserId: u.userId,
		}
	} else {
		request.User = &apiGateway.APITokenRequest_Username{
			Username: u.userId,
		}
	}

	// Get token from gateway
	response, err := u.apiTokenClient.RequestAPIToken(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(response.GetAccessToken())

	return nil
}

func (u *createAPITokenUsecase) Run(ctx context.Context) error {
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
