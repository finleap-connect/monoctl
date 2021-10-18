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
	"io"
	"strings"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// getTenantUsersUseCase provides the internal use-case of getting a tenant users model.
type getTenantUsersUseCase struct {
	useCaseBase
	tenantName   string
	tableFactory *output.TableFactory
}

func NewGetTenantUsersUseCase(config *config.Config, tenantName string, outputOptions *output.OutputOptions) UseCase {
	useCase := &getTenantUsersUseCase{
		useCaseBase: NewUseCaseBase("get-tenant-users", config),
		tenantName:  tenantName,
	}
	useCase.tableFactory = output.NewTableFactory().
		SetHeader([]string{"NAME", "EMAIL", "ROLES"}).
		SetSortColumn(outputOptions.SortOptions.SortByColumn).
		SetSortOrder(outputOptions.SortOptions.Order)

	return useCase
}

func (u *getTenantUsersUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	grpcClient := api.NewTenantClient(conn)

	tenant, err := grpcClient.GetByName(ctx, wrapperspb.String(u.tenantName))
	if err != nil {
		return err
	}

	tenantUserStream, err := grpcClient.GetUsers(ctx, wrapperspb.String(tenant.Id))
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		tenantUser, err := tenantUserStream.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		data = append(data, []interface{}{
			tenantUser.Name,
			tenantUser.Email,
			strings.Join(tenantUser.TenantRoles, ","),
		})
	}
	u.tableFactory.SetData(data) // Add Bulk Data
	u.tableFactory.ToTable().Render()

	return nil
}
