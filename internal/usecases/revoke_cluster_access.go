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
	m8Grpc "github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type revokeClusterAccessUseCase struct {
	useCaseBase
	conn                *grpc.ClientConn
	tenantName          string
	clusterName         string
	tenantClient        api.TenantClient
	clusterClient       api.ClusterClient
	clusterAccessClient api.ClusterAccessClient
	cmdHandlerClient    esApi.CommandHandlerClient
}

func NewRevokeClusterAccessUseCase(config *config.Config, tenantName, clusterName string) UseCase {
	useCase := &revokeClusterAccessUseCase{
		useCaseBase: NewUseCaseBase("revoke-cluster-access", config),
		tenantName:  tenantName,
		clusterName: clusterName,
	}
	return useCase
}

func (u *revokeClusterAccessUseCase) init(ctx context.Context) error {
	if u.initialized {
		return nil
	}

	conn, err := m8Grpc.CreateGrpcConnectionAuthenticatedFromConfig(ctx, u.config)
	if err != nil {
		return err
	}

	u.conn = conn
	u.clusterClient = api.NewClusterClient(conn)
	u.tenantClient = api.NewTenantClient(conn)
	u.clusterAccessClient = api.NewClusterAccessClient(conn)
	u.cmdHandlerClient = esApi.NewCommandHandlerClient(conn)
	u.setInitialized()

	return nil
}

func (u *revokeClusterAccessUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	err := u.init(ctx)
	if err != nil {
		return err
	}

	// Get tenant by name
	tenant, err := u.tenantClient.GetByName(ctx, wrapperspb.String(u.tenantName))
	if err != nil {
		return err
	}

	// Get cluster by name
	cluster, err := u.clusterClient.GetByName(ctx, wrapperspb.String(u.clusterName))
	if err != nil {
		return err
	}

	// Get binding
	binding, err := u.clusterAccessClient.GetTenantClusterMappingByTenantAndClusterId(ctx, &api.GetClusterMappingRequest{
		TenantId:  tenant.Id,
		ClusterId: cluster.Id,
	})
	if err != nil {
		return err
	}

	// Delete binding
	command := cmd.NewCommand(uuid.MustParse(binding.Id), commandTypes.DeleteTenantClusterBinding)
	_, err = u.cmdHandlerClient.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Printf("Access to cluster '%s' revokeed for tenant '%s'.\n", u.clusterName, u.tenantName)
	}

	return err
}
