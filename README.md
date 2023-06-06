# ksql, a SQL-like language tool for kubernetes

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/887f9700e424478e9be6dc88237e2f72)](https://app.codacy.com/gh/imuxin/ksql/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/887f9700e424478e9be6dc88237e2f72)](https://app.codacy.com/gh/imuxin/ksql/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

## Install

```bash
go install github.com/imuxin/ksql
```

## Goal #1: bring SQL lanugage for kubernetes command line tool

Examples:

- `SELECT`

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

- `DESC`

```go
ksql> DESC envoyfilters.networking.istio.io
+-------------------------------------------------------------------------------------------------+----------+
| SCHEMA                                                                                          | VERSION  |
+-------------------------------------------------------------------------------------------------+----------+
| type EnvoyFilter struct {                                                                       | v1alpha3 |
|     // Customizing Envoy configuration generated by Istio. See more details at:                 |          |
|     // https://istio.io/docs/reference/config/networking/envoy-filter.html                      |          |
|     spec struct {                                                                               |          |
|         // One or more patches with match conditions.                                           |          |
|         configPatches []struct {                                                                |          |
|             applyTo string                                                                      |          |
|             // Match on listener/route configuration/cluster.                                   |          |
|             match struct {                                                                      |          |
|                 // Match on envoy cluster attributes.                                           |          |
|                 cluster struct {                                                                |          |
|                     // The exact name of the cluster to match.                                  |          |
|                     name string                                                                 |          |
|                     // The service port for which this cluster was generated.                   |          |
|                     portNumber int                                                              |          |
|                     // The fully qualified service name for this cluster.                       |          |
|                     service string                                                              |          |
|                     // The subset associated with the service.                                  |          |
|                     subset string                                                               |          |
|                 }                                                                               |          |
|                 // The specific config generation context to match on.                          |          |
|                 context string                                                                  |          |
|                 // Match on envoy listener attributes.                                          |          |
|                 listener struct {                                                               |          |
|                     // Match a specific filter chain in a listener.                             |          |
|                     filterChain struct {                                                        |          |
|                         // Applies only to sidecars.                                            |          |
|                         applicationProtocols string                                             |          |
|                         // The destination_port value used by a filter chain's match condition. |          |
|                         destinationPort int                                                     |          |
|                         // The name of a specific filter to apply the patch to.                 |          |
|                         filter struct {                                                         |          |
|                             // The filter name to match on.                                     |          |
|                             name string                                                         |          |
|                             subFilter struct {                                                  |          |
|                                 // The filter name to match on.                                 |          |
|                                 name string                                                     |          |
|                             }                                                                   |          |
|                         }                                                                       |          |
|                         // The name assigned to the filter chain.                               |          |
|                         name string                                                             |          |
|                         // The SNI value used by a filter chain's match condition.              |          |
|                         sni string                                                              |          |
|                         // Applies only to `SIDECAR_INBOUND` context.                           |          |
|                         transportProtocol string                                                |          |
|                     }                                                                           |          |
|                     // Match a specific listener by its name.                                   |          |
|                     name string                                                                 |          |
|                     portName string                                                             |          |
|                     portNumber int                                                              |          |
|                 }                                                                               |          |
|                 // Match on properties associated with a proxy.                                 |          |
|                 proxy struct {                                                                  |          |
|                     metadata map[string]string                                                  |          |
|                     proxyVersion string                                                         |          |
|                 }                                                                               |          |
|                 // Match on envoy HTTP route configuration attributes.                          |          |
|                 routeConfiguration struct {                                                     |          |
|                     gateway string                                                              |          |
|                     // Route configuration name to match on.                                    |          |
|                     name string                                                                 |          |
|                     // Applicable only for GATEWAY context.                                     |          |
|                     portName string                                                             |          |
|                     portNumber int                                                              |          |
|                     vhost struct {                                                              |          |
|                         name string                                                             |          |
|                         // Match a specific route within the virtual host.                      |          |
|                         route struct {                                                          |          |
|                             // Match a route with specific action type.                         |          |
|                             action string                                                       |          |
|                             name string                                                         |          |
|                         }                                                                       |          |
|                     }                                                                           |          |
|                 }                                                                               |          |
|             }                                                                                   |          |
|             // The patch to apply along with the operation.                                     |          |
|             patch struct {                                                                      |          |
|                 // Determines the filter insertion order.                                       |          |
|                 filterClass string                                                              |          |
|                 // Determines how the patch should be applied.                                  |          |
|                 operation string                                                                |          |
|                 // The JSON config of the object being patched.                                 |          |
|                 value map[string]interface{}                                                    |          |
|             }                                                                                   |          |
|         }                                                                                       |          |
|         // Priority defines the order in which patch sets are applied within a context.         |          |
|         priority int                                                                            |          |
|         workloadSelector struct {                                                               |          |
|             labels map[string]string                                                            |          |
|         }                                                                                       |          |
|     }                                                                                           |          |
|     status map[string]interface{}                                                               |          |
| }                                                                                               |          |
+-------------------------------------------------------------------------------------------------+----------+
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

## Roadmap

- [x] Support `SELECT` stat
- [x] Support `FROM`
- [x] Support `AS` `LABEL` `NAMESPACE` `NAME`
- [x] Support `WHERE` expr
- [x] Support `DESC` expr
- [ ] Support `USE` stat
- [ ] Support `DELETE` stat
- [ ] Support `UPDATE` stat
- [ ] Support custom TABLE extensions
- [ ] ...
