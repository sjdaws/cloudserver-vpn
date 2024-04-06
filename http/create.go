package http

import (
    "net/http"

    "github.com/sjdaws/cloudserver-vpn/dns"
    "github.com/sjdaws/cloudserver-vpn/vps"
)

func (h *HTTP) create(response http.ResponseWriter, _ *http.Request) {
    server, err := vps.Create(h.env)
    if err != nil {
        errorResponse(response, err)
        return
    }

    if h.env.Cloudflare.Zone != "" {
        err = dns.Configure(h.env, server)
        if err != nil {
            errorResponse(response, err)
            return
        }
    }

    h.status(response, nil)
}
