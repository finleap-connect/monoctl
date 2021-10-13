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
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type createRoleBindingUseCase struct {
	useCaseBase
	mailaddress string
	role        string
	scope       string
	resource    string
}

func NewCreateRoleBindingUseCase(config *config.Config, mailaddress, role, scope, resource string) UseCase {
	useCase := &createRoleBindingUseCase{
		useCaseBase: NewUseCaseBase("create-rolebinding", config),
		mailaddress: mailaddress,
		role:        role,
		scope:       scope,
		resource:    resource,
	}
	return useCase
}

func (u *createRoleBindingUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	cmdHandlerExtensionClient := api.NewCommandHandlerExtensionsClient(conn)

	u.log.Info("Getting permission model from server...")
	permissionModel, err := cmdHandlerExtensionClient.GetPermissionModel(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	u.log.Info("Validating if role and scope are valid...")
	if !contains(permissionModel.Roles, u.role) {
		return fmt.Errorf("role '%s' does not exist", u.role)
	}
	if !contains(permissionModel.Scopes, u.scope) {
		return fmt.Errorf("scope '%s' does not exist", u.scope)
	}

	switch u.scope {
	case scopes.System.String():
	case scopes.Tenant.String():
		grpcClient := api.NewTenantClient(conn)
		tenant, err := grpcClient.GetByName(ctx, wrapperspb.String(u.resource))
		if err != nil {
			return err
		}
		u.resource = tenant.Id
	default:
		return fmt.Errorf("scope '%s' is not implemented", u.scope)
	}

	var resourceId uuid.UUID
	if resourceId, err = uuid.Parse(u.resource); len(u.resource) > 0 && err != nil {
		return err
	}

	userServiceClient := api.NewUserClient(conn)
	user, err := userServiceClient.GetByEmail(ctx, wrapperspb.String(u.mailaddress))
	if err != nil {
		return err
	}

	command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding)
	if _, err = cmd.AddCommandData(command,
		&cmdData.CreateUserRoleBindingCommandData{
			UserId:   user.Id,
			Role:     u.role,
			Scope:    u.scope,
			Resource: resourceId.String(),
		},
	); err != nil {
		return err
	}

	client := esApi.NewCommandHandlerClient(conn)
	_, err = client.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Println("Rolebinding created.")
	}

	return err
}

func contains(values []string, item string) bool {
	for _, value := range values {
		if value == item {
			return true
		}
	}
	return false
}
