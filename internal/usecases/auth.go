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
	"fmt"
	"strings"
	"text/template"

	_ "embed"

	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/grpc"
	"github.com/finleap-connect/monoctl/internal/spinner"
	"github.com/finleap-connect/monoctl/internal/version"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	monoctlAuth "github.com/finleap-connect/monoskope/pkg/auth"
	"github.com/pkg/browser"
	"golang.org/x/sync/errgroup"
)

// authUseCase provides the internal use-case of authentication.
type authUseCase struct {
	useCaseBase
	configManager *config.ClientConfigManager
	force         bool
	silent 		  bool
}

func NewAuthUsecase(configManager *config.ClientConfigManager, force, silent bool) UseCase {
	useCase := &authUseCase{
		useCaseBase:   NewUseCaseBase("authentication", configManager.GetConfig()),
		configManager: configManager,
		force:         force,
		silent: 	   silent,
	}
	return useCase
}

func (u *authUseCase) runAuthenticationFlow(ctx context.Context) error {
	u.log.Info("starting authentication")
	s := spinner.NewSpinner()
	defer s.Stop()

	conn, err := grpc.CreateGrpcConnection(ctx, u.config.Server)
	if err != nil {
		return fmt.Errorf("failed to connect the m8 control plane: %w", err)
	}
	defer conn.Close()
	gatewayClient := api.NewGatewayClient(conn)

	ready := make(chan string, 1)
	defer close(ready)

	indexPage, err := u.renderLocalServerSuccessHTML(u.config.Server, version.Version, version.Commit)
	if err != nil {
		return err
	}
	callbackServer, err := monoctlAuth.NewServer(&monoctlAuth.Config{
		LocalServerBindAddress: []string{
			"localhost:8000",
			"localhost:18000",
		},
		RedirectURLHostname:    "localhost",
		LocalServerSuccessHTML: indexPage,
		LocalServerReadyChan:   ready,
	})
	if err != nil {
		return err
	}
	defer callbackServer.Close()

	upstreamResponse, err := gatewayClient.RequestUpstreamAuthentication(ctx, &api.UpstreamAuthenticationRequest{
		CallbackUrl: callbackServer.RedirectURI,
	})
	if err != nil {
		return err
	}

	var authCode string
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			u.log.Info("Open " + url)
			if err := browser.OpenURL(url); err != nil {
				u.log.Error(err, "could not open the browser")
				return err
			}
			s.Stop()
			u.print("+-----------------------------------------------------------+\n")
			u.print("|monoctl has opened the browser for you to authenticate.    |\n")
			u.print("|It should show the log in screen of your identity provider |\n")
			u.print("|or the consent window if you are already logged in with it.|\n")
			u.print("+-----------------------------------------------------------+\n")
			u.print("Waiting for you to log in and give consent for OIDC flow...\n")
			s.Start()
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		var innerErr error
		authCode, innerErr = callbackServer.ReceiveCodeViaLocalServer(ctx, upstreamResponse.UpstreamIdpRedirect, upstreamResponse.State)
		return innerErr
	})
	if err := eg.Wait(); err != nil {
		u.log.Error(err, "authorization error: %s")
		return err
	}

	authResponse, err := gatewayClient.RequestAuthentication(ctx, &api.AuthenticationRequest{Code: authCode, State: upstreamResponse.State})
	if err != nil {
		return err
	}

	u.config.AuthInformation = &config.AuthInformation{
		Token:    authResponse.GetAccessToken(),
		Username: authResponse.GetUsername(),
	}
	if authResponse.Expiry != nil {
		expiry := authResponse.GetExpiry().AsTime()
		u.config.AuthInformation.Expiry = expiry
	}

	s.Stop()
	u.print("You're successfully authenticated as '%s'.\n", authResponse.GetUsername())
	u.print("---\n")
	u.print("\n")

	return u.configManager.SaveConfig()
}

func (u *authUseCase) Run(ctx context.Context) error {
	// Check if already authenticated
	if !u.force && u.config.HasAuthInformation() {
		u.log.Info("checking expiration of existing token")
		authInfo := u.config.AuthInformation
		if authInfo.IsValid() {
			u.log.Info("you have a valid auth token", "expiry", authInfo.Expiry.String())
			return nil
		}
		u.log.Info("your auth token has expired", "expiry", authInfo.Expiry)
	}
	return u.runAuthenticationFlow(ctx)
}

type IndexPageRenderData struct {
	M8Address string
	Version   string
}

// DefaultLocalServerSuccessHTML is a default response body on authorization success.
//go:embed CallbackServerSuccessPage.html
var DefaultLocalServerSuccessHTML string

func (u *authUseCase) renderLocalServerSuccessHTML(apiAddress string, version string, commit string) (string, error) {
	data := IndexPageRenderData{
		M8Address: apiAddress,
		Version:   version,
	}

	if strings.Contains(version, "-local") {
		data.Version = fmt.Sprintf("%s (commit %s)", version, commit)
	}

	// Create a new template and parse the document into it.
	t := template.Must(template.New("IndexPage").Parse(DefaultLocalServerSuccessHTML))

	outBuf := new(bytes.Buffer)
	err := t.Execute(outBuf, data)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil
}

func (u *authUseCase) print(format string, a ...interface{}) {
	if !u.silent {
		fmt.Printf(format, a...)
	}
}