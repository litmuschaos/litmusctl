module github.com/mayadata-io/kuberactl

go 1.14

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-resty/resty/v2 v2.3.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.3
	k8s.io/cli-runtime v0.19.2 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubectl v0.19.2 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.19.2

	k8s.io/client-go => k8s.io/client-go v0.19.2

)
