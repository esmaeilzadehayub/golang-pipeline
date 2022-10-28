# golang-pipeline
# golang-pipeline

# Overview

. Go version 1.18 or later. You might want to download and install Go first.

. Docker running locally. Follow the instructions to download and install Docker.

. An IDE or a text editor to edit files. We recommend using Visual Studio Code.

. Helm 

# Create a Dockerfile for the application

Next, we need to add a line in our Dockerfile that tells Docker what base image we would like to use for our application.

```Dockerfile

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /hostname

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
