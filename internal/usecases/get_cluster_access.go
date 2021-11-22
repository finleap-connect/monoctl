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
	"errors"
	"io"
	"time"

	"github.com/finleap-connect/monoctl/internal/config"
	m8Grpc "github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/output"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type getClusterAccess struct {
	useCaseBase
	conn                *grpc.ClientConn
	tenantName          string
	clusterName         string
	tenantClient        api.TenantClient
	clusterClient       api.ClusterClient
	clusterAccessClient api.ClusterAccessClient
	outputOptions       *output.OutputOptions
}

func NewGetClusterAccessUseCase(config *config.Config, outputOptions *output.OutputOptions, tenantName, clusterName string) UseCase {
	useCase := &getClusterAccess{
		useCaseBase:   NewUseCaseBase("get-cluster-access", config),
		tenantName:    tenantName,
		clusterName:   clusterName,
		outputOptions: outputOptions,
	}

	return useCase
}

func (u *getClusterAccess) init(ctx context.Context) error {
	if u.initialized {
		return nil
	}

	conn, err := m8Grpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.clusterClient = api.NewClusterClient(u.conn)
	u.tenantClient = api.NewTenantClient(conn)
	u.clusterAccessClient = api.NewClusterAccessClient(u.conn)
	u.setInitialized()

	return nil
}

func (u *getClusterAccess) byTenant(ctx context.Context) error {
	// Get tenant by name
	tenant, err := u.tenantClient.GetByName(ctx, wrapperspb.String(u.tenantName))
	if err != nil {
		return err
	}

	stream, err := u.clusterAccessClient.GetTenantClusterMappingsByTenantId(ctx, wrapperspb.String(tenant.Id))
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		access, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		// Get cluster by id
		cluster, err := u.clusterClient.GetById(ctx, wrapperspb.String(access.ClusterId))
		if err != nil {
			return err
		}

		dataRow := []interface{}{
			cluster.Name,
			time.Since(access.Metadata.Created.AsTime()),
		}
		if u.outputOptions.ShowDeleted && cluster.Metadata.Deleted.AsTime().Unix() != 0 {
			dataRow = append(dataRow, time.Since(access.Metadata.Deleted.AsTime()))
		}
		data = append(data, dataRow)
	}

	header := []string{"CLUSTER", "AGE"}
	if u.outputOptions.ShowDeleted {
		header = append(header, "DELETED")
	}

	output.NewTableFactory().
		SetHeader(header).
		SetColumnFormatter("AGE", output.DefaultAgeColumnFormatter()).
		SetColumnFormatter("DELETED", output.DefaultAgeColumnFormatter()).
		SetSortColumn(u.outputOptions.SortOptions.SortByColumn).
		SetSortOrder(u.outputOptions.SortOptions.Order).
		SetData(data).
		ToTable().
		Render()

	return nil
}

func (u *getClusterAccess) byCluster(ctx context.Context) error {
	// Get cluster by name
	cluster, err := u.clusterClient.GetByName(ctx, wrapperspb.String(u.clusterName))
	if err != nil {
		return err
	}

	stream, err := u.clusterAccessClient.GetTenantClusterMappingsByClusterId(ctx, wrapperspb.String(cluster.Id))
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		access, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		// Get tenant by id
		tenant, err := u.tenantClient.GetById(ctx, wrapperspb.String(access.TenantId))
		if err != nil {
			return err
		}

		dataRow := []interface{}{
			tenant.Name,
			time.Since(access.Metadata.Created.AsTime()),
		}
		if u.outputOptions.ShowDeleted && tenant.Metadata.Deleted.AsTime().Unix() != 0 {
			dataRow = append(dataRow, time.Since(access.Metadata.Deleted.AsTime()))
		}
		data = append(data, dataRow)
	}

	header := []string{"TENANT", "AGE"}
	if u.outputOptions.ShowDeleted {
		header = append(header, "DELETED")
	}

	output.NewTableFactory().
		SetHeader(header).
		SetColumnFormatter("AGE", output.DefaultAgeColumnFormatter()).
		SetColumnFormatter("DELETED", output.DefaultAgeColumnFormatter()).
		SetSortColumn(u.outputOptions.SortOptions.SortByColumn).
		SetSortOrder(u.outputOptions.SortOptions.Order).
		SetData(data).
		ToTable().
		Render()

	return nil
}

func (u *getClusterAccess) Run(ctx context.Context) error {
	err := u.init(ctx)
	if err != nil {
		return err
	}

	if len(u.tenantName) > 0 {
		return u.byTenant(ctx)
	} else if len(u.clusterName) > 0 {
		return u.byCluster(ctx)
	}

	return errors.New("neither tenant nor cluster has been specified")
}
