# Kubernetes Operator Demo

Reference: [Tutorial: Building CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial)

## API

### Create API

```bash
# install kubebuilder
arkade get kubebuilder

# install go 1.16
sudo snap refresh go --channel=1.16/stable

mkdir operator-demo
cd operator-demo

kubebuilder init --domain kibazen.cn --repo github.com/kibaamor/operator-demo

kubebuilder create api --group apps --version v1alpha1 --kind KDaemon --resource --controller
```

### Implement API

Edit file `api/v1alpha1/kdaemon_types.go`.

```diff
// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KDaemonSpec defines the desired state of KDaemon
type KDaemonSpec struct {
    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
    // Important: Run "make" to regenerate code after modifying this file

-    // Foo is an example field of KDaemon. Edit kdaemon_types.go to remove/update
-    Foo string `json:"foo,omitempty"`
+   // Pod image
+   Image string `json:"image,omitempty"`
}

// KDaemonStatus defines the observed state of KDaemon
type KDaemonStatus struct {
    // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
    // Important: Run "make" to regenerate code after modifying this file
+
+   // available replicas number
+   AvailableReplicas int `json:"availableReplicas,omitempty"`
}

```

Edit file `controllers/kdaemon_controller.go`.

```diff
import (
    "context"
+    "fmt"

+    v1 "k8s.io/api/core/v1"
+    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    appsv1alpha1 "github.com/kibaamor/operator-demo/api/v1alpha1"
)

...

func (r *KDaemonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
-    _ = log.FromContext(ctx)
-
-    // your logic here
+    log := log.FromContext(ctx)
+
+    kds := &appsv1alpha1.KDaemon{}
+    if err := r.Client.Get(ctx, req.NamespacedName, kds); err != nil {
+        log.Error(err, "failed to get KDaemon")
+        return ctrl.Result{
+            Requeue: true,
+        }, err
+    }
+    if kds.Spec.Image == "" {
+        err := fmt.Errorf("invalid image config")
+        log.Error(err, "can not deploy with empty image")
+        return ctrl.Result{}, err
+    }
+
+    nl := &v1.NodeList{}
+    if err := r.Client.List(ctx, nl); err != nil {
+        log.Error(err, "failed to get node list")
+        return ctrl.Result{
+            Requeue: true,
+        }, err
+    }
+
+    for _, n := range nl.Items {
+        p := &v1.Pod{
+            TypeMeta: metav1.TypeMeta{
+                APIVersion: "v1",
+                Kind:       "Pod",
+            },
+            ObjectMeta: metav1.ObjectMeta{
+                GenerateName: fmt.Sprintf("%s-", n.Name),
+                Namespace:    kds.Namespace,
+            },
+            Spec: v1.PodSpec{
+                Containers: []v1.Container{
+                    {
+                        Image: kds.Spec.Image,
+                        Name:  "kdaemon",
+                    },
+                },
+                NodeName: n.Name,
+            },
+        }
+        if err := r.Client.Create(ctx, p); err != nil {
+            log.Error(err, "failed create pod on Node", "node", n.Name)
+            return ctrl.Result{}, err
+        }
+    }

    return ctrl.Result{}, nil
}
```

### Test API

Edit file `config/samples/apps_v1alpha1_kdaemon.yaml`.

```diff
apiVersion: apps.kibazen.cn/v1alpha1
kind: KDaemon
metadata:
  name: kdaemon-sample
spec:
  # Add fields here
-  foo: bar
+  image: nginx
```

```bash
# build
go mod tidy
make manifests

# install CRD
make install

# run controller locally
make run
```

## Webhook

### Create webhook

```bash
kubebuilder create webhook --group apps --version v1alpha1 --kind KDaemon --defaulting --programmatic-validation
```

### Implement webhook

Edit file `api/v1alpha1/kdaemon_webhook.go`.

```diff
// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KDaemon) ValidateCreate() error {
    kdaemonlog.Info("validate create", "name", r.Name)

    // TODO(user): fill in your validation logic upon object creation.
+    if r.Spec.Image == "" {
+        return fmt.Errorf("image is required")
+    }
    return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KDaemon) ValidateUpdate(old runtime.Object) error {
    kdaemonlog.Info("validate update", "name", r.Name)

    // TODO(user): fill in your validation logic upon object update.
+    if r.Spec.Image == "" {
+        return fmt.Errorf("image is required")
+    }
    return nil
}
```

Update CRD.

```bash
# update CRD
make manifests
```

### Install cert-manager

Because the webhook server must provide https service to be used normally by k8s, and the easiest way to issue a certificate to webhook is to generate it through cert-manager. Thus, we need install cert-manager first.

```bash
# install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml
```

### Build docker image

Now we can proceed to the next steps.

