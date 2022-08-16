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
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
)

// getTenantsUseCase provides the internal use-case of getting the permission model.
type getTenantsUseCase struct {
	useCaseBase
	conn          *ggrpc.ClientConn
	client        api.TenantClient
	tableFactory  *output.TableFactory
	outputOptions *output.OutputOptions
}

func NewGetTenantsUseCase(config *config.Config, outputOptions *output.OutputOptions) UseCase {
	useCase := &getTenantsUseCase{
		useCaseBase:   NewUseCaseBase("get-tenants", config),
		outputOptions: outputOptions,
	}

	var header []string
	if outputOptions.Wide {
		header = append(header, "ID")
	}
	header = append(header, []string{"NAME", "PREFIX", "AGE"}...)
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

func (u *getTenantsUseCase) Run(ctx context.Context) error {
	err := u.setUp(ctx)
	if err != nil {
		return err
	}
	defer u.conn.Close()

	err = u.doRun(ctx)
	if err != nil {
		return err
	}

	tbl, err := u.tableFactory.ToTable()
	if err != nil {
		return err
	}

	tbl.Render()
	return nil
}

func (u *getTenantsUseCase) setUp(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	u.conn = conn
	u.client = api.NewTenantClient(conn)

	return nil
}

func (u *getTenantsUseCase) doRun(ctx context.Context) error {
	tenantStream, err := u.client.GetAll(ctx, &api.GetAllRequest{
		IncludeDeleted: u.outputOptions.ShowDeleted,
	})
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		tenant, err := tenantStream.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		var row []interface{}
		if u.outputOptions.Wide {
			row = append(row, tenant.Id)
		}
		row = append(row, []interface{}{
			tenant.Name,
			tenant.Prefix,
			time.Since(tenant.Metadata.Created.AsTime()),
		}...)
		if u.outputOptions.ShowDeleted && tenant.Metadata.Deleted.AsTime().Unix() != 0 {
			row = append(row, time.Since(tenant.Metadata.Deleted.AsTime()))
		}
		data = append(data, row)
	}
	u.tableFactory.SetData(data) // Add Bulk Data

	return nil
}
