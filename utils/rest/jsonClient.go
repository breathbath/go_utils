package rest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	io2 "github.com/breathbath/go_utils/utils/io"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type RequestContext struct {
	TargetUrl    string
	Method       string
	Body         string
	Headers      map[string]string
	ProxyUrl     string
	LoggingTopic string
	IsVerbose    bool
}

func (rc RequestContext) String() string {
	proxyUrl := ""
	if rc.ProxyUrl != "" {
		proxyUrl = ", proxy: " + rc.ProxyUrl
	}
	return fmt.Sprintf(
		"Request: method %s, url '%s', body '%s', headers: %v%s",
		rc.Method,
		rc.TargetUrl,
		rc.Body,
		rc.Headers,
		proxyUrl,
	)
}

type JsonClient struct{}

func NewJsonClient() JsonClient {
	return JsonClient{}
}

func (jc JsonClient) CallApi(requestContext RequestContext) ([]byte, error, *http.Response) {
	req, err := http.NewRequest(
		requestContext.Method,
		requestContext.TargetUrl,
		strings.NewReader(requestContext.Body),
	)
	if err != nil {
		return []byte{}, err, nil
	}

	for key, value := range requestContext.Headers {
		req.Header.Add(key, value)
	}

	connectionTimeout := 30 * time.Second
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: connectionTimeout,
	}

	if requestContext.ProxyUrl != "" {
		proxyUrl, err := url.Parse(requestContext.ProxyUrl)
		if err != nil {
			return []byte{}, err, nil
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	client := http.Client{Transport: transport}

	if requestContext.IsVerbose {
		dump, _ := httputil.DumpRequest(req, true)
		io2.OutputInfo(requestContext.LoggingTopic, "Input context: %s, raw request: %s", requestContext.String(), string(dump))
	} else {
		//hiding details here for security reasons (no sensitive data in logs)
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
		return []byte{}, fmt.Errorf("Request failed with error: %v", err), resp
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			io2.OutputError(closeErr, requestContext.LoggingTopic, "")
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		return []byte{}, fmt.Errorf("Empty body in the response, status: %d", resp.StatusCode), resp
	}
	if err != nil {
		return []byte{}, fmt.Errorf("Reading of the request body failed with error: %v, status: %d", err, resp.StatusCode), resp
	}

	io2.OutputInfo(requestContext.LoggingTopic, "Got response: '%s', status code: '%d'", string(respBody), resp.StatusCode)

	err = ValidateResponse(requestContext.TargetUrl, resp, respBody)
	return respBody, err, resp
}

func (jc JsonClient) Get(context RequestContext) ([]byte, error, *http.Response) {
	context.Method = http.MethodGet
	return jc.CallApi(context)
}

func (jc JsonClient) Post(context RequestContext) ([]byte, error, *http.Response) {
	context.Method = http.MethodPost
	return jc.CallApi(context)
}

func (jc JsonClient) ScanToTarget(context RequestContext, target interface{}) error {
	body, err, _ := jc.Get(context)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		err := fmt.Errorf("Cannot process response %s: %v", string(body), err)
		return err
	}

	return nil
}

func (jc JsonClient) ScanToTargetRecoveringOnProxyFailure(context RequestContext, target interface{}) error {
	err := jc.ScanToTarget(context, target)
	if err != nil && context.ProxyUrl != "" {
		io2.OutputWarning("", "Request failure: %v. Will try to repeat without proxy", err)
		context.ProxyUrl = ""
		err = jc.ScanToTarget(context, target)
	}

	return err
}
