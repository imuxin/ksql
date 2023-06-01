# ksql, a SQL-like language tool for kubernetes

## Quick start

install

```bash
go install github.com/imuxin/ksql
```

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

	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    "k8s.io",
		Version:  "v1alpha1",
		Resource: "tttt",
	}

	s := labels.NewSelector()
	req, err := labels.NewRequirement("key", selection.Equals, []string{"val"})
	if err != nil {
		return nil, err
	}
	s = s.Add(*req)

	us, err := client.Resource(gvr).List(context.TODO(), metav1.ListOptions{
		LabelSelector: s.String(),
	})
	if err != nil {
		return nil, err
	}

	var results []T
	for _, item := range us.Items {
		obj := &T{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
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
func list() ([]T, error)  {
	
}
```
</td>
</tr>
</table>

