# app-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### How we develop the operator
- check the development environment
```sh
docker version # 24.0.6 & API version: 1.43
go version #  go1.22.0 darwin/amd64
kubectl version --short # Client Version: v1.27.2 & Kustomize Version: v5.0.1
minikube version # v1.33.1
```
- create one project
```sh
mkdir MyOperatorProject
cd MyOperatorProject
mkdir app-operator
cd app-operator
kubebuilder init --domain=wlin.cn \
--repo=github.com/testcara/app-operator \
--owner cara
```
- analyze the current generated project files
```sh
carawang@carawangs-MacBook-Pro app-operator % tree
.
├── Dockerfile
├── Makefile # make cmds for build, deploy, test...
├── PROJECT # Kubebuilder related metadata including projectName, repo, version...
├── README.md
├── cmd
│   └── main.go # include the basic dependency and logics for Manager, Controller and Webhook...
├── config # config files for operator deplopyment
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_metrics_patch.yaml
│   │   └── metrics_service.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── network-policy
│   │   ├── allow-metrics-traffic.yaml
│   │   └── kustomization.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   └── rbac
│       ├── kustomization.yaml
│       ├── leader_election_role.yaml
│       ├── leader_election_role_binding.yaml
│       ├── metrics_auth_role.yaml
│       ├── metrics_auth_role_binding.yaml
│       ├── metrics_reader_role.yaml
│       ├── role.yaml
│       ├── role_binding.yaml
│       └── service_account.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── test
    ├── e2e
    │   ├── e2e_suite_test.go
    │   └── e2e_test.go
    └── utils
        └── utils.go

12 directories, 29 files
```
- create api
```sh
kubebuilder create api \
> --group apps --version v1 --kind App
```
and enter to 'y' to ensure it helps us to create resources and controllers
- analyze the generated files
```sh
carawang@carawangs-MacBook-Pro app-operator % tree
.
├── api
│   └── v1
│       ├── app_types.go # including Spec and Struct definition for our crd
│       ├── groupversion_info.go
│       └── zz_generated.deepcopy.go
├── bin
│   ├── controller-gen -> /Users/carawang/MyOperatorProject/app-operator/bin/controller-gen-v0.16.1
│   └── controller-gen-v0.16.1
├── config
│   ├── crd
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── rbac # resource permission files
│   │   ├── app_editor_role.yaml
│   │   ├── app_viewer_role.yaml
│   └── samples
│       ├── apps_v1_app.yaml # example cr files
│       └── kustomization.yaml
└── internal
    └── controller
       ├── app_controller.go # including the logic how we control our crd
       ├── app_controller_test.go
       └── suite_test.go
19 directories, 43 files
```
- Update the api and controller to meet our requirements


### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/app-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/app-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/app-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/app-operator/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024 cara.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

