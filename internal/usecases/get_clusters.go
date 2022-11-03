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
	ggrpc "google.golang.org/grpc"
)

// getClustersUseCase provides the internal use-case of getting the permission model.
type getClustersUseCase struct {
	useCaseBase
	conn          *ggrpc.ClientConn
	client        api_domain.ClusterClient
	tableFactory  *output.TableFactory
	outputOptions *output.OutputOptions
}

func NewGetClustersUseCase(config *config.Config, outputOptions *output.OutputOptions) UseCase {
	useCase := &getClustersUseCase{
		useCaseBase:   NewUseCaseBase("get-clusters", config),
		outputOptions: outputOptions,
	}

	var header []string
	if outputOptions.Wide {
		header = append(header, "ID")
	}
	header = append(header, []string{"NAME", "API SERVER ADDRESS", "AGE"}...)
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

func (u *getClustersUseCase) setUp(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.client = api_domain.NewClusterClient(u.conn)

	return nil
}

func (u *getClustersUseCase) doRun(ctx context.Context) error {
	clusterStream, err := u.client.GetAll(ctx, &api_domain.GetAllRequest{
		IncludeDeleted: u.outputOptions.ShowDeleted,
	})
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		cluster, err := clusterStream.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		var row []interface{}
		if u.outputOptions.Wide {
			row = append(row, cluster.Id)
		}
		row = append(row, []interface{}{
			cluster.Name,
			cluster.ApiServerAddress,
			time.Since(cluster.Metadata.Created.AsTime()),
		}...)
		if u.outputOptions.ShowDeleted && cluster.Metadata.Deleted.AsTime().Unix() != 0 {
			row = append(row, time.Since(cluster.Metadata.Deleted.AsTime()))
		}
		data = append(data, row)
	}
	u.tableFactory.SetData(data)

	return nil

}
func (u *getClustersUseCase) Run(ctx context.Context) error {
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
