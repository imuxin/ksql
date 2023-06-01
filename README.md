# ksql, a SQL-like language tool for kubernetes

## Install

```bash
go install github.com/imuxin/ksql
```

## Goal #1: bring SQL lanugage for kubernetes command line tool

Simple examples:

```bash
ksql> SELECT * FROM service NAMESPACE default NAME kubernetes
+------------+-----------+
| NAME       | NAMESPACE |
+------------+-----------+
| kubernetes | default   |
+------------+-----------+
ksql> SELECT "{ .metadata.name }" AS NAME, "{ .spec.clusterIP }" AS "CLUSTER-IP", "{ .spec.ports }" FROM svc NAMESPACE default NAME kubernetes
+------------+------------+------------------------------------------------------------------+
| NAME       | CLUSTER-IP | { .SPEC.PORTS }                                                  |
+------------+------------+------------------------------------------------------------------+
| kubernetes | 10.8.0.1   | [{"name":"https","port":443,"protocol":"TCP","targetPort":6443}] |
+------------+------------+------------------------------------------------------------------+
```

more usages, see EBNF description:
https://github.com/imuxin/ksql/blob/86e62709a6f3f1d7d6da94b02232623b8df04426/pkg/parser/parser_test.go#L11-L23

## Goal #2: make code `easier` to maintain
<table>
<tr>
<th><code>client-go</code></th>
<th><code>ksql</code></th>
</tr>
<tr>
<td>

```go
func list() ([]T, error) {
    kubeConfig := getKubeConfig()

    client, err := dynamic.
        NewForConfig(kubeConfig)
    if err != nil {
        return nil, err
    }

    gvr :=
        schema.GroupVersionResource{
            Group:    "k8s.io",
            Version:  "v1alpha1",
            Resource: "tttt",
        }

    s := labels.NewSelector()
    req, err := labels.NewRequirement(
        "key", selection.Equals, []string{"val"})
    if err != nil {
        return nil, err
    }
    s = s.Add(*req)

    us, err := client.
        Resource(gvr).
        List(context.TODO(), metav1.ListOptions{
            LabelSelector: s.String(),
        })
    if err != nil {
        return nil, err
    }

    var results []T
    for _, item := range us.Items {
        obj := &T{}
        if err := runtime.DefaultUnstructuredConverter.
            FromUnstructured(item.Object, obj); err != nil {
            return nil, err
        }
        results = append(results, *obj)
    }
    return results, nil
}
```
</td>
<td>

```go
import "github.com/imuxin/ksql/pkg/executor"

func list() ([]T, error) {
    kubeConfig := getKubeConfig()
    sql := `SELECT * FROM tttt.v1alpha1.k8s.io LABEL key = val`
    return executor.Execute[T](sql, kubeConfig)
}
```
</td>
</tr>
</table>
