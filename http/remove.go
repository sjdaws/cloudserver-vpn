package http

import (
    "net/http"

    "github.com/sjdaws/cloudserver-vpn/vps"
)

func (h *HTTP) remove(response http.ResponseWriter, _ *http.Request) {
    var active []int
    var err error

    servers, err := vps.ListActiveVPS(h.env)
    if err != nil {
        errorResponse(response, err)
        return
    }

    for _, server := range servers {
        active = append(active, server.ID)
    }

    for _, serverID := range active {
        err = vps.Destroy(h.env, serverID)
        if err != nil {
            errorResponse(response, err)
            return
        }
    }

    h.status(response, nil)
}
