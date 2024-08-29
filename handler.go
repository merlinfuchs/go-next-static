package gonextstatic

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
)

var dynamicRouteParamRE = regexp.MustCompile(`\\\[(\w+)\\\]`)
var dynamicRouteCatchAllRE = regexp.MustCompile(`\\\[\.\.\.(\w+)\\\]`)

type Route struct {
	RoutePattern *regexp.Regexp
	FilePath     string
}

func NewHandler(f fs.StatFS) (*Handler, error) {
	routes := make([]Route, 0)

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			route := strings.TrimSuffix(path, ".html")
			route = regexp.QuoteMeta(route)

			route = dynamicRouteParamRE.ReplaceAllString(route, "[a-zA-Z0-9_.=-]+")
			route = dynamicRouteCatchAllRE.ReplaceAllString(route, "[a-zA-Z0-9_.=/-]+")

			route = "^" + route + "$"

			routes = append(routes, Route{
				RoutePattern: regexp.MustCompile(route),
				FilePath:     path,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// We want to sort the routes so that the most specific routes come first.
	// This is not perfect and may not work for complex routes.
	slices.SortFunc(routes, func(i, j Route) int {
		a := strings.Index(i.FilePath, "[")
		b := strings.Index(j.FilePath, "[")

		if a == -1 && b == -1 {
			return 0
		}

		if a == -1 {
			return -1
		}

		if b == -1 {
			return 1
		}

		return b - a
	})

	for i := range routes {
		r := routes[i]
		fmt.Println(r.RoutePattern, r.FilePath)
	}

	return &Handler{
		f:      f,
		routes: routes,
	}, nil
}

type Handler struct {
	f      fs.StatFS
	routes []Route
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	if path == "" {
		http.ServeFileFS(res, req, h.f, "index.html")
		return
	}

	_, err := h.f.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		http.ServeFileFS(res, req, h.f, path)
		return
	}

	for _, r := range h.routes {
		if r.RoutePattern.MatchString(path) {
			http.ServeFileFS(res, req, h.f, r.FilePath)
			return
		}
	}

	http.ServeFileFS(res, req, h.f, "index.html")
}
