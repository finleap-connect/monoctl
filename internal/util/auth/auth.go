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

package auth

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoctl/internal/spinner"
	"github.com/juju/clock"
	"github.com/juju/mutex/v2"
	"time"

	"github.com/finleap-connect/monoctl/cmd/monoctl/flags"
	"github.com/finleap-connect/monoctl/internal/config"
	"github.com/finleap-connect/monoctl/internal/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	authFlowNamedMutex         = "monoctl-auth-flow"
	namedMutexDetectionTimeOut = 500 * time.Millisecond
)

func RetryOnAuthFailSilently(ctx context.Context, configManager *config.ClientConfigManager, f func(ctx context.Context) error) error {
	return retryOnAuthFail(ctx, configManager, true, f)
}

func RetryOnAuthFail(ctx context.Context, configManager *config.ClientConfigManager, f func(ctx context.Context) error) error {
	return retryOnAuthFail(ctx, configManager, false, f)
}

func retryOnAuthFail(ctx context.Context, configManager *config.ClientConfigManager, silent bool, f func(ctx context.Context) error) error {
	// Make sure no other process run the authentication flow.
	lock, err := acquireLock(ctx, silent)
	if err != nil {
		return err
	}
	defer lock.Release()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*2) // special timeout for login flow with consent can take longer
	defer cancel()

	if err := LoadConfigAndAuth(ctx, configManager, flags.ForceAuth, silent); err != nil {
		return fmt.Errorf("init failed: %w", err)
	}

	if err := f(ctx); err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.Unauthenticated {
			if err := LoadConfigAndAuth(ctx, configManager, true, silent); err != nil {
				return fmt.Errorf("init failed: %w", err)
			}
			return f(ctx)
		}
		return err
	}
	return nil
}

func LoadConfigAndAuth(ctx context.Context, configManager *config.ClientConfigManager, force, silent bool) error {
	if err := configManager.LoadConfig(); err != nil {
		return fmt.Errorf("failed loading monoconfig: %w", err)
	}
	return usecases.NewAuthUsecase(configManager, force, silent).Run(ctx)
}

func acquireLock(_ context.Context, silent bool) (mutex.Releaser, error) {
	var s *spinner.Spinner
	if !silent {
		s = spinner.NewSpinner()
		defer s.Stop()
	}

	acquired := make(chan struct{})
	defer close(acquired)
	go func() {
		select {
		case <-acquired:
			return
		case <-time.After(namedMutexDetectionTimeOut):
			if !silent {
				s.Stop()
				fmt.Printf("Another monoctl instance is already running the authentication flow...\n")
				s.Start()
			}
		}
	}()

	return mutex.Acquire(mutex.Spec{
		Name:  authFlowNamedMutex,
		Clock: clock.WallClock,
		Delay: time.Second,
	})
}
