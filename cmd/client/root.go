package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

var (
	verbose bool
	host    string
)

var rootCmdOpts struct {
	Verbose          bool
	Host             string
	Kubeconfig       string
	K8sNamespace     string
	K8sLabelSelector string
	K8sPodPort       string
	DialMode         string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbm",
	Short: "Bhojpur DBM is an elastic database management system powered by Kubernetes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("verbose logging enabled")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type dialMode string

const (
	dialModeHost       = "host"
	dialModeKubernetes = "kubernetes"
)

func init() {
	dbmHost := os.Getenv("DBM_HOST")
	if dbmHost == "" {
		dbmHost = "localhost:7777"
	}
	dbmKubeconfig := os.Getenv("KUBECONFIG")
	if dbmKubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.WithError(err).Warn("cannot determine user's home directory")
		} else {
			dbmKubeconfig = filepath.Join(home, ".kube", "config")
		}
	}
	dbmNamespace := os.Getenv("DBM_K8S_NAMESPACE")
	dbmLabelSelector := os.Getenv("DBM_K8S_LABEL")
	if dbmLabelSelector == "" {
		dbmLabelSelector = "app.kubernetes.io/name=dbm"
	}
	dbmPodPort := os.Getenv("DBM_K8S_POD_PORT")
	if dbmPodPort == "" {
		dbmPodPort = "7777"
	}
	dialMode := os.Getenv("DBM_DIAL_MODE")
	if dialMode == "" {
		dialMode = string(dialModeHost)
	}

	rootCmd.PersistentFlags().BoolVar(&rootCmdOpts.Verbose, "verbose", false, "en/disable verbose logging")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.DialMode, "dial-mode", dialMode, "dial mode that determines how we connect to Bhojpur DBM. Valid values are \"host\" or \"kubernetes\" (defaults to DBM_DIAL_MODE env var).")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Host, "host", dbmHost, "[host dial mode] Bhojpur DBM host to talk to (defaults to DBM_HOST env var)")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Kubeconfig, "kubeconfig", dbmKubeconfig, "[kubernetes dial mode] kubeconfig file to use (defaults to KUEBCONFIG env var)")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.K8sNamespace, "k8s-namespace", dbmNamespace, "[kubernetes dial mode] Kubernetes namespace in which to look for the Bhojpur DBM pods (defaults to DBM_K8S_NAMESPACE env var, or configured kube context namespace)")
	// The following are such specific flags that really only matters if one doesn't use the stock helm charts.
	// They can still be set using an env var, but there's no need to clutter the CLI with them.
	rootCmdOpts.K8sLabelSelector = dbmLabelSelector
	rootCmdOpts.K8sPodPort = dbmPodPort
}

type closableGrpcClientConnInterface interface {
	grpc.ClientConnInterface
	io.Closer
}

func dial() (res closableGrpcClientConnInterface) {
	var err error
	switch rootCmdOpts.DialMode {
	case dialModeHost:
		res, err = grpc.Dial(rootCmdOpts.Host, grpc.WithInsecure())
	case dialModeKubernetes:
		res, err = dialKubernetes()
	default:
		log.Fatalf("unknown dial mode: %s", rootCmdOpts.DialMode)
	}

	if err != nil {
		log.WithError(err).Fatal("cannot connect to Bhojpur DBM server")
	}
	return
}

func dialKubernetes() (closableGrpcClientConnInterface, error) {
	kubecfg, namespace, err := getKubeconfig(rootCmdOpts.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("cannot load kubeconfig %s: %w", rootCmdOpts.Kubeconfig, err)
	}
	if rootCmdOpts.K8sNamespace != "" {
		namespace = rootCmdOpts.K8sNamespace
	}

	clientSet, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, err
	}

	pod, err := findDbmPod(clientSet, namespace, rootCmdOpts.K8sLabelSelector)
	if err != nil {
		return nil, fmt.Errorf("cannot find Bhojpur DBM pod: %w", err)
	}

	localPort, err := findFreeLocalPort()

	ctx, cancel := context.WithCancel(context.Background())
	readychan, errchan := forwardPort(ctx, kubecfg, namespace, pod, fmt.Sprintf("%d:%s", localPort, rootCmdOpts.K8sPodPort))
	select {
	case err := <-errchan:
		cancel()
		return nil, err
	case <-readychan:
	}

	res, err := grpc.Dial(fmt.Sprintf("localhost:%d", localPort), grpc.WithInsecure())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("cannot dial forwarded connection: %w", err)
	}

	return closableConn{
		ClientConnInterface: res,
		Closer:              func() error { cancel(); return nil },
	}, nil
}

type closableConn struct {
	grpc.ClientConnInterface
	Closer func() error
}

func (c closableConn) Close() error {
	return c.Closer()
}

func findFreeLocalPort() (int, error) {
	const (
		start = 30000
		end   = 60000
	)
	for p := start; p <= end; p++ {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			l.Close()
			return p, nil
		}
	}
	return 0, fmt.Errorf("no free local port found")
}

// GetKubeconfig loads kubernetes connection config from a kubeconfig file
func getKubeconfig(kubeconfig string) (res *rest.Config, namespace string, err error) {
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{},
	)
	namespace, _, err = cfg.Namespace()
	if err != nil {
		return nil, "", err
	}

	res, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, namespace, err
	}

	return res, namespace, nil
}

// findDbmPod returns the first pod we found for a particular component
func findDbmPod(clientSet kubernetes.Interface, namespace, selector string) (podName string, err error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pod in %s with label component=%s", namespace, selector)
	}
	return pods.Items[0].Name, nil
}

// ForwardPort establishes a TCP port forwarding to a Kubernetes pod
func forwardPort(ctx context.Context, config *rest.Config, namespace, pod, port string) (readychan chan struct{}, errchan chan error) {
	errchan = make(chan error, 1)
	readychan = make(chan struct{}, 1)

	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		errchan <- err
		return
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, pod)
	hostIP := strings.TrimLeft(config.Host, "https://")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan := make(chan struct{}, 1)
	fwdReadyChan := make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	forwarder, err := portforward.New(dialer, []string{port}, stopChan, fwdReadyChan, out, errOut)
	if err != nil {
		panic(err)
	}

	var once sync.Once
	go func() {
		err := forwarder.ForwardPorts()
		if err != nil {
			errchan <- err
		}
		once.Do(func() { close(readychan) })
	}()

	go func() {
		select {
		case <-readychan:
			// we're out of here
		case <-ctx.Done():
			close(stopChan)
		}
	}()

	go func() {
		for range fwdReadyChan {
		}

		if errOut.Len() != 0 {
			errchan <- fmt.Errorf(errOut.String())
			return
		}

		once.Do(func() { close(readychan) })
	}()

	return
}
