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

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type createUserUseCase struct {
	useCaseBase
	username    string
	mailaddress string
}

func NewCreateUserUseCase(config *config.Config, username, mailaddress string) UseCase {
	useCase := &createUserUseCase{
		useCaseBase: NewUseCaseBase("create-user", config),
		username:    username,
		mailaddress: mailaddress,
	}
	return useCase
}

func (u *createUserUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()

	command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateUser, &cmdData.CreateUserCommandData{
		Name:  u.username,
		Email: u.mailaddress,
	})

	client := eventsourcing.NewCommandHandlerClient(conn)
	_, err = client.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Printf("User '%s' with email address '%s' created.\n", u.username, u.mailaddress)
	}

	return err
}
