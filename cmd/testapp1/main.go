package main

import (
	"context"
	"flag"
	"fmt"
	golog "log"

	"github.com/go-logr/zapr"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"k8c.io/kcp-ctrl-runtime-test/pkg/controller/testctrl"
	kdplog "k8c.io/kcp-ctrl-runtime-test/pkg/log"

	kcpdevv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/apis/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	ctrlruntimelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	klog.InitFlags(flag.CommandLine)

	ctx := context.Background()

	opts := NewOptions()
	opts.AddFlags(pflag.CommandLine)

	// ctrl-runtime will have added its --kubeconfig to Go's flag set
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := opts.Validate(); err != nil {
		golog.Fatalf("Invalid command line: %v", err)
	}

	log := kdplog.NewFromOptions(opts.LogOptions)
	sugar := log.Sugar()

	// set the logger used by sigs.k8s.io/controller-runtime
	ctrlruntimelog.SetLogger(zapr.NewLogger(log.WithOptions(zap.AddCallerSkip(1))))

	if err := run(ctx, sugar, opts); err != nil {
		sugar.Fatalw("Servlet has encountered an error", zap.Error(err))
	}
}

func run(ctx context.Context, log *zap.SugaredLogger, opts *Options) error {
	log.Info("Moin!")

	localScheme := runtime.NewScheme()
	kcpScheme := runtime.NewScheme()

	kubeconfig := ctrlruntime.GetConfigOrDie()
	log.Infow("--kubeconfig (manager) info", "host", kubeconfig.Host)

	// create the ctrl-runtime manager (this points to a kind cluster)
	mgr, err := manager.New(kubeconfig, manager.Options{
		Scheme: localScheme,
		BaseContext: func() context.Context {
			return ctx
		},
	})
	if err != nil {
		return fmt.Errorf("failed to setup manager: %w", err)
	}

	if err := corev1.AddToScheme(localScheme); err != nil {
		return fmt.Errorf("failed to register local scheme %s: %w", corev1.SchemeGroupVersion, err)
	}

	// init the kcp cluster connection
	kcpRestConfig, err := loadKubeconfig(opts.KCPKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load platform kubeconfig: %w", err)
	}

	log.Infow("--kcp-kubeconfig (cluster) info", "host", kcpRestConfig.Host)

	if err := kcpdevv1alpha1.AddToScheme(kcpScheme); err != nil {
		return fmt.Errorf("failed to register kcp scheme %s: %w", kcpdevv1alpha1.SchemeGroupVersion, err)
	}

	kcpCluster, err := cluster.New(kcpRestConfig, func(o *cluster.Options) {
		o.Scheme = kcpScheme
		// o.MapperProvider = kcp.NewClusterAwareMapperProvider
		// o.NewClient = kcp.NewClusterAwareClient
		// o.NewCache = kcp.NewClusterAwareCache
		// o.NewAPIReader = kcp.NewClusterAwareAPIReader
	})
	if err != nil {
		return fmt.Errorf("failed to initialize kcp cluster: %w", err)
	}

	if err := mgr.Add(kcpCluster); err != nil {
		return fmt.Errorf("failed to add kcp cluster to manager: %w", err)
	}

	if err := testctrl.Add(mgr, kcpCluster, log, 4); err != nil {
		return fmt.Errorf("failed to add test controller: %w", err)
	}

	log.Info("Starting app…")

	return mgr.Start(ctx)
}

func loadKubeconfig(filename string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = filename

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil).ClientConfig()
}
