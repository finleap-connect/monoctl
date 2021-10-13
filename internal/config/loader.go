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
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"gitlab.figo.systems/platform/monoskope/monoctl/internal/util"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	RecommendedConfigPathEnvVar = "MONOSKOPECONFIG"
	RecommendedHomeDir          = ".monoskope"
	RecommendedFileName         = "config"
	FileMode                    = 0644
)

var (
	RecommendedConfigDir = path.Join(util.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile  = path.Join(RecommendedConfigDir, RecommendedFileName)
)

type ClientConfigManager struct {
	// Logger interface
	log          logger.Logger
	config       *Config
	configPath   string
	explicitFile string
}

// NewLoader is a convenience function that returns a new ClientConfigManager object with defaults
func NewLoader() *ClientConfigManager {
	return &ClientConfigManager{
		log: logger.WithName("client-config-loader"),
	}
}

// NewLoaderFromConfig is mainly intended to be used for tests. It allows setting the config without the need of an external file.
func NewLoaderFromConfig(conf *Config) *ClientConfigManager {
	return &ClientConfigManager{
		config: conf,
		log:    logger.WithName("client-config-loader"),
	}
}

// NewLoaderFromExplicitFile is a convenience function that returns a new ClientConfigManager object with explicitFile set
func NewLoaderFromExplicitFile(explicitFile string) *ClientConfigManager {
	// Expand home directory if exists
	if strings.HasPrefix(explicitFile, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		explicitFile = filepath.Join(dir, explicitFile[2:])
	}

	loader := NewLoader()
	loader.explicitFile = explicitFile
	return loader
}

// saveConfig saves the configuration
func (l *ClientConfigManager) saveConfig(filename string, force bool, config *Config) error {
	exists, err := util.FileExists(filename)
	if err != nil {
		return err
	}
	if exists && !force {
		return ErrAlreadyInitialized
	}
	l.config = config
	l.configPath = filename
	return l.SaveToFile(config, l.configPath, FileMode)
}

// loadConfig checks if the given file exists and loads it's contents
func (l *ClientConfigManager) loadConfig(filename string) error {
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return err
	}
	l.config, err = l.LoadFromFile(filename)

	// Load token from keyring
	l.config.LoadToken()

	return err
}

// GetConfig returns the previously loaded config
func (l *ClientConfigManager) GetConfig() *Config {
	return l.config
}

// GetConfigLocation returns the location of the previously loaded config
func (l *ClientConfigManager) GetConfigLocation() string {
	return l.configPath
}

func (l *ClientConfigManager) InitConfig(config *Config, force bool) error {
	if l.explicitFile != "" {
		return l.saveConfig(l.explicitFile, force, config)
	}

	envVarFile := os.Getenv(RecommendedConfigPathEnvVar)
	if len(envVarFile) != 0 {
		return l.saveConfig(envVarFile, force, config)
	}

	return l.saveConfig(RecommendedHomeFile, force, config)
}

func (l *ClientConfigManager) SaveConfig() error {
	if l.configPath == "" || l.config == nil {
		return ErrNoConfigExists
	}

	// Store token in keyring
	if err := l.config.StoreToken(); err != nil {
		return err
	}

	return l.SaveToFile(l.config, l.configPath, FileMode)
}

// LoadConfig loads the config either from env or home file.
func (l *ClientConfigManager) LoadConfig() error {
	if l.explicitFile != "" {
		if err := l.loadConfig(l.explicitFile); err != nil {
			return err
		}
		l.configPath = l.explicitFile
		return nil
	}

	// Load config from envvar path if provided
	envVarFile := os.Getenv(RecommendedConfigPathEnvVar)
	if len(envVarFile) != 0 {
		if err := l.loadConfig(envVarFile); err != nil {
			return err
		}
		l.configPath = envVarFile
		return nil
	}

	// Load recommended home file if present
	if err := l.loadConfig(RecommendedHomeFile); err != nil {
		if os.IsNotExist(err) {
			return ErrNoConfigExists
		}
		return err
	}
	l.configPath = RecommendedHomeFile

	return nil
}

// LoadFromFile takes a filename and deserializes the contents into Config object
func (l *ClientConfigManager) LoadFromFile(filename string) (*Config, error) {
	monoconfigBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := l.LoadFromBytes(monoconfigBytes)
	if err != nil {
		return nil, err
	}
	l.log.Info("Config loaded from file", "filename", filename)

	return config, nil
}

// LoadFromBytes takes a byte slice and deserializes the contents into Config object.
// Encapsulates deserialization without assuming the source is a file.
func (*ClientConfigManager) LoadFromBytes(data []byte) (*Config, error) {
	config := NewConfig()

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveToFile takes a config, serializes the contents and stores them into a file.
func (l *ClientConfigManager) SaveToFile(config *Config, filename string, permission os.FileMode) error {
	var err error

	// Create directory with any necessary parents
	err = os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	// Marshal config
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write config to file
	err = ioutil.WriteFile(filename, bytes, permission)
	if err != nil {
		return err
	}

	l.log.Info("Config saved to file", "filename", filename)

	return nil
}
