# Overview

Dummy Controller is a simple example Kubernetes controller that creates an nginx pod for every Dummy custom resource created.

```yaml
apiVersion: interview.interview.com/v1alpha1
kind: Dummy
metadata:
  name: dummy1
  namespace: example
spec:
  message: "I'm just a dummy"
```

The message in the spec is used only to populate a string in the status. Additionally, the status also contains a copy of the pod's current phase

```yaml
status:
  specEcho: "I'm just a dummy"
  podStatus: "Running"
```

The image built for this repo is pushed to [Docker Hub](https://hub.docker.com/r/lilshah/dummy-controller).

# Installing Kind
To run the image, you will need a Kubernetes cluster (for example, one created with [Minikube](https://minikube.sigs.k8s.io/docs/start/) or [Kind](https://kind.sigs.k8s.io/)).
You can use [Homebrew](https://brew.sh/) (mac, linux) or [Chocolatey](https://chocolatey.org/) (windows) to install Kind. Homebrew is installed via this command:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

To install Chocolatey on windows, you must ensure `Get-ExecutionPolicy` is not `Restricted`. Use `Bypass` to bypass the policy to get things installed or `AllSigned` for quite a bit more security.

In an elevated Powershell environment, run `Get-ExecutionPolicy`. If it returns `Restricted`, then run `Set-ExecutionPolicy AllSigned` or `Set-ExecutionPolicy Bypass -Scope Process`.

Now run the following command:

```sh
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Chocolatey should be ready to use now.

To install Kind, run the relevant command:

MacOS/Linux: `brew install kind`

Windows: `choco install kind`

You can now create and run a local cluster with `kind create cluster`

# Test the controller

The simplest method to install the controller for testing is to clone this repo, and from within it, run the following command:

```sh
make deploy
```

Give it a few minutes to download and create the pod. You can check pod status by using:

```sh
kubectl get pods -n tmp-homework-ms-system
```

A sample Dummy CR is also provided within the repo:

```sh
kubectl apply -f config/samples/interview_v1alpha1_dummy.yaml
```

You will need to install [kubectl](https://kubernetes.io/docs/tasks/tools/) to run the above commands.
