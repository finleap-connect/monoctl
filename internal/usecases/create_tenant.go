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

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/spinner"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"golang.org/x/oauth2"
)

type createTenantUseCase struct {
	useCaseBase
	name   string
	prefix string
}

func NewCreateTenantUseCase(config *config.Config, name, tenant string) UseCase {
	useCase := &createTenantUseCase{
		useCaseBase: NewUseCaseBase("create-tenant", config),
		name:        name,
		prefix:      tenant,
	}
	return useCase
}

func (u *createTenantUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()

	command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant)
	if _, err := cmd.AddCommandData(command, &cmdData.CreateTenantCommandData{
		Name:   u.name,
		Prefix: u.prefix,
	}); err != nil {
		return err
	}

	client := esApi.NewCommandHandlerClient(conn)
	_, err = client.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Printf("Tenant '%s' created.\n", u.name)
	}

	return err
}
