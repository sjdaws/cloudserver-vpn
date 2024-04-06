package http

import (
    "net/http"

    "github.com/sjdaws/cloudserver-vpn/dns"
    "github.com/sjdaws/cloudserver-vpn/env"
    "github.com/sjdaws/cloudserver-vpn/vps"
)

type Status struct {
    DNS  bool   `json:"dns,omitempty"`
    ID   int    `json:"id"`
    IP   string `json:"ip"`
    Name string `json:"name"`
}

// status returns VPS and optionally DNS status for active VPS
func (h *HTTP) status(response http.ResponseWriter, _ *http.Request) {
    statuses, err := getVPSStatus(h.env)
    if err != nil {
        errorResponse(response, err)
        return
    }

    if h.env.Cloudflare.Zone != "" {
        statuses = getDNSStatus(h.env, statuses)
    }

    sendResponse(response, statuses)
}

// getDNSStatus resolves dns for specified VPS
func getDNSStatus(env env.Env, statuses []Status) []Status {
    for id, status := range statuses {
        content, _ := dns.Retrieve(env, status.Name)
        statuses[id].DNS = false
        if content == status.IP {
            statuses[id].DNS = true
        }
    }

    return statuses
}

// getVPSStatus gets the status of active VPS
func getVPSStatus(env env.Env) ([]Status, error) {
    servers, err := vps.ListActiveVPS(env)
    if err != nil {
        return nil, err
    }

    statuses := make([]Status, 0)
    for _, server := range servers {
        var ip string
        for _, ipData := range server.IPs {
            if ipData.Primary == true {
                ip = ipData.IP
                break
            }
        }

        statuses = append(statuses, Status{
            ID:   server.ID,
            IP:   ip,
            Name: server.Name,
        })
    }

    return statuses, nil
}
