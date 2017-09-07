# deisrel

[![Build Status](https://travis-ci.org/deis/deisrel.svg?branch=master)](https://travis-ci.org/deis/deisrel)
[![codebeat badge](https://codebeat.co/badges/46e06b60-7e4c-4daf-875b-c7c07ee56035)](https://codebeat.co/projects/github-com-deis-deisrel)

[Download for 64 Bit Linux](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)

[Download for 64 Bit Darwin](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a
developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to
deploy and manage applications on your own servers.

For more information about the Deis Workflow, please visit the main project page at
<https://github.com/deisthree/workflow>.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

# About

`deisrel` is a utility tool for assisting with Deis product releases. It makes is easy to check what components need to be released before making a workflow release.

# Installing deisrel

You can install the latest version of `deisrel` from the following links:

- [linux-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-linux-amd64)
- [darwin-amd64](https://storage.googleapis.com/deisrel/deisrel-latest-darwin-amd64)

Alternatively, you can compile this project from source using Go 1.6+:

	$ git clone https://github.com/deisthree/deisrel
	$ cd deisrel
	$ make bootstrap build
	$ ./deisrel

Once done, you can then move the client binary anywhere on your PATH:

	$ mv deisrel /usr/local/bin/

# Usage

deisrel requires two files, the `requirements.lock` file from a [Kubernetes Helm](https://github.com/kubernetes/helm) Workflow chart,
and a repository mapping file, allowing deisrel to map between a repository and chart name.

You may fetch the `requirements.lock` file from a given chart in the following manner:

	$ helm repo add deis https://charts.deis.com/workflow
	$ helm fetch --untar deis/workflow --version v2.9.0
	$ cat workflow/requirements.lock


As the Workflow component list has stabilized, the repository mapping file (`map.json`) at the root of this project may be used.  (However, it will need to be updated when/if the list changes.)  The key is the repository name, and the value is the chart name, as found in the `requirements.lock` file.

With these two files, you can use deisrel to generate a report:

```console
$ deisrel path/to/requirements.lock path/to/repomapping.json
builder                  v2.6.1 -> v2.7.1 (dirty)
	builder has unreleased changes. See https://github.com/deisthree/builder/compare/v2.7.1...master
controller               v2.10.0 -> v2.11.0 (dirty)
	controller has unreleased changes. See https://github.com/deisthree/controller/compare/v2.11.0...master
database                 v2.5.0 -> v2.5.1 (clean)
dockerbuilder            v2.5.2 -> v2.6.0 (clean)
fluentd                  v2.6.0 -> v2.6.1 (clean)
logger                   v2.4.1 -> v2.4.2 (clean)
minio                    v2.3.4 -> v2.3.5 (clean)
monitor                  v2.7.0 -> v2.7.1 (clean)
nsqd                     v2.2.5 -> v2.2.6 (clean)
redis                    v2.2.4 -> v2.2.5 (clean)
registry                 v2.3.2 -> v2.3.3 (clean)
registry-proxy           v1.1.1 -> v1.1.1 (dirty)
	registry-proxy has unreleased changes. See https://github.com/deisthree/registry-proxy/compare/v1.1.1...master
registry-token-refresher v1.1.0 -> v1.1.1 (clean)
router                   v2.8.1 -> v2.9.0 (clean)
slugbuilder              v2.4.8 -> v2.4.9 (clean)
slugrunner               v2.2.4 -> v2.2.4 (dirty)
	slugrunner has unreleased changes. See https://github.com/deisthree/slugrunner/compare/v2.2.4...master
workflow-manager         v2.4.4 -> v2.4.5 (clean)
```

It's also possible to output the report in json using the `-o json` argument for easier machine parsing.

Github has a very low ratelimit for unauthenticated api requests. A "[personal access token][]" can be used to bypass this restriction. Create a GitHub token with "repo" permission and set it to `$GITHUB_ACCESS_TOKEN`. `deisrel` will
use it when making requests.

# Development

The Deis project welcomes contributions from all developers. The high level process for development matches many other open source projects. See below for an outline.

* Fork this repository
* Make your changes
* [Submit a pull request][prs] (PR) to this repository with your changes, and unit tests whenever possible
	* If your PR fixes any [issues][issues], make sure you write `Fixes #1234` in your PR description (where `#1234` is the number of the issue you're closing)
* The Deis core contributors will review your code. After each of them sign off on your code, they'll label your PR with `LGTM1` and `LGTM2` (respectively). Once that happens, a contributor will merge it

# License

Copyright 2013, 2014, 2015, 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


[issues]: https://github.com/deisthree/deisrel/issues
[personal access token]: https://github.com/settings/tokens
[prs]: https://github.com/deisthree/deisrel/pulls
[workflow]: https://github.com/deisthree/workflow
