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
	api_commandhandler "github.com/finleap-connect/monoskope/pkg/api/domain"
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

	header := []string{"ID", "ROLE", "SCOPE", "RESOURCE", "AGE"}
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

func (u *getRoleBindingsUseCase) Run(ctx context.Context) error {
	conn, err := grpc.CreateGrpcConnectionAuthenticated(ctx, u.config.Server, &oauth2.Token{AccessToken: u.config.AuthInformation.Token})
	if err != nil {
		return err
	}
	defer conn.Close()
	grpcClient := api_commandhandler.NewUserClient(conn)

	user, err := grpcClient.GetByEmail(ctx, wrapperspb.String(u.email))
	if err != nil {
		return err
	}

	roleBindingstream, err := grpcClient.GetRoleBindingsById(ctx, wrapperspb.String(user.Id))
	if err != nil {
		return err
	}

	var data [][]interface{}
	for {
		// Read next
		rolebinding, err := roleBindingstream.Recv()
		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		row := []interface{}{
			rolebinding.Id,
			rolebinding.Role,
			rolebinding.Scope,
			rolebinding.Resource,
			time.Since(rolebinding.GetMetadata().GetCreated().AsTime()),
		}
		if u.showDeleted && rolebinding.GetMetadata().GetDeleted().AsTime().Unix() != 0 {
			row = append(row, time.Since(rolebinding.GetMetadata().GetDeleted().AsTime()))
		}
		data = append(data, row)
	}

	u.tableFactory.SetData(data) // Add Bulk Data
	u.tableFactory.ToTable().Render()

	return nil
}
