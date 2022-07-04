module github.com/crowdstrike/kubectl-falcon

go 1.15

require (
	github.com/containers/image/v5 v5.10.1
	github.com/crowdstrike/gofalcon v0.0.0-20210201183550-10e0ebcd5c85
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.4.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	k8s.io/cli-runtime v0.25.0-alpha.2
)

replace (
	github.com/go-openapi/analysis => github.com/go-openapi/analysis v0.19.5
	github.com/go-openapi/loads => github.com/go-openapi/loads v0.19.5
	// pin to older versions of go-openapi as these are used by current cli-runtime (through kustomize)
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.19.5
)
