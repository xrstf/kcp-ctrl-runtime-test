# kcp controller-runtime test

This sample app demonstrates how the controller-runtime fork always uses a kcp-aware cache, even if not explicitly configured. This causes the underlying cache to always include the logicalcluster name in cache keys, but since the clients do not modify any Go contexts to inject logicalcluster names, the clients will never find objects in the cache and so whenever the controller attempts to create an APIResourceSchema in kcp, it results in a HTTP 409 Conflict error.

## How to use

1. Create a local Kubernetes cluster using kind:

   ```bash
   export KUBECONFIG=kind.kubeconfig
   kind create cluster
   ```

2. Start a local kcp instance:

   ```bash
   cd ~/your/kcp/clone
   kcp start &
   cd -
   cp ~/your/kcp/clone/.kcp/admin.kubeconfig kcp.kubeconfig
   ```

3. Ensure the `kcp.kubeconfig` points to a pre-existing workspace (of type Universal), e.g. by modifying the server URL to `https://..../clusters/root:org1`.

4. Build and run the test app:

   ```bash
   make build run
   ```

## Example Log

```
_build/testapp1 \
  --kubeconfig kind.kubeconfig \
  --kcp-kubeconfig kcp.kubeconfig \
  -v 6
{"level":"info","time":"2023-10-05T14:51:35.557+0200","caller":"testapp1/main.go:57","msg":"Moin!"}
I1005 14:51:35.558383 2154146 loader.go:373] Config loaded from file:  kind.kubeconfig
{"level":"info","time":"2023-10-05T14:51:35.558+0200","caller":"testapp1/main.go:63","msg":"--kubeconfig (manager) info","host":"https://127.0.0.1:34177"}
I1005 14:51:35.563715 2154146 round_trippers.go:553] GET https://127.0.0.1:34177/api?timeout=32s 200 OK in 4 milliseconds
I1005 14:51:35.564736 2154146 round_trippers.go:553] GET https://127.0.0.1:34177/apis?timeout=32s 200 OK in 0 milliseconds
{"level":"info","time":"2023-10-05T14:51:35.565+0200","logger":"controller-runtime.metrics","caller":"manager/manager.go:413","msg":"Metrics server is starting to listen","addr":":8080"}
I1005 14:51:35.566350 2154146 loader.go:373] Config loaded from file:  kcp.kubeconfig
{"level":"info","time":"2023-10-05T14:51:35.566+0200","caller":"testapp1/main.go:86","msg":"--kcp-kubeconfig (cluster) info","host":"https://192.168.178.62:6443/clusters/root:org1"}
I1005 14:51:35.568759 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/api?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.569216 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis?timeout=32s 200 OK in 0 milliseconds
I1005 14:51:35.571363 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/authentication.k8s.io/v1?timeout=32s 200 OK in 0 milliseconds
I1005 14:51:35.572083 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/topology.kcp.io/v1alpha1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572566 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/authorization.k8s.io/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572604 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/rbac.authorization.k8s.io/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572614 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apps/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572779 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/certificates.k8s.io/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572941 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apiextensions.k8s.io/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.572958 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/core.kcp.io/v1alpha1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.573050 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1?timeout=32s 200 OK in 3 milliseconds
I1005 14:51:35.573170 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/api/v1?timeout=32s 200 OK in 3 milliseconds
I1005 14:51:35.573210 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/events.k8s.io/v1?timeout=32s 200 OK in 3 milliseconds
I1005 14:51:35.573287 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/coordination.k8s.io/v1?timeout=32s 200 OK in 3 milliseconds
I1005 14:51:35.573255 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apiresource.kcp.io/v1alpha1?timeout=32s 200 OK in 1 milliseconds
I1005 14:51:35.573486 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/admissionregistration.k8s.io/v1?timeout=32s 200 OK in 2 milliseconds
I1005 14:51:35.573503 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/tenancy.kcp.io/v1alpha1?timeout=32s 200 OK in 3 milliseconds
I1005 14:51:35.581547 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/admissionregistration.k8s.io/v1alpha1?timeout=32s 200 OK in 0 milliseconds
{"level":"info","time":"2023-10-05T14:51:35.581+0200","caller":"testapp1/main.go:111","msg":"Starting app…"}
I1005 14:51:35.581955 2154146 shared_informer.go:330] caches populated
I1005 14:51:35.581984 2154146 shared_informer.go:330] caches populated
{"level":"info","time":"2023-10-05T14:51:35.581+0200","caller":"runtime/asm_amd64.s:1650","msg":"Starting server","path":"/metrics","kind":"metrics","addr":"[::]:8080"}
{"level":"info","time":"2023-10-05T14:51:35.582+0200","caller":"controller/controller.go:242","msg":"Starting EventSource","controller":"test-controller","source":"kind source: *v1.ConfigMap"}
{"level":"info","time":"2023-10-05T14:51:35.582+0200","caller":"controller/controller.go:242","msg":"Starting Controller","controller":"test-controller"}
I1005 14:51:35.582165 2154146 reflector.go:273] Starting reflector *v1.ConfigMap (10h20m7.171066157s) from k8s.io/client-go@v12.0.0+incompatible/tools/cache/reflector.go:217
I1005 14:51:35.582171 2154146 reflector.go:309] Listing and watching *v1.ConfigMap from k8s.io/client-go@v12.0.0+incompatible/tools/cache/reflector.go:217
I1005 14:51:35.583066 2154146 round_trippers.go:553] GET https://127.0.0.1:34177/api/v1/configmaps?limit=500&resourceVersion=0 200 OK in 0 milliseconds
I1005 14:51:35.583909 2154146 round_trippers.go:553] GET https://127.0.0.1:34177/api/v1/configmaps?allowWatchBookmarks=true&resourceVersion=5429&timeoutSeconds=383&watch=true 200 OK in 0 milliseconds
I1005 14:51:35.682224 2154146 shared_informer.go:330] caches populated
I1005 14:51:35.682276 2154146 shared_informer.go:330] caches populated
{"level":"info","time":"2023-10-05T14:51:35.682+0200","caller":"controller/controller.go:242","msg":"Starting workers","controller":"test-controller","worker count":4}
{"level":"info","time":"2023-10-05T14:51:35.682+0200","logger":"test-controller","caller":"testctrl/controller.go:93","msg":"Processing","configmap":"default/kube-root-ca.crt"}
{"level":"info","time":"2023-10-05T14:51:35.682+0200","logger":"test-controller","caller":"testctrl/controller.go:93","msg":"Processing","configmap":"kube-system/extension-apiserver-authentication"}
{"level":"info","time":"2023-10-05T14:51:35.682+0200","logger":"test-controller","caller":"testctrl/controller.go:93","msg":"Processing","configmap":"kube-system/kube-proxy"}
{"level":"info","time":"2023-10-05T14:51:35.682+0200","logger":"test-controller","caller":"testctrl/controller.go:93","msg":"Processing","configmap":"kube-node-lease/kube-root-ca.crt"}
I1005 14:51:35.682542 2154146 reflector.go:273] Starting reflector *v1alpha1.APIResourceSchema (9h11m52.137201308s) from k8s.io/client-go@v12.0.0+incompatible/tools/cache/reflector.go:217
I1005 14:51:35.682549 2154146 reflector.go:309] Listing and watching *v1alpha1.APIResourceSchema from k8s.io/client-go@v12.0.0+incompatible/tools/cache/reflector.go:217
I1005 14:51:35.685253 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas?limit=500&resourceVersion=0 200 OK in 2 milliseconds
I1005 14:51:35.688092 2154146 round_trippers.go:553] GET https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas?allowWatchBookmarks=true&resourceVersion=16653&timeoutSeconds=516&watch=true 200 OK in 0 milliseconds
I1005 14:51:35.783339 2154146 shared_informer.go:330] caches populated
I1005 14:51:35.783385 2154146 shared_informer.go:330] caches populated
{"level":"info","time":"2023-10-05T14:51:35.783+0200","logger":"test-controller","caller":"testctrl/controller.go:125","msg":"Creating new ARS, cause could not fetch existing one","configmap":"kube-system/extension-apiserver-authentication","error":"APIResourceSchema.apis.kcp.io \"v42.foos.tremors.valley\" not found"}
I1005 14:51:35.783384 2154146 shared_informer.go:330] caches populated
{"level":"info","time":"2023-10-05T14:51:35.783+0200","logger":"test-controller","caller":"testctrl/controller.go:125","msg":"Creating new ARS, cause could not fetch existing one","configmap":"default/kube-root-ca.crt","error":"APIResourceSchema.apis.kcp.io \"v42.foos.tremors.valley\" not found"}
I1005 14:51:35.783390 2154146 shared_informer.go:330] caches populated
{"level":"info","time":"2023-10-05T14:51:35.783+0200","logger":"test-controller","caller":"testctrl/controller.go:125","msg":"Creating new ARS, cause could not fetch existing one","configmap":"kube-node-lease/kube-root-ca.crt","error":"APIResourceSchema.apis.kcp.io \"v42.foos.tremors.valley\" not found"}
{"level":"info","time":"2023-10-05T14:51:35.783+0200","logger":"test-controller","caller":"testctrl/controller.go:125","msg":"Creating new ARS, cause could not fetch existing one","configmap":"kube-system/kube-proxy","error":"APIResourceSchema.apis.kcp.io \"v42.foos.tremors.valley\" not found"}
I1005 14:51:35.804110 2154146 round_trippers.go:553] POST https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas 409 Conflict in 20 milliseconds
I1005 14:51:35.804141 2154146 round_trippers.go:553] POST https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas 409 Conflict in 20 milliseconds
I1005 14:51:35.804152 2154146 round_trippers.go:553] POST https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas 409 Conflict in 20 milliseconds
I1005 14:51:35.804199 2154146 round_trippers.go:553] POST https://192.168.178.62:6443/clusters/root:org1/apis/apis.kcp.io/v1alpha1/apiresourceschemas 409 Conflict in 20 milliseconds
{"level":"info","time":"2023-10-05T14:51:35.804+0200","logger":"test-controller","caller":"testctrl/controller.go:148","msg":"result of creating ARS","configmap":"kube-node-lease/kube-root-ca.crt","error":"apiresourceschemas.apis.kcp.io \"v42.foos.tremors.valley\" already exists"}
{"level":"info","time":"2023-10-05T14:51:35.804+0200","logger":"test-controller","caller":"testctrl/controller.go:148","msg":"result of creating ARS","configmap":"default/kube-root-ca.crt","error":"apiresourceschemas.apis.kcp.io \"v42.foos.tremors.valley\" already exists"}
{"level":"info","time":"2023-10-05T14:51:35.804+0200","logger":"test-controller","caller":"testctrl/controller.go:93","msg":"Processing","configmap":"kube-system/kube-root-ca.crt"}
```
