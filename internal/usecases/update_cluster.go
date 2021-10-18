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

	"github.com/google/uuid"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type updateClusterUseCase struct {
	useCaseBase
	name                string
	newDisplayName      string
	newApiServerAddress string
	newCaCertBundle     []byte
}

func NewUpdateClusterUseCase(config *config.Config, name, newDisplayName, newApiServerAddress string, newCaCertBundle []byte) UseCase {
	useCase := &updateClusterUseCase{
		useCaseBase:         NewUseCaseBase("update-cluster", config),
		name:                name,
		newDisplayName:      newDisplayName,
		newApiServerAddress: newApiServerAddress,
		newCaCertBundle:     newCaCertBundle,
	}
	return useCase
}

func (u *updateClusterUseCase) Run(ctx context.Context) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()

	clusterClient := api.NewClusterClient(conn)
	cluster, err := clusterClient.GetByName(ctx, wrapperspb.String(u.name))
	if err != nil {
		return err
	}

	command := cmd.CreateCommand(uuid.MustParse(cluster.Id), commandTypes.UpdateCluster)
	cmdData := new(cmdData.UpdateCluster)
	cmdData.CaCertBundle = u.newCaCertBundle
	if u.newDisplayName != "" {
		cmdData.DisplayName = wrapperspb.String(u.newDisplayName)
	}
	if u.newApiServerAddress != "" {
		cmdData.ApiServerAddress = wrapperspb.String(u.newApiServerAddress)
	}
	if _, err := cmd.AddCommandData(command, cmdData); err != nil {
		return err
	}

	cmdHandlerClient := esApi.NewCommandHandlerClient(conn)
	_, err = cmdHandlerClient.Execute(ctx, command)

	s.Stop()
	if err == nil {
		fmt.Printf("Cluster '%s' updated.\n", u.name)
	}

	return err
}
