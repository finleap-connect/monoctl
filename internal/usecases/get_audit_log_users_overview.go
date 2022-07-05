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
	"github.com/finleap-connect/monoctl/internal/config"
	m8Grpc "github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"time"
)

// getAuditLogUsersOverviewUseCase provides the internal use-case of getting the audit overview of all users.
type getAuditLogUsersOverviewUseCase struct {
	useCaseBase
	conn           *grpc.ClientConn
	tableFactory   *output.TableFactory
	outputOptions  *output.OutputOptions
	auditLogClient api.AuditLogClient
	timestamp      time.Time
}

func NewGetAuditLogUsersOverviewUseCase(config *config.Config, outputOptions *output.OutputOptions, timestamp time.Time) UseCase {
	useCase := &getAuditLogUsersOverviewUseCase{
		useCaseBase:   NewUseCaseBase("get-audit-log-users-overview", config),
		outputOptions: outputOptions,
		timestamp:     timestamp,
	}

	header := []string{"NAME", "EMAIL", "ROLES", "TENANTS", "CLUSTERS", "DETAILS"}

	useCase.tableFactory = output.NewTableFactory().
		SetHeader(header).
		SetSortColumn(outputOptions.SortOptions.SortByColumn).
		SetSortOrder(outputOptions.SortOptions.Order).
		SetExportFormat(outputOptions.ExportOptions.Format).
		SetExportFile(outputOptions.ExportOptions.File)

	return useCase
}

func (u *getAuditLogUsersOverviewUseCase) Run(ctx context.Context) error {
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

func (u *getAuditLogUsersOverviewUseCase) setUp(ctx context.Context) error {
	conn, err := m8Grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	u.conn = conn
	u.auditLogClient = api.NewAuditLogClient(conn)

	return nil
}

func (u *getAuditLogUsersOverviewUseCase) doRun(ctx context.Context) error {
	overviewStream, err := u.auditLogClient.GetUsersOverview(ctx, &api.GetUsersOverviewRequest{Timestamp: timestamppb.New(u.timestamp)})
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		overview, err := overviewStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		dataLine := []interface{}{
			overview.Name,
			overview.Email,
			overview.Roles,
			overview.Tenants,
			overview.Clusters,
			overview.Details,
		}
		data = append(data, dataLine)
	}
	u.tableFactory.SetData(data)

	return nil
}
