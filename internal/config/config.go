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

package config

import (
	"errors"
	"fmt"
	"time"

	keyring "github.com/zalando/go-keyring"
	"gopkg.in/yaml.v2"
)

var (
	ErrEmptyServer        = errors.New("has no server defined")
	ErrNoConfigExists     = errors.New("no valid monoconfig found")
	ErrAlreadyInitialized = errors.New("a configuartion already exists")
)

const (
	monoctlService = "monoskope/monoctl"
)

// Config holds the information needed to build connect to remote monoskope instance as a given user
type Config struct {
	// Server is the address of the Monoskope Gateway (https://hostname:port).
	Server string `yaml:"server"`
	// AuthInformation contains information to authenticate against Monoskope
	AuthInformation *AuthInformation `yaml:"authInformation,omitempty"`
	// ClusterAuthInformation contains information to authenticate against K8s clusters
	ClusterAuthInformation map[string]*AuthInformation `yaml:"clusterAuthInformation,omitempty"`
}

// NewConfig is a convenience function that returns a new Config object with defaults
func NewConfig() *Config {
	return &Config{
		ClusterAuthInformation: make(map[string]*AuthInformation),
	}
}

// Validate validates if the config is valid
func (c *Config) Validate() error {
	if c.Server == "" {
		return ErrEmptyServer
	}
	return nil
}

func (c *Config) StoreToken() error {
	if c.HasAuthInformation() {
		err := keyring.Set(monoctlService, c.AuthInformation.Username, c.AuthInformation.Token)
		if err != nil {
			return err
		}
	}
	for key, authInfo := range c.ClusterAuthInformation {
		err := keyring.Set(monoctlService, key, authInfo.Token)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) LoadToken() {
	if c.HasAuthInformation() {
		token, err := keyring.Get(monoctlService, c.AuthInformation.Username)
		if err == nil {
			c.AuthInformation.Token = token
		}
	}

	for key, authInfo := range c.ClusterAuthInformation {
		token, err := keyring.Get(monoctlService, key)
		if err != nil {
			delete(c.ClusterAuthInformation, key)
		} else {
			authInfo.Token = token
		}
	}
}

// HasAuthInformation checks if the the config contains AuthInformation
func (c *Config) HasAuthInformation() bool {
	return c.AuthInformation != nil
}

func (c *Config) GetClusterAuthInformation(clusterId, username, role string) *AuthInformation {
	return c.ClusterAuthInformation[fmt.Sprintf("%s/%s/%s", clusterId, username, role)]
}

func (c *Config) SetClusterAuthInformation(clusterId, username, role, token string, expiry time.Time) {
	c.ClusterAuthInformation[fmt.Sprintf("%s/%s/%s", clusterId, username, role)] = &AuthInformation{
		Username: username,
		Token:    token,
		Expiry:   expiry,
	}
}

func (c *Config) String() (string, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type AuthInformation struct {
	Username string    `yaml:"username,omitempty"`
	Token    string    `yaml:"-"`
	Expiry   time.Time `yaml:"expiry,omitempty"`
}

// IsValid checks that Token is not empty and is not expired with an offset
func (a *AuthInformation) IsValid() bool {
	return a.HasToken() && !a.IsTokenExpired()
}

// IsValid checks that Token is not empty and is not expired
func (a *AuthInformation) IsValidExact() bool {
	return a.HasToken() && !a.IsTokenExpiredExact()
}

// HasToken checks that Token is not empty
func (a *AuthInformation) HasToken() bool {
	return a.Token != ""
}

// IsTokenExpired checks if the auth token is expired with an offset
func (a *AuthInformation) IsTokenExpired() bool {
	return a.Expiry.IsZero() || a.Expiry.Before(time.Now().UTC().Add(5*time.Minute)) // check if token is valid for at least five more minutes
}

// IsTokenExpiredExact checks if the auth token is expired
func (a *AuthInformation) IsTokenExpiredExact() bool {
	return a.Expiry.IsZero() || a.Expiry.Before(time.Now().UTC().Add(1*time.Second)) // check if token is valid exactly
}
