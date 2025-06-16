package utils

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, mids ...Middleware) http.Handler {
	for _, mid := range mids {
		handler = mid(handler)
	}
	return handler
}

//rl := middlewares.NewRateLimiter(5, time.Minute)

//hppOptions := middlewares.HPPOptions{
//	CheckQuery:                  true,
//	CheckBody:                   true,
//	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
//	Whitelist:                   []string{"name", "sortBy", "sortOrder", "age", "city"},
//}
//serveMux := applyMiddlewares(mux,
//	middlewares.Hpp(hppOptions),
//	middlewares.Compression,
//	middlewares.SecurityHeaders,
//	middlewares.ResponseTimeMiddleware,
//	rl.Middleware,
//	middlewares.Cors)
// create custom server
