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
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	api_domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// getRoleBindingsUseCase provides the internal use-case of getting the permission model.
type getRoleBindingsUseCase struct {
	useCaseBase
	email         string
	showDeleted   bool
	tableFactory  *output.TableFactory
	outputOptions *output.OutputOptions
}

func NewGetRoleBindingsUseCase(config *config.Config, email string, outputOptions *output.OutputOptions) UseCase {
	useCase := &getRoleBindingsUseCase{
		useCaseBase:   NewUseCaseBase("get-rolebindings", config),
		email:         email,
		outputOptions: outputOptions,
	}

	var header []string
	if outputOptions.Wide {
		header = append(header, "ID")
	}
	header = append(header, []string{"ROLE", "SCOPE", "RESOURCE", "AGE"}...)
	if outputOptions.ShowDeleted {
		header = append(header, "DELETED")
	}

	useCase.tableFactory = output.NewTableFactory().
		SetHeader(header).
		SetColumnFormatter("AGE", output.DefaultAgeColumnFormatter()).
		SetColumnFormatter("DELETED", output.DefaultAgeColumnFormatter()).
		SetSortColumn(outputOptions.SortOptions.SortByColumn).
		SetSortOrder(outputOptions.SortOptions.Order).
		SetExportFormat(outputOptions.ExportOptions.Format).
		SetExportFile(outputOptions.ExportOptions.File)

	return useCase
}

func (u *getRoleBindingsUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	users := api_domain.NewUserClient(conn)
	tenants := api_domain.NewTenantClient(conn)

	user, err := users.GetByEmail(ctx, wrapperspb.String(u.email))
	if err != nil {
		return err
	}

	roleBindingsStream, err := users.GetRoleBindingsById(ctx, wrapperspb.String(user.Id))
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		rb, err := roleBindingsStream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		resource := rb.Resource
		switch rb.Scope {
		case string(scopes.System):
		case string(scopes.Tenant):
			tenant, err := tenants.GetById(ctx, wrapperspb.String(rb.Resource))
			if err != nil {
				return err
			}
			resource = tenant.Name
		}

		var row []interface{}
		if u.outputOptions.Wide {
			row = append(row, rb.Id)
		}
		row = append(row, []interface{}{
			rb.Role,
			rb.Scope,
			resource,
			time.Since(rb.GetMetadata().GetCreated().AsTime()),
		}...)
		if u.showDeleted && rb.GetMetadata().GetDeleted().AsTime().Unix() != 0 {
			row = append(row, time.Since(rb.GetMetadata().GetDeleted().AsTime()))
		}
		data = append(data, row)
	}

	u.tableFactory.SetData(data) // Add Bulk Data
	tbl, err := u.tableFactory.ToTable()
	if err != nil {
		return err
	}

	tbl.Render()
	return nil
}
