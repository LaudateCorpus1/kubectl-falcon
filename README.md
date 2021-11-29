## Deprecated

Kubectl-falcon project has been deprecated. Please consider adopting the following projects instead:
 - falcon-operator - Kubernetes operator to manage deployment of Falcon sensors to the cluster
 - falcon-helm - helm chart to manage deployment of Falcon sensors to the cluster
 - cloud-tools-image - collection of tools for Falcon Container manipulation in cloud settings
 

# kubectl-falcon ![Build CI](https://github.com/CrowdStrike/kubectl-falcon/workflows/Build%20CI/badge.svg)

[Falcon](https://www.crowdstrike.com/) is the [CrowdStrike](https://www.crowdstrike.com/)
platform purpose-built to stop breaches via a unified set of cloud-delivered
technologies that prevent all types of attacks — including malware and much
more.

kubectl-falcon is plug-in to kubectl command-line tool that helps with manipulation of
Falcon Container.

kubectl-falcon is an open source project, not CrowdStrike product. As such it carries
no formal support, expressed or implied.

This is pre-release version of kubectl-falcon.


## Installation

```
go get -u -v github.com/crowdstrike/kubectl-falcon
```

## Usage
```
kubectl falcon help
```
