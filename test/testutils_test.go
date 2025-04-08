package test

import (
	"bytes"
	"database/sql"
	_ "embed"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"snippetbox/internal/modules"
	"snippetbox/internal/platform/database"
	"snippetbox/internal/platform/web"
	"snippetbox/internal/ui"
	"testing"
	"time"
)

//go:embed testdata/setup.sql
var setupDB string

//go:embed testdata/teardown.sql
var teardownDB string

func newTestDB(t *testing.T) (database.DB, error) {
	db, err := sql.Open("mysql", "test_web:pass@/test?parseTime=true&multiStatements=true")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(setupDB)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		_, err := db.Exec(teardownDB)
		if err != nil {
			t.Fatal(err)
		}

		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	return db, nil
}

var csrfTokenRgx = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func ExtractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRgx.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatalf("missing CSRF token in response body")
	}

	return html.UnescapeString(matches[1])
}

func NewApplication(t *testing.T) *web.Application {
	tc, err := ui.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	db, err := newTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	app := &web.Application{
		Router:         httprouter.New(),
		ErrorLog:       log.New(io.Discard, "", 0),
		InfoLog:        log.New(io.Discard, "", 0),
		TemplateCache:  tc,
		FormDecoder:    form.NewDecoder(),
		SessionManager: sessionManager,
		DB:             db,
	}
	modules.Init(app)

	return app
}

type TestServer struct {
	*httptest.Server
}

func (ts *TestServer) Get(t *testing.T, path string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *TestServer) PostForm(t *testing.T, path string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+path, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}
