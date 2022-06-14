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

	_ "embed"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	domApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
)

/*
 * you can retrieve the CA cert bundle via management UIs like Gardener or GCP. It will also be available
 * in the .kubeconfig if you have log into the cluster via an admin token.
 */

type createClusterUseCase struct {
	useCaseBase
	displayName      string
	name             string
	apiServerAddress string
	caCertBundle     []byte

	conn           *ggrpc.ClientConn
	cHandlerClient esApi.CommandHandlerClient
	clusterClient  domApi.ClusterClient
}

func NewCreateClusterUseCase(config *config.Config, name, label, apiServerAddress string, caCertBundle []byte) UseCase {
	useCase := &createClusterUseCase{
		useCaseBase:      NewUseCaseBase("create-cluster", config),
		displayName:      name,
		name:             label,
		apiServerAddress: apiServerAddress,
		caCertBundle:     caCertBundle,
	}
	return useCase
}

func (u *createClusterUseCase) setUp(ctx context.Context) error {
	var err error

	u.conn, err = grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}

	u.cHandlerClient = esApi.NewCommandHandlerClient(u.conn)
	u.clusterClient = domApi.NewClusterClient(u.conn)

	return nil
}

func (u *createClusterUseCase) doCreate(ctx context.Context) (*esApi.CommandReply, error) {
	s := spinner.NewSpinner()
	defer s.Stop()

	// this is a create command; use nil as input, the correct ID will be contained in the reply
	command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster)
	if _, err := cmd.AddCommandData(command, &cmdData.CreateCluster{
		DisplayName:      u.displayName,
		Name:             u.name,
		ApiServerAddress: u.apiServerAddress,
		CaCertBundle:     u.caCertBundle,
	}); err != nil {
		return nil, err
	}

	reply, err := u.cHandlerClient.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Printf("Cluster '%s' created.\n", u.displayName)
	}

	return reply, err
}

// Prepare some data to insert into the template.
type ClusterRenderData struct {
	ApiServerAddress string
	M8Endpoint       string
	Jwt              string
	ClusterName      string
}

func (u *createClusterUseCase) Run(ctx context.Context) error {
	err := u.setUp(ctx)
	if err != nil {
		return err
	}
	defer u.conn.Close()

	_, err = u.doCreate(ctx)
	return err
}
