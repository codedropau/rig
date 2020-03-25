package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/goji/httpauth"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Run the web server.
func Run(addr string, clientset kubernetes.Interface, refresh time.Duration, username, password string) error {
	pods := &PodList{}

	// Background task to refresh the list of Pods.
	go func(pods *PodList) {
		limiter := time.Tick(refresh)

		for {
			<-limiter

			list, err := getPodList(clientset)
			if err != nil {
				panic(err)
			}

			*pods = list
		}
	}(pods)

	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		pod, exist := getPod(r.Host, pods)
		if !exist {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment cannot be found. Try rebuilding to reinstate.")
			return
		}

		if pod.Status == corev1.PodPending {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment is currently building. Check back shortly.")
			return
		}

		if pod.Status == corev1.PodFailed {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment is currently in a failed state.")
			return
		}

		if pod.Status == corev1.PodUnknown {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment status is unknown. Please try rebuilding to reinstate.")
			return
		}

		if pod.Status == corev1.PodSucceeded {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "Environment has shutdown. Please rebuild to reinstate.")
			return
		}

		// @todo, Determine if we remove the hardcoded "http://".
		endpoint := fmt.Sprintf("http://%s", pod.IP)

		url, err := url.Parse(endpoint)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(url)

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)

		proxy.ServeHTTP(w, r)
	})

	http.Handle("/", httpauth.SimpleBasicAuth(username, password)(handler))

	return http.ListenAndServe(addr, nil)
}
