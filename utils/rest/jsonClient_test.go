package rest

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var postData []string

func init() {
	postData = []string{}
}

func TestRequestContextToString(t *testing.T) {
	rc := RequestContext{
		TargetUrl: "ya.ru",
		Method: "GET",
		Body: "Lala",
		Headers: map[string]string{"head1": "headVal1"},
		ProxyUrl: "someProx.ru",
		LoggingTopic: "lala",
		IsVerbose: true,
	}

	expectedStr := "Request: method GET, url 'ya.ru', body 'Lala', headers: map[head1:headVal1], proxy: someProx.ru"
	assert.Equal(t, expectedStr, rc.String())
}

func TestJsonClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	httpServer := startHTTPServer(t, 18090)
	defer stopHTTPServer(httpServer)

	t.Run("testGet", testGet)
	t.Run("testPost", testPost)
	t.Run("testScan", testScan)
}

func testGet(t *testing.T) {
	cl := NewJsonClient()
	rc := RequestContext{
		TargetUrl: "http://localhost:18090",
		Method: "GET",
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
	postData = []string{}
	cl := NewJsonClient()
	rc := RequestContext{
		TargetUrl: "http://localhost:18090",
		Method: "POST",
		Body: "Accept me please",
	}

	body, err, resp := cl.Post(rc)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"key":"val"}`, string(body))
	assert.Len(t, postData, 1)
	assert.Equal(t, "Accept me please", postData[0])
}

func testScan(t *testing.T) {
	item := struct{
		Key string `json:"key"`
	}{}

	rc := RequestContext{
		TargetUrl: "http://localhost:18090",
		Method: "GET",
	}

	cl := NewJsonClient()
	err := cl.ScanToTarget(rc, &item)
	assert.NoError(t, err)
	assert.Equal(t, "val", item.Key)
}

func startHTTPServer(t *testing.T, port int) *http.Server {
	var handlerFunc http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte(`{"key":"val"}`))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(body) > 0 {
			postData = append(postData, string(body))
		}
	}
	hostAddress := ":" + strconv.Itoa(port)
	httpServer := &http.Server{
		Addr:    hostAddress,
		Handler: handlerFunc,
	}

	errChan := make(chan error)
	go func(httpServer *http.Server) {
		err := httpServer.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}(httpServer)

	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(time.Millisecond * 100):
	}
	log.Printf("[INFO] HTTP Server started at [%v] \n", hostAddress)

	return httpServer
}

func stopHTTPServer(httpServer *http.Server) {
	err := httpServer.Close()
	if err != nil {
		log.Println(err)
	}
}