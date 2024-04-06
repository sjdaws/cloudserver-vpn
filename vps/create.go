package vps

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "strings"

    "github.com/sjdaws/cloudserver-vpn/env"
)

type IP struct {
    IP      string `json:"ip"`
    Primary bool   `json:"is_primary"`
}

type Server struct {
    Data     ServerData `json:"data,omitempty"`
    FQDNs    []string   `json:"fqdns"`
    IPTypes  []string   `json:"ip_types"`
    Location int        `json:"location"`
    Name     string     `json:"name"`
    OS       int        `json:"os"`
    Plan     int        `json:"plan"`
    Project  int        `json:"project"`
    UserData string     `json:"user_data"`
}

type ServerData struct {
    ID   int    `json:"id"`
    IPs  []IP   `json:"ips"`
    Name string `json:"name"`
}

type VPS struct {
    ID   int
    IP   string
    Name string
}

const userdataTemplate = `#cloud-config
write_files:
- content: bmV0LmlwdjQuY29uZi5hbGwucHJveHlfYXJwPTEKbmV0LmlwdjQuaXBfZm9yd2FyZD0xCg==
  encoding: b64
  owner: root:root
  path: /etc/sysctl.d/wireguard.conf
  permissions: '0644'
- content: %s
  encoding: b64
  owner: root:root
  path: /etc/wireguard/wg0.conf
  permissions: '0600'
- content: IyEvc2Jpbi9vcGVucmMtcnVuCgpkZXBlbmQoKSB7CiAgICBuZWVkIGxvY2FsbW91bnQgbmV0CiAgICB1c2UgZG5zCiAgICBhZnRlciBib290bWlzYwp9CgpjaGVja2NvbmZpZygpIHsKICAgICMgVE9ETzogZG9lcyB3aXJlZ3VhcmQgbW9kdWxlIGlzIGxvYWRlZAogICAgcmV0dXJuIDAKfQoKc3RhcnQoKSB7CiAgICBlYmVnaW4gIlN0YXJ0aW5nIFdpcmVndWFyZCIKCiAgICBjaGVja2NvbmZpZyB8fCByZXR1cm4gMQogICAgCiAgICB3Zy1xdWljayB1cCB3ZzAKICAgIGVlbmQgJD8KfQoKc3RvcCgpIHsKICAgIGViZWdpbiAiU3RvcHBpbmcgV2lyZWd1YXJkIgogICAgd2ctcXVpY2sgZG93biB3ZzAKICAgIGVlbmQgJD8KfQo=
  encoding: b64
  owner: root:root
  path: /etc/init.d/wireguard
  permissions: '0755'
runcmd:
  - sysctl -p /etc/sysctl.d/wireguard.conf
  - apk add wireguard-tools
  - rc-update add wireguard default
  - rc-service wireguard start`

// Create a new virtual private server
func Create(env env.Env) (*VPS, error) {
    errs := env.ValidateCreateEnv()
    if len(errs) > 0 {
        return nil, fmt.Errorf("unable to create new server:\n - %s", strings.Join(errs, "\n - "))
    }

    projectID, err := findOrCreateProject(env)
    if err != nil {
        return nil, err
    }

    payload, err := json.Marshal(&Server{
        FQDNs:    []string{env.Server.FQDN},
        IPTypes:  []string{"IPv4"},
        Location: 1,
        Name:     env.Server.FQDN,
        OS:       15,
        Plan:     29,
        Project:  projectID,
        UserData: fmt.Sprintf(userdataTemplate, generateEncodedWireguardConfiguration(env)),
    })
    if err != nil {
        return nil, fmt.Errorf("unable to marshal new server configuration: %v", err)
    }

    request, err := http.NewRequest("POST", fmt.Sprintf("%s/servers", apiURL), bytes.NewBuffer(payload))
    if err != nil {
        return nil, fmt.Errorf("unable to create new server request: %v", err)
    }

    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.CloudServer.ApiKey))
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("unable to create new server: %v", err)
    }
    defer closeBody(response.Body)

    body, _ := io.ReadAll(response.Body)

    if response.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("unable to create new server, invalid status: %s - %s", response.Status, string(body))
    }

    var result Server
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("error reading response from server: %v", err)
    }

    var primaryIP string
    for _, ip := range result.Data.IPs {
        if ip.Primary {
            primaryIP = ip.IP
            break
        }
    }

    if result.Data.ID == 0 || primaryIP == "" {
        return nil, errors.New("unable to detect whether server was successfully created: perform a manual check")
    }

    return &VPS{
        ID:   result.Data.ID,
        IP:   primaryIP,
        Name: result.Data.Name,
    }, nil
}
