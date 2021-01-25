module github.com/mayadata-io/kuberactl

go 1.14

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-resty/resty/v2 v2.3.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mayadata-io/cli-utils v0.0.0-20210119141112-84fe44fff4e7
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	k8s.io/apimachinery v0.20.1 // indirect
	k8s.io/cli-runtime v0.20.1 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubectl v0.19.2 // indirect
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
