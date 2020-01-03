package main

import (
	"github.com/elitenomad/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

/*
	Create a newTest application which returns instance of
	an instance of application struct containing mocked
	dependencies
 */
func newTestApplication(t *testing.T) *application  {
	/*
		Initialize template Cache
	 */
	templateCache, err := newTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	/*
		Create session manager instance, with the same settings as production
	 */
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return  &application{
		infoLog: log.New(ioutil.Discard, "", 0),
		errorLog: log.New(ioutil.Discard, "", 0),
		session: session,
		snippets: &mock.SnippetModel{},
		users: &mock.UserModel{},
		templateCache: templateCache,
	}
}

/*
	Define a custom testServer type which anonymously embeds the http.TestServer
 */
type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}