package main

import (
	"fmt"
	"flag"
	"net"
	"net/http"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getRemoteIP(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

func checkNodeHealth(node *v1.Node) (bool) {
	for _, t := range node.Spec.Taints {
		if t.Key == v1.TaintNodeNotReady {
			return false
			}
		if t.Key == v1.TaintNodeUnreachable {
			return false
		}
		if t.Key == v1.TaintNodeUnschedulable {
			return false
		}
		if t.Key == v1.TaintNodeNetworkUnavailable {
			return false
		}
	}
	return true
}

func main() {
	var kubeconfig *string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	incluster := flag.Bool("incluster", false, "use incluster config")
	flag.Parse()

	var config *rest.Config
	var err error
	if *incluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		hostname := req.URL.Query()["host"]
		if len(hostname) != 1 {
			http.Error(w, "Missing \"host\" parameter", 400)
			return
		}

		node, err := clientset.CoreV1().Nodes().Get(hostname[0], metav1.GetOptions{})
		if errors.IsNotFound(err) {
			error := fmt.Sprintf("Node not found: %s", hostname[0])
			http.Error(w, error, 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		healthy := checkNodeHealth(node)
		if (healthy) {
			fmt.Fprintf(w, "Node is healthy: %s\n", node.Name)
		} else {
			error := fmt.Sprintf("Node is NOT healthy: %s", node.Name)
			http.Error(w, error, http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8090", nil)
}
