package http

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/sjdaws/cloudserver-vpn/env"
)

type HTTP struct {
    env env.Env
}

const httpPort = 5252

// New creates a new HTTP instance
func New(env env.Env) *HTTP {
    return &HTTP{
        env: env,
    }
}

// Start configures and starts an HTTP server
func (h *HTTP) Start() error {
    errs := h.env.ValidateServeEnv()
    if len(errs) > 0 {
        return fmt.Errorf("unable to start http server:\n - %s", strings.Join(errs, "\n - "))
    }

    port := h.env.HTTP.Port
    if port == 0 && h.env.HTTP.PortAlpha != "0" {
        port = httpPort
    }

    http.HandleFunc("/create", h.create)
    http.HandleFunc("/remove", h.remove)
    http.HandleFunc("/status", h.status)
    log.Printf("listening on port %d", port)

    return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// errorResponse writes an error to the response buffer
func errorResponse(response http.ResponseWriter, err error) {
    response.WriteHeader(http.StatusInternalServerError)
    _, err = response.Write([]byte(err.Error()))
    if err != nil {
        log.Printf("unable to write http response: %v", err)
    }
}

// sendResponse marshals and sends an HTTP response
func sendResponse(response http.ResponseWriter, v any) {
    response.WriteHeader(http.StatusOK)
    body, err := json.Marshal(v)
    if err != nil {
        errorResponse(response, err)
        return
    }

    _, err = response.Write(body)
    if err != nil {
        errorResponse(response, err)
        return
    }
}
