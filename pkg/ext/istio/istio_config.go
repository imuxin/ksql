package istio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/ext/kube"
	utilkube "github.com/imuxin/ksql/pkg/util/kube"
)

const (
	IstioSystem = "istio-system"
)

var _ abs.Plugin = &IstioConfig{}

// nolint
type IstioConfig struct {
	APIServerPlugin kube.APIServerPlugin
}

// nolint
func NewPlugin(
	table, namespace string,
	names []string,
	restConfig *rest.Config,
	selector labels.Selector) abs.Plugin {
	return IstioConfig{
		APIServerPlugin: kube.NewPlugin(table, namespace, names, restConfig, selector).(kube.APIServerPlugin),
	}
}

func (t IstioConfig) todo() ([]abs.Object, error) {
	config, err := t.APIServerPlugin.RestConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// istio=pilot
	pods, err := clientset.CoreV1().Pods(IstioSystem).List(context.TODO(),
		v1.ListOptions{
			LabelSelector: "app=istiod",
			FieldSelector: "status.phase=Running",
		})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Pods: %v", err)
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no running Istio pods in %q", IstioSystem)
	}
	pf, err := utilkube.NewPortForwarder(config, pods.Items[0].Name, pods.Items[0].Namespace, "", 0, 8080)
	if err != nil {
		return nil, fmt.Errorf("unable to create portforward: %v", err)
	}
	if err := pf.Start(); err != nil {
		return nil, err
	}
	defer pf.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/%s", pf.Address(), "debug/configz"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []abs.Object
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (t IstioConfig) Download() ([]abs.Object, error) {
	return t.todo()
}

func (t IstioConfig) Delete(list []abs.Object) ([]abs.Object, error) { return nil, nil } //nolint
