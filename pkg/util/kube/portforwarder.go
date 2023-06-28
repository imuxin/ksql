// Copied from https://github.com/istio/istio/blob/master/pkg/kube/portforwarder.go

package kube

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	defaultLocalAddress = "localhost"
)

// PortForwarder manages the forwarding of a single port.
type PortForwarder interface {
	// Start runs this forwarder.
	Start() error

	// Address returns the local forwarded address. Only valid while the forwarder is running.
	Address() string

	// Close this forwarder and release an resources.
	Close()

	// WaitForStop blocks until connection closed (e.g. control-C interrupt)
	WaitForStop()
}

var _ PortForwarder = &forwarder{}

type forwarder struct {
	stopCh       chan struct{}
	restConfig   *rest.Config
	podName      string
	ns           string
	localAddress string
	localPort    int
	podPort      int
	address      string
}

func (f *forwarder) Start() error {
	errCh := make(chan error, 1)
	readyCh := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-f.stopCh:
				return
			default:
			}
			// Build a new port forwarder.
			fw, err := f.buildK8sPortForwarder(readyCh)
			if err != nil {
				errCh <- err
				return
			}
			if err = fw.ForwardPorts(); err != nil {
				errCh <- err
				return
			}
			// At this point, either the stopCh has been closed, or port forwarder connection is broken.
			// the port forwarder should have already been ready before.
			// No need to notify the ready channel anymore when forwarding again.
			readyCh = nil
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failure running port forward process: %v", err)
	case <-readyCh:
		// The forwarder is now ready.
		return nil
	}
}

func (f *forwarder) Address() string {
	return f.address
}

func (f *forwarder) Close() {
	close(f.stopCh)
	// Closing the stop channel should close anything
	// opened by f.forwarder.ForwardPorts()
}

func (f *forwarder) WaitForStop() {
	<-f.stopCh
}

func NewPortForwarder(config *rest.Config, podName, ns, localAddress string, localPort, podPort int) (PortForwarder, error) {
	if localAddress == "" {
		localAddress = defaultLocalAddress
	}
	f := &forwarder{
		stopCh:       make(chan struct{}),
		restConfig:   config,
		podName:      podName,
		ns:           ns,
		localAddress: localAddress,
		localPort:    localPort,
		podPort:      podPort,
	}
	if f.localPort == 0 {
		var err error
		reserved, err := defaultAvailablePort()
		if err != nil {
			return nil, fmt.Errorf("failure allocating port: %v", err)
		}
		f.localPort = int(reserved)
	}

	sLocalPort := fmt.Sprintf("%d", f.localPort)
	f.address = net.JoinHostPort(localAddress, sLocalPort)
	return f, nil
}

func (f *forwarder) buildK8sPortForwarder(readyCh chan struct{}) (*portforward.PortForwarder, error) {
	f.restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	f.restConfig.GroupVersion = &schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	f.restConfig.APIPath = "/api"
	restClient, err := rest.RESTClientFor(f.restConfig)
	if err != nil {
		return nil, err
	}

	req := restClient.Post().Resource("pods").Namespace(f.ns).Name(f.podName).SubResource("portforward")
	serverURL := req.URL()

	roundTripper, upgrader, err := roundTripperFor(f.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failure creating roundtripper: %v", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, serverURL)

	fw, err := portforward.NewOnAddresses(dialer,
		[]string{f.localAddress},
		[]string{fmt.Sprintf("%d:%d", f.localPort, f.podPort)},
		f.stopCh,
		readyCh,
		io.Discard,
		os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed establishing port-forward: %v", err)
	}

	// Run the same check as k8s.io/kubectl/pkg/cmd/portforward/portforward.go
	// so that we will fail early if there is a problem contacting API server.
	podGet := restClient.Get().Resource("pods").Namespace(f.ns).Name(f.podName)
	obj, err := podGet.Do(context.TODO()).Get()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving: %v in the %q namespace", err, f.ns)
	}
	us, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed getting pod, object type is %T", obj)
	}

	pod := &corev1.Pod{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(us.Object, pod); err != nil {
		return nil, fmt.Errorf("failed converting pod: %v", err)
	}

	if pod.Status.Phase != corev1.PodRunning {
		return nil, fmt.Errorf("pod is not running. Status=%v", pod.Status.Phase)
	}

	return fw, nil
}

func defaultAvailablePort() (uint16, error) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("127.0.0.1", "0"))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	return uint16(port), l.Close()
}
