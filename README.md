# Kubernetes Operator Demo

practice for <https://github.com/kibaamor/101/blob/master/module11/operator/kubebuilder.md>

```bash
# install kubebuilder
$ arkade get kubebuilder

# install go 1.16
$ sudo snap refresh go --channel=1.16/stable

$ mkdir operator-demo

$ cd operator-demo

$ kubebuilder init --domain kibazen.cn --repo github.com/kibaamor/operator-demo

$ kubebuilder create api --group apps --version v1alpha1 --kind KDaemon --resource --controller
```

Edit file `api/v1alpha1/kdaemon_types.go`

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

Edit file `controllers/kdaemon_controller.go`

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

Edit file `config/samples/apps_v1alpha1_kdaemon.yaml`

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
$ go mod tidy
$ make

# install CRD
$ make install

# run controller locally
$ make run
```
