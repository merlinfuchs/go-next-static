package gonextstatic

import (
	"io/fs"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	Path                string
	ExpectedCode        int
	ExpectedBody        string
	ExpectedContentType string
}{
	{
		Path:                "/",
		ExpectedCode:        200,
		ExpectedBody:        "index.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/test.txt",
		ExpectedCode:        200,
		ExpectedBody:        "text.txt",
		ExpectedContentType: "text/plain; charset=utf-8",
	},
	{
		Path:                "/not/found",
		ExpectedCode:        200,
		ExpectedBody:        "index.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/abc",
		ExpectedCode:        200,
		ExpectedBody:        "[arg3].html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/noarg",
		ExpectedCode:        200,
		ExpectedBody:        "noarg.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/abc/notfound",
		ExpectedCode:        200,
		ExpectedBody:        "index.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/abc/page",
		ExpectedCode:        200,
		ExpectedBody:        "[arg1]/page.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
	{
		Path:                "/abc/abc/page",
		ExpectedCode:        200,
		ExpectedBody:        "[arg1]/[arg2]/page.html",
		ExpectedContentType: "text/html; charset=utf-8",
	},
}

func TestParse(t *testing.T) {

	t.Parallel()

	testData := os.DirFS("test_data")

	handler, err := NewHandler(testData.(fs.StatFS))
	require.NoError(t, err)
	require.NotNil(t, handler)

	for _, testCase := range testCases {
		t.Run("path:"+testCase.Path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", testCase.Path, nil)

			handler.ServeHTTP(res, req)

			body := res.Body.String()
			body = strings.TrimSpace(body)

			assert.Equal(t, testCase.ExpectedCode, res.Code)
			assert.Equal(t, testCase.ExpectedBody, body)
			assert.Equal(t, testCase.ExpectedContentType, res.Header().Get("Content-Type"))
		})
	}
}
