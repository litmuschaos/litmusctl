module github.com/mayadata-io/kuberactl

go 1.14

require (
	github.com/go-resty/resty/v2 v2.3.0
	github.com/mayadata-io/cli-utils v0.0.0-20210113084301-86afa3dfe6c6
	github.com/spf13/cobra v1.1.1
	k8s.io/api v0.20.1 // indirect
	k8s.io/apimachinery v0.20.1 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace (
	// github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0
	github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
	k8s.io/api => k8s.io/api v0.19.2

	k8s.io/client-go => k8s.io/client-go v0.19.2

)
