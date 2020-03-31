package main

import (
	"fmt"
	"flag"
	"log"
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
	// parse arguments
	var kubeconfig *string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	incluster := flag.Bool("incluster", false, "use incluster config")
	addr := flag.String("addr", ":8991", "Address to listen on")
	nodeName := flag.String("node", os.Getenv("NODE_NAME"), "(optional) node name to check")
	sickCode := flag.Int("sick-code", 404, "response code used when node is sick")
	flag.Parse()

	// load kubeconfig
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

	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// determine host to check
		host := *nodeName
		hostParam := req.URL.Query()["host"]
		if len(hostParam) == 1 {
			host = hostParam[0]
		}
		if host == "" {
			http.Error(w, "Missing \"host\" parameter", 400)
			return
		}

		// get node details from apiserver
		node, err := clientset.CoreV1().Nodes().Get(host, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			error := fmt.Sprintf("Unknown node: %s", host)
			http.Error(w, error, 400)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// verify if node is healthy
		healthy := checkNodeHealth(node)
		if (healthy) {
			fmt.Fprintf(w, "Node is healthy: %s\n", node.Name)
		} else {
			error := fmt.Sprintf("Node is SICK: %s", node.Name)
			log.Println(error)
			http.Error(w, error, *sickCode)
		}
	})

	log.Println("listen on", *addr)
	http.ListenAndServe(*addr, nil)
}
