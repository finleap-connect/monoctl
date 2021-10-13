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

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/output"
	api_commandhandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	ggrpc "google.golang.org/grpc"
)

// getClustersUseCase provides the internal use-case of getting the permission model.
type getClustersUseCase struct {
	useCaseBase
	conn          *ggrpc.ClientConn
	client        api_commandhandler.ClusterClient
	tableFactory  *output.TableFactory
	outputOptions *output.OutputOptions
}

func NewGetClustersUseCase(config *config.Config, outputOptions *output.OutputOptions) UseCase {
	useCase := &getClustersUseCase{
		useCaseBase:   NewUseCaseBase("get-clusters", config),
		outputOptions: outputOptions,
	}

	header := []string{"NAME", "DISPLAY NAME", "DNS ADDRESS", "AGE"}
	if outputOptions.ShowDeleted {
		header = append(header, "DELETED")
	}

	useCase.tableFactory = output.NewTableFactory().
		SetHeader(header).
		SetColumnFormatter("AGE", output.DefaultAgeColumnFormatter()).
		SetColumnFormatter("DELETED", output.DefaultAgeColumnFormatter()).
		SetSortColumn(outputOptions.SortOptions.SortByColumn).
		SetSortOrder(outputOptions.SortOptions.Order)

	return useCase
}

func (u *getClustersUseCase) setUp(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.client = api_commandhandler.NewClusterClient(u.conn)

	return nil
}

func (u *getClustersUseCase) doRun(ctx context.Context) error {
	clusterStream, err := u.client.GetAll(ctx, &api_commandhandler.GetAllRequest{
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

		dataRow := []interface{}{
			cluster.Name,
			cluster.DisplayName,
			cluster.ApiServerAddress,
			time.Since(cluster.Metadata.Created.AsTime()),
		}
		if u.outputOptions.ShowDeleted && cluster.Metadata.Deleted.AsTime().Unix() != 0 {
			dataRow = append(dataRow, time.Since(cluster.Metadata.Deleted.AsTime()))
		}
		data = append(data, dataRow)
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
	u.tableFactory.ToTable().Render()

	return nil
}
