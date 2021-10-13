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

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/output"
	api_commandhandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetPermissionModelUseCase provides the internal use-case of getting the permission model.
type getPoliciesUseCase struct {
	useCaseBase
	tableFactory *output.TableFactory
}

func NewGetPoliciesUseCase(config *config.Config, outputOptions *output.OutputOptions) UseCase {
	useCase := &getPoliciesUseCase{
		useCaseBase: NewUseCaseBase("get-policies", config),
	}

	useCase.tableFactory = output.NewTableFactory().
		SetHeader([]string{"COMMAND", "ROLE", "SCOPE"}).
		SetSortColumn(outputOptions.SortOptions.SortByColumn).
		SetSortOrder(outputOptions.SortOptions.Order)

	return useCase
}

func (u *getPoliciesUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	grpcClient := api_commandhandler.NewCommandHandlerExtensionsClient(conn)

	u.log.Info("Getting permission model from server...")
	policyOverview, err := grpcClient.GetPolicyOverview(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	var data [][]interface{}
	for _, p := range policyOverview.Policies {
		data = append(data, []interface{}{
			p.GetCommand(),
			p.GetRole(),
			p.GetScope(),
		})
	}

	u.tableFactory.SetData(data) // Add Bulk Data
	u.tableFactory.ToTable().Render()

	return nil
}
