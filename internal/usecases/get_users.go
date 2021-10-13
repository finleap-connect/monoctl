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
	api_commandhandler "github.com/finleap-connect/monoskope/pkg/api/domain"
	"golang.org/x/oauth2"
)

// getUsersUseCase provides the internal use-case of getting the permission model.
type getUsersUseCase struct {
	useCaseBase
	tableFactory  *output.TableFactory
	outputOptions *output.OutputOptions
}

func NewGetUsersUseCase(config *config.Config, outputOptions *output.OutputOptions) UseCase {
	useCase := &getUsersUseCase{
		useCaseBase:   NewUseCaseBase("get-users", config),
		outputOptions: outputOptions,
	}

	header := []string{"NAME", "EMAIL", "AGE"}
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

func (u *getUsersUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	grpcClient := api_commandhandler.NewUserClient(conn)

	userStream, err := grpcClient.GetAll(ctx, &api_commandhandler.GetAllRequest{
		IncludeDeleted: u.outputOptions.ShowDeleted,
	})
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		user, err := userStream.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		dataline := []interface{}{
			user.Name,
			user.Email,
			time.Since(user.GetMetadata().GetCreated().AsTime()),
		}

		if u.outputOptions.ShowDeleted && user.Metadata.Deleted.AsTime().Unix() != 0 {
			dataline = append(dataline, time.Since(user.Metadata.Deleted.AsTime()))
		}

		data = append(data, dataline)
	}

	u.tableFactory.SetData(data) // Add Bulk Data
	u.tableFactory.ToTable().Render()

	return nil
}
