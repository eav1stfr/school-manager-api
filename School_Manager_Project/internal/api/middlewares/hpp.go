package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	Whitelist                   []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				// filter the body params
				filterBodyParams(r, options.Whitelist)
			}
			if options.CheckQuery && r.URL.Query() != nil {
				filterQueryParams(r, options.Whitelist)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()

	for key, value := range query {
		if len(value) > 1 {
			query.Set(key, value[0]) // first value
		}
		if !isWhiteListed(key, whitelist) {
			query.Del(key)
		}
	}
	r.URL.RawQuery = query.Encode()
}

func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	for key, value := range r.Form {
		if len(value) > 1 {
			r.Form.Set(key, value[0]) // first value
		}
		if !isWhiteListed(key, whitelist) {
			delete(r.Form, key)
		}
	}
}

func isWhiteListed(param string, whitelist []string) bool {
	for _, value := range whitelist {
		if param == value {
			return true
		}
	}
	return false
}
