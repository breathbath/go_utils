package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	io2 "github.com/breathbath/go_utils/v2/pkg/io"
	"github.com/stretchr/testify/require"

	http2 "github.com/breathbath/go_utils/v2/pkg/http"
	"github.com/stretchr/testify/assert"
)

type RequestMock struct {
	Method     string
	URL        *url.URL
	Header     http.Header
	Body       string
	RequestURI string
}

func NewRequestMock(r *http.Request) RequestMock {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		io2.OutputError(err, "", "")
	}

	return RequestMock{
		Method:     r.Method,
		URL:        r.URL,
		Header:     r.Header,
		Body:       string(body),
		RequestURI: r.RequestURI,
	}
}

var requests = []RequestMock{}
var serverAddr string
var proxyServerAddr string

func TestRequestContextToString(t *testing.T) {
	rc := &RequestContext{
		TargetURL:    "ya.ru",
		Method:       "GET",
		Body:         "Lala",
		Headers:      map[string]string{"head1": "headVal1"},
		ProxyURL:     "someProx.ru",
		LoggingTopic: "lala",
		IsVerbose:    true,
	}

	expectedStr := "Request: method GET, url 'ya.ru', body 'Lala', headers: map[head1:headVal1], proxy: someProx.ru"
	assert.Equal(t, expectedStr, rc.String())
}

func TestJsonClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	httpServer := startHTTPServer()
	defer httpServer.Close()

	proxyServer := startProxyServer()
	defer proxyServer.Close()

	t.Run("testGet", testGet)
	t.Run("testPost", testPost)
	t.Run("testScan", testScan)
	t.Run("testHeaders", testHeaders)
	t.Run("testInvalidMethod", testInvalidMethod)
	t.Run("testServerErrors", testServerErrors)
	t.Run("testInvalidAddress", testInvalidAddress)
	t.Run("testProxy", testProxy)
}

func testHeaders(t *testing.T) {
	requests = []RequestMock{}

	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "GET",
		Headers:   map[string]string{"head1": "headVal1"},
	}

	cl := NewJSONClient()
	_, resp, err := cl.Get(context.Background(), rc)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Len(t, requests, 1)
	req := requests[0]
	headerVal := req.Header.Get("head1")
	assert.Equal(t, "headVal1", headerVal)
}

func testGet(t *testing.T) {
	requests = []RequestMock{}

	cl := NewJSONClient()
	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "GET",
	}

	body, resp, err := cl.Get(context.Background(), rc)
	require.NoError(t, err)

	defer resp.Body.Close()

	if err != nil {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"key":"val"}`, string(body))
}

func testPost(t *testing.T) {
	requests = []RequestMock{}

	cl := NewJSONClient()
	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "POST",
		Body:      "Accept me please",
	}

	body, resp, err := cl.Post(context.Background(), rc)
	require.NoError(t, err)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"key":"val"}`, string(body))
	assert.Len(t, requests, 1)

	req := requests[0]

	assert.Equal(t, "Accept me please", req.Body)
}

func testScan(t *testing.T) {
	requests = []RequestMock{}

	item := struct {
		Key string `json:"key"`
	}{}

	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "GET",
	}

	cl := NewJSONClient()
	err := cl.ScanToTarget(context.Background(), rc, &item)
	assert.NoError(t, err)
	assert.Equal(t, "val", item.Key)

	wrongItem := struct {
		Key int `json:"key"`
	}{}
	err = cl.ScanToTarget(context.Background(), rc, &wrongItem)
	assert.EqualError(t, err, "cannot process response {\"key\":\"val\"}: json: cannot unmarshal string into Go struct field .key of type int")

	requestWithExpectedErrResp := &RequestContext{
		TargetURL: serverAddr + "?err=400",
		Method:    "GET",
	}
	err = cl.ScanToTarget(context.Background(), requestWithExpectedErrResp, &item)
	assert.IsType(t, BadResponseCodeError{}, err)

	badRespErr := err.(BadResponseCodeError)
	assert.Equal(t, 400, badRespErr.resp.StatusCode)
}

func testInvalidMethod(t *testing.T) {
	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "мама",
	}
	cl := NewJSONClient()
	_, resp, err := cl.CallAPI(context.Background(), rc)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	assert.EqualError(t, err, `net/http: invalid method "мама"`)
}

func testServerErrors(t *testing.T) {
	rc := &RequestContext{
		TargetURL: serverAddr + "?err=500&body=lals",
		Method:    "GET",
	}
	cl := NewJSONClient()
	_, resp, err := cl.CallAPI(context.Background(), rc)
	assert.IsType(t, BadResponseCodeError{}, err)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	badRespErr := err.(BadResponseCodeError)
	assert.Equal(t, 500, badRespErr.resp.StatusCode)
	assert.Equal(t, "lals", string(badRespErr.respBody))
}

func testInvalidAddress(t *testing.T) {
	rc := &RequestContext{
		TargetURL: "",
		Method:    "GET",
	}
	cl := NewJSONClient()
	_, resp, err := cl.Get(context.Background(), rc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `request failed with error`)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
}

func testProxy(t *testing.T) {
	requests = []RequestMock{}

	rc := &RequestContext{
		TargetURL: serverAddr,
		Method:    "GET",
		ProxyURL:  proxyServerAddr,
	}
	cl := NewJSONClient()
	_, resp, err := cl.Get(context.Background(), rc)
	assert.NoError(t, err)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	assert.Len(t, requests, 1)
	req := requests[0]

	assert.Equal(t, serverAddr, strings.Trim(req.RequestURI, "/"))
	assert.NotEqual(t, proxyServerAddr, strings.Trim(req.RequestURI, "/"))

	rc2 := &RequestContext{
		TargetURL: serverAddr,
		Method:    "GET",
		ProxyURL:  ":::",
	}

	item := struct {
		Key string `json:"key"`
	}{}
	err = cl.ScanToTargetRecoveringOnProxyFailure(context.Background(), rc2, &item)
	assert.NoError(t, err)
}

func startHTTPServer() *httptest.Server {
	var handlerFunc http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		rm := NewRequestMock(r)
		requests = append(requests, rm)

		errCode := http2.GetRequestValueInt(r, "err", 0)
		if errCode > 0 {
			rw.WriteHeader(int(errCode))
			_, err := rw.Write([]byte(http2.GetRequestValueString(r, "body", "")))
			if err != nil {
				io2.OutputError(err, "", "")
			}
			return
		}

		_, err := rw.Write([]byte(`{"key":"val"}`))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	srv := httptest.NewServer(handlerFunc)
	serverAddr = srv.URL

	return srv
}

func startProxyServer() *httptest.Server {
	var handlerFunc http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		rm := NewRequestMock(r)
		requests = append(requests, rm)
	}
	srv := httptest.NewServer(handlerFunc)
	proxyServerAddr = srv.URL

	return srv
}
