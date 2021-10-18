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
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"golang.org/x/oauth2"
)

type deleteUserRoleBindingUseCase struct {
	useCaseBase
	id string
}

func NewDeleteUserRoleBindingUseCase(config *config.Config, id string) UseCase {
	useCase := &deleteUserRoleBindingUseCase{
		useCaseBase: NewUseCaseBase("delete-userrolebinding", config),
		id:          id,
	}
	return useCase
}

func (u *deleteUserRoleBindingUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()

	// Validate id is valid uuid
	id, err := uuid.Parse(u.id)
	if err != nil {
		return err
	}

	command := cmd.CreateCommand(id, commandTypes.DeleteUserRoleBinding)
	cmdHandlerClient := esApi.NewCommandHandlerClient(conn)
	_, err = cmdHandlerClient.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Println("UserRoleBinding deleted.")
	}

	return err
}
