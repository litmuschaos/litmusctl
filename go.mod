module github.com/litmuschaos/litmusctl

go 1.14

require (
	github.com/argoproj/argo v2.5.2+incompatible
	github.com/go-resty/resty/v2 v2.3.0
	github.com/spf13/cobra v1.1.1
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
