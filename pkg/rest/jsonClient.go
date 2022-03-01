package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	io2 "github.com/breathbath/go_utils/v3/pkg/io"
)

type RequestContext struct {
	TargetURL    string
	Method       string
	Body         string
	Headers      map[string]string
	ProxyURL     string
	LoggingTopic string
	IsVerbose    bool
}

func (rc *RequestContext) String() string {
	proxyURL := ""
	if rc.ProxyURL != "" {
		proxyURL = ", proxy: " + rc.ProxyURL
	}
	return fmt.Sprintf(
		"Request: method %s, url '%s', body '%s', headers: %v%s",
		rc.Method,
		rc.TargetURL,
		rc.Body,
		rc.Headers,
		proxyURL,
	)
}

type JSONClient struct{}

func NewJSONClient() *JSONClient {
	return &JSONClient{}
}

func (jc *JSONClient) CallAPI(ctx context.Context, requestContext *RequestContext) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		requestContext.Method,
		requestContext.TargetURL,
		strings.NewReader(requestContext.Body),
	)
	if err != nil {
		return []byte{}, nil, err
	}

	for key, value := range requestContext.Headers {
		req.Header.Add(key, value)
	}

	connectionTimeout := 30 * time.Second
	transport := &http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: connectionTimeout,
	}

	if requestContext.ProxyURL != "" {
		proxyURL, e := url.Parse(requestContext.ProxyURL)
		if e != nil {
			return []byte{}, nil, e
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := http.Client{Transport: transport}

	if requestContext.IsVerbose {
		dump, _ := httputil.DumpRequest(req, true)
		io2.OutputInfo(requestContext.LoggingTopic, "Input context: %s, raw request: %s", requestContext.String(), string(dump))
	} else {
		// hiding details here for security reasons (no sensitive data in logs)
		io2.OutputInfo(requestContext.LoggingTopic, "Calling api")
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			closeErr := resp.Body.Close()
			if closeErr != nil {
				io2.OutputError(closeErr, requestContext.LoggingTopic, "")
			}
		}
		return []byte{}, resp, fmt.Errorf("request failed with error: %v", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			io2.OutputError(closeErr, requestContext.LoggingTopic, "")
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err == io.EOF {
		return []byte{}, resp, fmt.Errorf("empty body in the response, status: %d", resp.StatusCode)
	}
	if err != nil {
		return []byte{}, resp, fmt.Errorf("reading of the request body failed with error: %v, status: %d", err, resp.StatusCode)
	}

	io2.OutputInfo(requestContext.LoggingTopic, "Got response: '%s', status code: '%d'", string(respBody), resp.StatusCode)

	err = ValidateResponse(requestContext.TargetURL, resp, respBody)
	return respBody, resp, err
}

func (jc *JSONClient) Get(ctx context.Context, req *RequestContext) ([]byte, *http.Response, error) {
	req.Method = http.MethodGet
	return jc.CallAPI(ctx, req)
}

func (jc *JSONClient) Post(ctx context.Context, req *RequestContext) ([]byte, *http.Response, error) {
	req.Method = http.MethodPost
	return jc.CallAPI(ctx, req)
}

func (jc *JSONClient) ScanToTarget(ctx context.Context, req *RequestContext, target interface{}) error {
	body, resp, err := jc.Get(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, target); err != nil {
		e := fmt.Errorf("cannot process response %s: %v", string(body), err)
		return e
	}

	return nil
}

func (jc *JSONClient) ScanToTargetRecoveringOnProxyFailure(ctx context.Context, req *RequestContext, target interface{}) error {
	err := jc.ScanToTarget(ctx, req, target)
	if err != nil && req.ProxyURL != "" {
		io2.OutputWarning("", "Request failure: %v. Will try to repeat without proxy", err)
		req.ProxyURL = ""
		err = jc.ScanToTarget(ctx, req, target)
	}

	return err
}
