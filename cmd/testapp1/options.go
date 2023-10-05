package main

import (
	"errors"

	"github.com/spf13/pflag"

	"k8c.io/kcp-ctrl-runtime-test/pkg/log"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type Options struct {
	// NB: Not actually defined here, as ctrl-runtime registers its
	// own --kubeconfig flag that is required to make its GetConfigOrDie()
	// work.
	// KubeconfigFile string

	// KCPKubeconfig is the kubeconfig that gives access to kcp.
	KCPKubeconfig string

	LogOptions log.Options
}

func NewOptions() *Options {
	return &Options{
		LogOptions: log.NewDefaultOptions(),
	}
}

func (o *Options) AddFlags(flags *pflag.FlagSet) {
	o.LogOptions.AddPFlags(flags)
	flags.StringVar(&o.KCPKubeconfig, "kcp-kubeconfig", o.KCPKubeconfig, "kubeconfig file for kcp")
}

func (o *Options) Validate() error {
	errs := []error{}

	if err := o.LogOptions.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(o.KCPKubeconfig) == 0 {
		errs = append(errs, errors.New("--kcp-kubeconfig is required"))
	}

	return utilerrors.NewAggregate(errs)
}
