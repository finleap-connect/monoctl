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

package k8s

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/finleap-connect/monoctl/internal/prompt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

const kubeConfigEnvVar = "KUBECONFIG"

type KubeConfig struct {
	log        logger.Logger
	ConfigPath string
}

func NewKubeConfig() *KubeConfig {
	return &KubeConfig{
		log: logger.WithName("KubeConfig"),
	}
}

func (k *KubeConfig) LoadConfig() (*api.Config, error) {
	if k.ConfigPath == "" {
		k.ConfigPath = os.Getenv(kubeConfigEnvVar)
	}

	if k.ConfigPath != "" {
		fileList := filepath.SplitList(k.ConfigPath)
		if len(fileList) > 1 {
			var err error
			_, k.ConfigPath, err = prompt.SelectWithAdd(
				"please specify a config file",
				"Create new",
				fileList,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if k.ConfigPath == "" {
		k.ConfigPath = path.Join("~", ".kube", "config")
	}

	k.ConfigPath = os.ExpandEnv(k.ConfigPath)
	if strings.HasPrefix(k.ConfigPath, "~/") {
		homeDir := homedir.HomeDir()
		if homeDir == "" {
			return nil, errors.New("couldn't determine home directory")
		}
		k.ConfigPath = filepath.Join(homeDir, k.ConfigPath[2:])
	}

	if _, err := os.Stat(k.ConfigPath); err != nil {
		if os.IsNotExist(err) {
			return api.NewConfig(), nil
		} else {
			return nil, err
		}
	}

	config, err := clientcmd.LoadFromFile(k.ConfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (k *KubeConfig) StoreConfig(conf *api.Config) error {
	return clientcmd.WriteToFile(*conf, k.ConfigPath)
}

func (k *KubeConfig) SetPath(path string) {
	k.ConfigPath = path
}
