package testctrl

import (
	"context"
	"fmt"
	"strings"

	kcpdevv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/apis/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/kontext"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName = "test-controller"
)

type Reconciler struct {
	localClient        ctrlruntimeclient.Client
	kcpCluster         cluster.Cluster
	log                *zap.SugaredLogger
	recorder           record.EventRecorder
	logicalClusterName logicalcluster.Name
}

var kcpAwareEnqueueRequestForObject = handler.EnqueueRequestsFromMapFunc(func(o ctrlruntimeclient.Object) []reconcile.Request {
	annotations := o.GetAnnotations()
	cluster := annotations["kcp.io/cluster"]

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: fmt.Sprintf("%s!%s", cluster, o.GetNamespace()),
			Name:      o.GetName(),
		},
	}}
})

func decodeKcpAwareRequest(req reconcile.Request) (types.NamespacedName, logicalcluster.Name) {
	parts := strings.SplitN(req.Namespace, "!", 2)
	cluster := parts[0]
	namespace := parts[1]

	return types.NamespacedName{
		Namespace: namespace,
		Name:      req.Name,
	}, logicalcluster.Name(cluster)
}

// Add creates a new controller and adds it to the given manager.
func Add(
	mgr manager.Manager,
	kcpCluster cluster.Cluster,
	logicalClusterName logicalcluster.Name,
	log *zap.SugaredLogger,
	numWorkers int,
) error {
	reconciler := &Reconciler{
		localClient:        mgr.GetClient(),
		kcpCluster:         kcpCluster,
		log:                log.Named(ControllerName),
		recorder:           mgr.GetEventRecorderFor(ControllerName),
		logicalClusterName: logicalClusterName,
	}

	ctrlOptions := controller.Options{
		Reconciler:              reconciler,
		MaxConcurrentReconciles: numWorkers,
	}
	c, err := controller.New(ControllerName, mgr, ctrlOptions)
	if err != nil {
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.With("configmap", request)
	log.Info("Processing")

	cm := &corev1.ConfigMap{}
	if err := r.localClient.Get(ctx, request.NamespacedName, cm); err != nil {
		return reconcile.Result{}, err // normally we'd ignore NotFound errs here, but for testing we do not
	}

	if cm.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	result, err := r.reconcile(ctx, log)
	if err != nil {
		r.recorder.Event(cm, corev1.EventTypeWarning, "ReconcilingError", err.Error())
	}
	if result == nil {
		result = &reconcile.Result{}
	}

	return *result, err
}

func (r *Reconciler) reconcile(ctx context.Context, log *zap.SugaredLogger) (*reconcile.Result, error) {
	arsName := "v42.foos.tremors.valley"
	wsCtx := kontext.WithCluster(ctx, r.logicalClusterName)
	kcpClient := r.kcpCluster.GetClient()

	ars := &kcpdevv1alpha1.APIResourceSchema{}
	if err := kcpClient.Get(wsCtx, types.NamespacedName{Name: arsName}, ars); err != nil {
		log.Infow("Creating new ARS, cause could not fetch existing one", zap.Error(err))

		// create the ARS if it was missing
		ars.Name = arsName
		ars.Spec.Group = "tremors.valley"
		ars.Spec.Scope = apiextensionsv1.NamespaceScoped
		ars.Spec.Names = apiextensionsv1.CustomResourceDefinitionNames{
			Plural:   "foos",
			Singular: "foo",
			Kind:     "Foo",
			ListKind: "FooList",
		}
		ars.Spec.Versions = []kcpdevv1alpha1.APIResourceVersion{
			{
				Name:    "v1",
				Storage: true,
				Served:  true,
				Schema: runtime.RawExtension{
					Raw: []byte(`{"type":"object","properties":{}}`),
				},
			},
		}

		log.Infow("result of creating ARS", zap.Error(kcpClient.Create(wsCtx, ars)))
	} else {
		log.Info("Found ARS, nothing to do.")
	}

	return nil, nil
}
