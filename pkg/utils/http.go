package utils

import (
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"net"
	"net/http"
	"time"
)

var (
	tracedTr = &nethttp.Transport{
		RoundTripper: tr,
	}
	tr = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       0,
	}
	tracedClient = &http.Client{
		Transport: tracedTr,
	}
)

func TracedRequest(r *http.Request) (*http.Response, error) {
	request, httpTracer := nethttp.TraceRequest(opentracing.GlobalTracer(), r)
	defer httpTracer.Finish()
	response, err := tracedClient.Do(request)
	if err != nil {
		return response, err
	}

	return response, err
}
