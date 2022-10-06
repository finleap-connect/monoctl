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
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type createRoleBindingUseCase struct {
	useCaseBase
	mailAddresses []string
	role          string
	scope         string
	resource      string
}

func NewCreateRoleBindingUseCase(config *config.Config, mailAddresses []string, role, scope, resource string) UseCase {
	useCase := &createRoleBindingUseCase{
		useCaseBase:   NewUseCaseBase("create-rolebinding", config),
		mailAddresses: mailAddresses,
		role:          role,
		scope:         scope,
		resource:      resource,
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
	case string(scopes.System):
	case string(scopes.Tenant):
		grpcClient := api.NewTenantClient(conn)
		tenant, err := grpcClient.GetByName(ctx, wrapperspb.String(u.resource))
		if err != nil {
			return err
		}
		u.resource = tenant.Id
	default:
		return fmt.Errorf("scope '%s' is not implemented", u.scope)
	}

	if len(u.resource) > 0 {
		if _, err = uuid.Parse(u.resource); err != nil {
			return err
		}
	}

	userServiceClient := api.NewUserClient(conn)

	for _, mailAddress := range u.mailAddresses {
		user, err := userServiceClient.GetByEmail(ctx, wrapperspb.String(mailAddress))
		if err != nil {
			return err
		}
		command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateUserRoleBinding, &cmdData.CreateUserRoleBindingCommandData{
			UserId:   user.Id,
			Role:     u.role,
			Scope:    u.scope,
			Resource: wrapperspb.String(u.resource),
		})

		client := esApi.NewCommandHandlerClient(conn)
		_, err = client.Execute(ctx, command)
		if err != nil {
			return err
		}

		fmt.Printf("Role binding created for user '%s' with role '%s' in scope '%s'.", mailAddress, u.role, u.scope)
	}

	s.Stop()

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
