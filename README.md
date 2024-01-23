# `ghcr.io/amp-buildpacks/aleo`

A Cloud Native Buildpack that provides the Aleo Tool Suite

## Configuration

| Environment Variable      | Description                                                                                                                                                                                                                                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$BP_ALEO_VERSION` | Configure the version of snarkOS to install. It can be a specific version or a wildcard like `1.*`. It defaults to the latest `1.*` version. |
| `$BPL_ENABLE_ALEO_DEPLOY` | Enable Aleo deployment features. Default: `true` |
| `$BPL_ALEO_DEPLOY_PRIVATE_KEY` | The Aleo wallet private key. |
| `$BPL_ALEO_DEPLOY_API_URL` | The endpoint to query node state from Aleo network. Default: `https://api.explorer.aleo.org/v1` |
| `$BPL_ALEO_DEPLOY_PRIORITY_FEE` | The priority fee in microcredits of deploy. Default: `100` |

## Usage

### 1. To use this buildpack, simply run:

```shell
pack build <image-name> \
    --path <aleo-samples-path> \
    --buildpack ghcr.io/amp-buildpacks/aleo \
    --builder paketobuildpacks/builder-jammy-full
```

For example:

```shell
pack build aleo-sample \
    --path ./samples/leo/snarkvm \
    --buildpack ghcr.io/amp-buildpacks/aleo \
    --builder paketobuildpacks/builder-jammy-full
```

### 2. To run the image, simply run:

```shell
docker run -u <uid>:<gid> -it <image-name>
```

For example:

```shell
docker run -u 1001:cnb -it aleo-sample
```

## Contributing

If anything feels off, or if you feel that some functionality is missing, please
check out the [contributing
page](https://docs.amphitheatre.app/contributing/). There you will find
instructions for sharing your feedback, building the tool locally, and
submitting pull requests to the project.

## License

Copyright (c) The Amphitheatre Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Credits

Heavily inspired by https://github.com/paketo-community/rustup
