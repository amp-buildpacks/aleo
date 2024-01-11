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
	"fmt"
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

func NewAleo(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) Aleo {
	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Build:  true,
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

func (r Aleo) Name() string {
	return r.LayerContributor.LayerName()
}