Build docker image and load it init minikube.

```bash
# build docker image locally
make docker-build IMG=kibazen.cn/operator-demo:kdaemon-v1alpha1

# load docker image into minikube
minikube image load kibazen.cn/operator-demo:kdaemon-v1alpha1
```

### Enable webhook and cert manager

Enable the webhook and cert manager configuration through kustomize.

Edit file `config/default/kustomization.yaml` and it should look like the following:

```yaml
# Adds namespace to all resources.
namespace: operator-demo-system

# Value of this field is prepended to the
# names of all resources, e.g. a deployment named
# "wordpress" becomes "alices-wordpress".
# Note that it should also match with the prefix (text before '-') of the namespace
# field above.
namePrefix: operator-demo-

# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

bases:
- ../crd
- ../rbac
- ../manager
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
- ../certmanager
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'.
# - ../prometheus

patchesStrategicMerge:
# Protect the /metrics endpoint by putting it behind auth.
# If you want your controller-manager to expose the /metrics
# endpoint w/o any authn/z, please comment the following line.
- manager_auth_proxy_patch.yaml

# Mount the controller config file for loading manager configurations
# through a ComponentConfig type
#- manager_config_patch.yaml

# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
- manager_webhook_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
# Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
# 'CERTMANAGER' needs to be enabled to use ca injection
- webhookcainjection_patch.yaml

# the following config is for teaching kustomize how to do var substitution
vars:
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
- name: CERTIFICATE_NAMESPACE # namespace of the certificate CR
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1
    name: serving-cert # this name should match the one in certificate.yaml
  fieldref:
    fieldpath: metadata.namespace
- name: CERTIFICATE_NAME
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1
    name: serving-cert # this name should match the one in certificate.yaml
- name: SERVICE_NAMESPACE # namespace of the service
  objref:
    kind: Service
    version: v1
    name: webhook-service
  fieldref:
    fieldpath: metadata.namespace
- name: SERVICE_NAME
  objref:
    kind: Service
    version: v1
    name: webhook-service
```

Edit file `config/crd/kustomization.yamll` and it should look like the following:

```yaml
# This kustomization.yaml is not intended to be run by itself,
# since it depends on service name and namespace that are out of this kustomize package.
# It should be run by config/default
resources:
- bases/apps.kibazen.cn_kdaemons.yaml
#+kubebuilder:scaffold:crdkustomizeresource

patchesStrategicMerge:
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix.
# patches here are for enabling the conversion webhook for each CRD
- patches/webhook_in_kdaemons.yaml
#+kubebuilder:scaffold:crdkustomizewebhookpatch

# [CERTMANAGER] To enable webhook, uncomment all the sections with [CERTMANAGER] prefix.
# patches here are for enabling the CA injection for each CRD
- patches/cainjection_in_kdaemons.yaml
#+kubebuilder:scaffold:crdkustomizecainjectionpatch

# the following config is for teaching kustomize how to do kustomization for CRDs.
configurations:
- kustomizeconfig.yaml
```

### Deploy webhook

```bash
make deploy IMG=kibazen.cn/operator-demo:kdaemon-v1alpha1
```

### Test webhook

Edit file `config/samples/apps_v1alpha1_kdaemon.yaml`.

```diff
apiVersion: apps.kibazen.cn/v1alpha1
kind: KDaemon
metadata:
  name: kdaemon-sample
spec:
  # Add fields here
-  image: nginx
```

#### Test update

```bash
$ kubectl apply -f ./config/samples/apps_v1alpha1_kdaemon.yaml
Error from server (image is required): error when applying patch:
{"metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"apps.kibazen.cn/v1alpha1\",\"kind\":\"KDaemon\",\"metadata\":{\"annotations\":{},\"name\":\"kdaemon-sample\",\"namespace\":\"default\"},\"spec\":null}\n"}},"spec":null}
to:
Resource: "apps.kibazen.cn/v1alpha1, Resource=kdaemons", GroupVersionKind: "apps.kibazen.cn/v1alpha1, Kind=KDaemon"
Name: "kdaemon-sample", Namespace: "default"
for: "./config/samples/apps_v1alpha1_kdaemon.yaml": error when patching "./config/samples/apps_v1alpha1_kdaemon.yaml": admission webhook "vkdaemon.kb.io" denied the request: image is required
```

#### Test create

```bash
# delete it before create
$ kubectl delete -f ./config/samples/apps_v1alpha1_kdaemon.yaml

$ kubectl create -f ./config/samples/apps_v1alpha1_kdaemon.yaml
Error from server (image is required): error when creating "./config/samples/apps_v1alpha1_kdaemon.yaml": admission webhook "vkdaemon.kb.io" denied the request: image is required
```
