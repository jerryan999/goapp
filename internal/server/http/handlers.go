package http

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/webgo/v6"

	"github.com/jerryan999/goapp/internal/api"
)

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	api  *api.API
	home *template.Template
}

func (h *Handlers) routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:          "helloworld",
			Pattern:       "",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.HelloWorld)},
			TrailingSlash: true,
		},
		{
			Name:          "health",
			Pattern:       "/-/health",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.Health)},
			TrailingSlash: true,
		},
		{
			Name:          "create-user",
			Pattern:       "/users",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{errWrapper(h.CreateUser)},
			TrailingSlash: true,
		},
		{
			Name:          "read-user-byemail",
			Pattern:       "/users/:email",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.ReadUserByEmail)},
			TrailingSlash: true,
		},
	}
}

// Health is the HTTP handler to return the status of the app including the version, and other details
// This handler uses webgo to respond to the http request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) error {
	out, err := h.api.Health()
	if err != nil {
		return err
	}
	webgo.R200(w, out)
	return nil
}

// HelloWorld is a helloworld HTTP handler
func (h *Handlers) HelloWorld(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		{
			webgo.SendResponse(w, "hello world", http.StatusOK)
		}
	default:
		{
			buff := bytes.NewBufferString("")
			err := h.home.Execute(
				buff,
				struct {
					Message string
				}{
					Message: "welcome to the home page!",
				},
			)
			if err != nil {
				return errors.InternalErr(err, "Inter server error")
			}

			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			_, err = w.Write(buff.Bytes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func errWrapper(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		webgo.SendError(w, msg, status)
	}
}

func panicRecoverer(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		fmt.Println(string(debug.Stack()))
		webgo.R500(w, errors.DefaultMessage)
	}()

	next(w, r)
}

func loadHomeTemplate(basePath string) (*template.Template, error) {
	t := template.New("index.html")
	home, err := t.ParseFiles(
		fmt.Sprintf("%s/index.html", basePath),
	)
	if err != nil {
		return nil, err
	}

	return home, nil
}
