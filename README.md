# golang-pipeline


# Overview

. Go version 1.18 or later. You might want to download and install Go first.

. Docker running locally. Follow the instructions to download and install Docker.

. An IDE or a text editor to edit files. We recommend using Visual Studio Code.

. Helm 

# Create a Dockerfile for the application

Next, we need to add a line in our Dockerfile that tells Docker what base image we would like to use for our application.

```Dockerfile

FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /hostname

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /hostname /hostname

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/hostname"]
EXPOSE 8000

CMD [ "/hostname" ]

```
# Build the image

Now that we’ve created our Dockerfile, let’s build an image from it. The docker build command creates Docker images from the Dockerfile and a “context”.

```bash

docker build --tag golang-hostname .

```

<img width="1027" alt="Screenshot 2022-10-29 at 00 33 03" src="https://user-images.githubusercontent.com/28998255/198745773-2e2c92c6-367c-405a-b4c7-f62af01da80a.png">


# View local images

```bash

docker image ls
```

# Tag images

Please set your repo 
```bash
docker image tag golang-hostname:latest YOURREPO.localgolang-hostname:v1.0
docker push YOURREPO.localgolang-hostname:v1.0
```
# Helm Charts

👷 Collection of Helm Charts for Kubernetes deployments ☸️

[[_TOC_]]

## Intro

[Helm](https://helm.sh/) is generally described as the package manager for Kubernetes and Charts are the format of packaging an application for Kubernetes.

A Helm Chart consists at least of the following components:

```plain
example-chart
├── Chart.yaml                # Contains information about the Chart
├── templates                 # Contains all Kubernetes manifests to be rendered
│   └── example-manifest.yaml
└── values.yaml               # Contains all variables the the manifests are rendered with
```

Creating a Helm Chart is as simple as gathering all your Kubernetes manifests in the `templates` folder and replacing everything you want to customize with variables ( `{{ .Values.<VARIABLE> }}`) you then place in your `values.yaml` file. Add the info about your Chart to `Chart.yaml` and you are ready to install it to Kubernetes with `helm install <DEPLOYMENT_NAME> <CHART_LOCATION> `.

## Best Practices

To create a Helm Chart quickly you can just run `helm create <CHART_NAME>`, which will create a generic Helm Chart for a web service. You can then customize the Chart to your liking.

This method ensures, that all Charts have the same structure, at least for their base, which makes it easier for CI and CD pipelines to work with the Charts.


### Testing your Chart

If you develop a new Helm Chart it's good to test it before deploying or pushing it. There are two basic ways to do this.

#### 1. `helm lint`

With `helm lint` you can lint a given Helm Chart to check for general errors and the following of best practices.
Simply run the following command:

```bash
helm lint <CHART> -f <PATH/TO/VALUES_FILE>
# for example
helm lint demo -f demo/values.yaml
```

#### 2. `helm template`

With `helm template` you can template a given Helm Chart, which will output the compiled Kubernetes manifests.
With this you see if your Helm Chart is able to compile and you can check the output to see what will be applied to Kubernetes.
Simply run the following command:

```bash
helm template <RELEASE_NAME> <CHART> -f <PATH/TO/VALUES_FILE>
# for example
helm template test-demo demo -f demo/values.yaml
```
# ArgoCD


[ArgoCD](https://argoproj.github.io/argo-cd/) is a declarative, GitOps continous delivery tool for Kubernetes. This means you describe the state of a resource in an ArgoCD config file and ArgoCD takes care that this resource in the cluster always comlies with the description you gave in the configfile.

By using ArgoCD and the concept of GitOps we get some major benefits:

- Our Kubernetes manifests stay valid
- We have a single source of truth for the application state on Kubernetes
- Easy versioning and rollback of changes to the application
.


## Workflow

The complete workflow of our CI + CD summarises to the following steps:

1. A developer pushes code to his projects repository
2. Gitlab CI performes the CI tasks specified in the `.gitlab-ci.yml` file, i.e.:
   - Running tests against the source code
   - Building a Dockerfile
   - Pushing the Dockerimage to a registry
   - Security & vulnerability scanning
3. After the above tasks one more CI task will update the corresponding Helm chart for this project, with the new Dockerimage tag
4. ArgoCD watches it's configured apps and will notice the changes in the Helm chart and will update the current deployment with the new version

![argocd-workflow](https://user-images.githubusercontent.com/28998255/198868788-86970581-c210-4b55-8f64-5615abb57523.png)


## Create a new ArgoCD managed app

ArgoCD is configured to manage all resources in the app directory.
Because ArgoCD provides Kubernetes CRDs for it's configuration we can place our
new ArgoCD app in the `apps` directory and ArgoCD will automatically pick it up.
The directory structure inside the `apps` directory represents our clusters,
but is just for better overview, the placement of an application definition
inside a specific directory does not mean it will be deployed to that cluster,
this is defined in the `spec.destination` section of the application definition.

An ArgoCD app definition looks something like the following:

```yaml
apiVersion: argoproj.io/v1alpha1                              # The Argo CRD version
kind: Application                                             # The Argo CRD kind
metadata:
  name: demo-app                                              # Name of the ArgoCD app
  namespace: argocd                                           # Namespace to deploy the app definition to (should always be `argocd`)
spec:
  project: development                                        # The ArgoCD project to deploy the app in
  source:
    repoURL: YourRepo  # The repo URL to get the kubernetes manifests from
    targetRevision: HEAD                                      # The branch to use
    path: demo                                                # The path inside the repo storing the manifests
    helm:
      valueFiles:                                             # A list of values files, inside the helm chart, to use for this deployment
        - values.yaml
  destination:
    server: YourServer    # The target cluster (Check Rancher cluster URLs)
    namespace: demo                                           # The namespace in which to deploy the manifests
  syncPolicy:
    automated:
      prune: true                                             # Prune the app on sync
      selfHeal: true                                          # Try to heal the app on failure
    syncOptions:
      - CreateNamespace=true                                  # Create the specified namespace, if non existent
```

To create a new ArgoCD managed app just create a YAML file like the one above in the `apps` directory and customize the parts you want to change.
