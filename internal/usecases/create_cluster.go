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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
	"time"

	_ "embed"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/grpc"
	"gitlab.figo.systems/platform/monoskope/monoctl/internal/spinner"
	domApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	jwt              string

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

func (u *createClusterUseCase) queryJwt(ctx context.Context, id string) error {
	s := spinner.NewSpinner()
	defer s.Stop()

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 10 * time.Second
	err := backoff.Retry(func() error {
		wrapper, err := u.clusterClient.GetBootstrapToken(ctx, &wrapperspb.StringValue{Value: id})
		if err != nil {
			return err
		}
		u.jwt = wrapper.GetValue()
		if u.jwt == "" {
			return fmt.Errorf("u.jwt not set yet")
		}
		return nil
	}, params)
	if err != nil {
		fmt.Printf("Receiving Cluster bootstrap token for cluster '%s' failed: %s.\n", u.displayName, err)
		return err
	}

	s.Stop()
	fmt.Printf("Cluster '%s' created.\n", u.displayName)

	return err
}

// Prepare some data to insert into the template.
type ClusterRenderData struct {
	ApiServerAddress string
	M8Endpoint       string
	Jwt              string
	ClusterName      string
}

func (u *createClusterUseCase) renderOutput() (string, error) {

	encodedJwt := base64.StdEncoding.EncodeToString([]byte(u.jwt))
	var data = ClusterRenderData{
		ApiServerAddress: u.apiServerAddress,
		M8Endpoint:       u.config.Server,
		Jwt:              encodedJwt,
		ClusterName:      u.displayName,
	}

	return RenderClusterConfig(data)
}

// template.for certificate resource and issuer
//go:embed m8_operator_bootstrap.yaml
var resource string

func RenderClusterConfig(data ClusterRenderData) (string, error) {
	// Create a new template and parse the document into it.
	t := template.Must(template.New("resource").Parse(resource))

	outBuf := new(bytes.Buffer)
	err := t.Execute(outBuf, data)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil
}

func (u *createClusterUseCase) Run(ctx context.Context) error {
	err := u.setUp(ctx)
	if err != nil {
		return err
	}
	defer u.conn.Close()

	reply, err := u.doCreate(ctx)
	if err != nil {
		return err
	}

	err = u.queryJwt(ctx, reply.AggregateId)
	if err != nil {
		return err
	}

	out, err := u.renderOutput()
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, out)

	return nil
}
