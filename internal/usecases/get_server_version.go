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

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
)

// getServerVersionUseCase provides the internal use-case of getting the server version.
type getServerVersionUseCase struct {
	useCaseBase
}

func NewGetServerVersionUseCase(config *config.Config) UseCase {
	useCase := &getServerVersionUseCase{
		useCaseBase: NewUseCaseBase("get-server-version", config),
	}
	return useCase
}

func (u *getServerVersionUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	grpcClient := api_common.NewServiceInformationServiceClient(conn)

	u.log.Info("Getting service information from server...")
	serverInfo, err := grpcClient.GetServiceInformation(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	var serviceInfos []string
	for {
		// Read next
		serverInfo, err := serverInfo.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		// Append
		serviceInfos = append(serviceInfos, serverInfo.GetName())
		serviceInfos = append(serviceInfos, fmt.Sprintf("  version     : %s", serverInfo.GetVersion()))
		serviceInfos = append(serviceInfos, fmt.Sprintf("  commit      : %s", serverInfo.GetCommit()))
	}

	for _, version := range serviceInfos {
		fmt.Print(version)
		fmt.Println()
	}

	return nil
}
