// Copyright (c) The Amphitheatre Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aleo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

// Aleo will handle installing the `leo` tool & adding it to the PATH
type Aleo struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

type AleoApp struct {
	Program string `json:"program"`
}

func NewAleo(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) Aleo {
	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Cache:  true,
		Launch: true,
	})
	return Aleo{
		LayerContributor: contributor,
	}
}

func (r Aleo) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	r.LayerContributor.Logger = r.Logger
	return r.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		bin := filepath.Join(layer.Path, "bin")

		r.Logger.Bodyf("Expanding %s to %s", artifact.Name(), bin)
		if err := crush.Extract(artifact, bin, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand %s\n%w", artifact.Name(), err)
		}

		file := filepath.Join(bin, "snarkos")
		r.Logger.Bodyf("Setting %s as executable", file)
		if err := os.Chmod(file, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to chmod %s\n%w", file, err)
		}

		r.Logger.Bodyf("Setting %s in PATH", bin)
		if err := os.Setenv("PATH", sherpa.AppendToEnvVar("PATH", ":", bin)); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to set $PATH\n%w", err)
		}
		return layer, nil
	})
}

func (r Aleo) BuildProcessTypes(cr libpak.ConfigurationResolver, app libcnb.Application) ([]libcnb.Process, error) {
	processes := []libcnb.Process{}

	enableDeploy := cr.ResolveBool("BPL_ENABLE_ALEO_DEPLOY")
	if enableDeploy {
		deployProcess, err := r.makeDeployCommand(cr, app)
		if err != nil {
			return processes, err
		}
		processes = append(processes, deployProcess)
	}
	return processes, nil
}

func (r Aleo) Name() string {
	return r.LayerContributor.LayerName()
}

// snarkos developer deploy "${APPNAME}.aleo" --private-key "${PRIVATE_KEY}" \
// --path "./build/" \
// --query "${API_URL}" \
// --broadcast "${API_URL}/testnet3/transaction/broadcast" \
// --priority-fee 100
func (r Aleo) makeDeployCommand(cr libpak.ConfigurationResolver, app libcnb.Application) (libcnb.Process, error) {
	process := libcnb.Process{}

	aleoApp, err := r.ReadAppConfig(app.Path)
	if err != nil {
		return process, err
	}

	privateKey, enableKey := cr.Resolve("BPL_ALEO_DEPLOY_PRIVATE_KEY")
	if !enableKey {
		return process, fmt.Errorf("BPL_ALEO_DEPLOY_PRIVATE_KEY must to be set")
	}

	apiUrl, _ := cr.Resolve("BPL_ALEO_DEPLOY_API_URL")
	priorityFee, _ := cr.Resolve("BPL_ALEO_DEPLOY_PRIORITY_FEE")

	process.Type = "web"
	process.Default = true
	process.Command = "snarkos"
	process.Arguments = []string{
		"developer",
		"deploy", aleoApp.Program,
		"--private-key", privateKey,
		"--path", filepath.Join(app.Path, "build"),
		"--query", apiUrl,
		"--broadcast", apiUrl + "/testnet3/transaction/broadcast",
		"--priority-fee", priorityFee,
	}
	return process, nil
}

func (r Aleo) ReadAppConfig(appDir string) (AleoApp, error) {
	aleoApp := AleoApp{}

	fileName := filepath.Join(appDir, "program.json")
	file, err := os.Open(fileName)
	if err != nil {
		return aleoApp, fmt.Errorf("unable to determine if '%s' exists\n%w", fileName, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return aleoApp, fmt.Errorf("unable to read '%s'\n%w", fileName, err)
	}

	if err := json.Unmarshal(data, &aleoApp); err != nil {
		return aleoApp, fmt.Errorf("unable to convert '%s'\n%w", fileName, err)
	}
	return aleoApp, nil
}
