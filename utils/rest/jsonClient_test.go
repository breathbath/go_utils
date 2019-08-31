package rest

import (
	http2 "github.com/breathbath/go_utils/utils/http"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type RequestMock struct {
	Method     string
	Url        *url.URL
	Header     http.Header
	Body       string
	RequestUri string
}

func NewRequestMock(r *http.Request) RequestMock {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.OutputError(err, "", "")
	}

	return RequestMock{
		Method: r.Method,
		Url:    r.URL,
		Header: r.Header,
		Body:   string(body),
		RequestUri: r.RequestURI,
	}
}

var requests []RequestMock
var serverAddr string
var proxyServerAddr string

func init() {
	requests = []RequestMock{}
}

func TestRequestContextToString(t *testing.T) {
	rc := RequestContext{
		TargetUrl:    "ya.ru",
		Method:       "GET",
		Body:         "Lala",
		Headers:      map[string]string{"head1": "headVal1"},
		ProxyUrl:     "someProx.ru",
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

	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "GET",
		Headers:   map[string]string{"head1": "headVal1"},
	}

	cl := NewJsonClient()
	_, err, _ := cl.Get(rc)
	assert.NoError(t, err)

	assert.Len(t, requests, 1)
	req := requests[0]
	headerVal := req.Header.Get("head1")
	assert.Equal(t, "headVal1", headerVal)
}

func testGet(t *testing.T) {
	requests = []RequestMock{}

	cl := NewJsonClient()
	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "GET",
	}

	body, err, resp := cl.Get(rc)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"key":"val"}`, string(body))
}

func testPost(t *testing.T) {
	requests = []RequestMock{}

	cl := NewJsonClient()
	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "POST",
		Body:      "Accept me please",
	}

	body, err, resp := cl.Post(rc)
	assert.NoError(t, err)
	if err != nil {
		return
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

	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "GET",
	}

	cl := NewJsonClient()
	err := cl.ScanToTarget(rc, &item)
	assert.NoError(t, err)
	assert.Equal(t, "val", item.Key)

	wrongItem := struct {
		Key int `json:"key"`
	}{}
	err = cl.ScanToTarget(rc, &wrongItem)
	assert.EqualError(t, err, "Cannot process response {\"key\":\"val\"}: json: cannot unmarshal string into Go struct field .key of type int")

	requestWithExpectedErrResp := RequestContext{
		TargetUrl: serverAddr + "?err=400",
		Method:    "GET",
	}
	err = cl.ScanToTarget(requestWithExpectedErrResp, &item)
	assert.IsType(t, BadResponseCodeError{}, err)

	badRespErr := err.(BadResponseCodeError)
	assert.Equal(t, 400, badRespErr.resp.StatusCode)
}

func testInvalidMethod(t *testing.T) {
	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "мама",
	}
	cl := NewJsonClient()
	_, err, _ := cl.CallApi(rc)
	assert.EqualError(t, err, `net/http: invalid method "мама"`)
}

func testServerErrors(t *testing.T) {
	rc := RequestContext{
		TargetUrl: serverAddr + "?err=500&body=lals",
		Method:    "GET",
	}
	cl := NewJsonClient()
	_, err, _ := cl.CallApi(rc)
	assert.IsType(t, BadResponseCodeError{}, err)

	badRespErr := err.(BadResponseCodeError)
	assert.Equal(t, 500, badRespErr.resp.StatusCode)
	assert.Equal(t, "lals", string(badRespErr.respBody))
}

func testInvalidAddress(t *testing.T) {
	rc := RequestContext{
		TargetUrl: "",
		Method:    "GET",
	}
	cl := NewJsonClient()
	_, err, _ := cl.Get(rc)
	assert.Contains(t, err.Error(), `Request failed with error`)
}

func testProxy(t *testing.T) {
	requests = []RequestMock{}

	rc := RequestContext{
		TargetUrl: serverAddr,
		Method:    "GET",
		ProxyUrl:  proxyServerAddr,
	}
	cl := NewJsonClient()
	_, err, _ := cl.Get(rc)
	assert.NoError(t, err)

	assert.Len(t, requests, 1)
	req := requests[0]

	assert.Equal(t, serverAddr, strings.Trim(req.RequestUri, "/"))
	assert.NotEqual(t, proxyServerAddr, strings.Trim(req.RequestUri, "/"))

	rc2 := RequestContext{
		TargetUrl: serverAddr,
		Method:    "GET",
		ProxyUrl:  ":::",
	}

	item := struct {
		Key string `json:"key"`
	}{}
	err = cl.ScanToTargetRecoveringOnProxyFailure(rc2, &item)
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
				io.OutputError(err, "", "")
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
