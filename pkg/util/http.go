package util

import (
	"bytes"
	"context"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/wg815737157/paper-work/pkg/log"
	"io"
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
func GetWithContext(ctx context.Context, url string) []byte {
	logger := log.SugarLogger()
	userInfoRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil
	}
	client := http.DefaultClient
	res, err := client.Do(userInfoRequest)
	if err != nil {
		logger.Error(err)
		return nil
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resBytes
}

func PostWithContext(ctx context.Context, url string, body []byte) []byte {
	logger := log.SugarLogger()
	userInfoRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		logger.Error(err)
		return nil
	}
	client := http.DefaultClient
	res, err := client.Do(userInfoRequest)
	if err != nil {
		logger.Error(err)
		return nil
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resBytes
}
