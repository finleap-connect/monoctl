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

package grpc

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/config"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateGrpcConnection(ctx context.Context, url string) (*ggrpc.ClientConn, error) {
	factory := grpc.NewGrpcConnectionFactory(url).WithOSCaTransportCredentials()
	return factory.WithRetry().WithBlock().Connect(ctx)
}

func CreateGrpcConnectionAuthenticated(ctx context.Context, url string, token *oauth2.Token) (*ggrpc.ClientConn, error) {
	factory := grpc.NewGrpcConnectionFactory(url).WithOSCaTransportCredentials()
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		factory = factory.WithPerRPCCredentials(oauth.NewOauthAccess(token))
	}
	return factory.WithRetry().WithBlock().Connect(ctx)
}

func CreateGrpcConnectionAuthenticatedFromConfig(ctx context.Context, config *config.Config) (*ggrpc.ClientConn, error) {
	conn, err := CreateGrpcConnectionAuthenticated(ctx, config.Server, &oauth2.Token{AccessToken: config.AuthInformation.Token})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// dummy connections to be used for testing
func CreateDummyGrpcConnection() *ggrpc.ClientConn {
	return &ggrpc.ClientConn{}
}
