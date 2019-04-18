package function

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	t.Skip()

	values := url.Values{}
	values.Add("url", "https://i.ytimg.com/vi/mSnMTMlJzME/maxresdefault.jpg")
	values.Add("name", "Emma Watson")
	os.Setenv("machinebox_url", "https://gateway.ofc.matiaspan.me/function/machine-box")

	r := &http.Request{
		Form: values,
	}
	res := &mockResponse{
		buf: bytes.NewBuffer([]byte{}),
	}

	Handle(res, r)

	b, err := ioutil.ReadAll(res.buf)
	if err != nil {
		t.Fatalf("Could not read buffer: %s", err)
	}
	fmt.Println(string(b))
}

type mockResponse struct {
	buf io.ReadWriter
}

func (m *mockResponse) WriteHeader(code int) {}

func (m *mockResponse) Write(b []byte) (n int, err error) {
	return m.buf.Write(b)
}

func (m *mockResponse) Header() http.Header {
	return http.Header{}
}
