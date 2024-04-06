package env

import (
    "fmt"
    "regexp"

    "github.com/3th1nk/cidr"
)

// ValidateCreateEnv ensures all the required information is specified before attempting to create a VPN
func (e Env) ValidateCreateEnv() []string {
    var errs []string

    if e.Cloudflare.Zone != "" && e.Cloudflare.ApiKey == "" {
        errs = append(errs, "CLOUDFLARE_APIKEY is mandatory when CLOUDFLARE_ZONE is set")
    }

    if e.Cloudflare.Zone != "" && len(e.Cloudflare.ApiKey) != 40 {
        errs = append(errs, fmt.Sprintf("CLOUDFLARE_APIKEY '%s' is not valid", e.Cloudflare.ApiKey))
    }

    if e.CloudServer.ApiKey == "" {
        errs = append(errs, "CLOUDSERVER_APIKEY is mandatory")
    }

    if e.Server.Name == "" {
        errs = append(errs, "SERVER_NAME is mandatory")
    } else if !regexp.MustCompile(`^\w[\w.-]*\w$`).MatchString(e.Server.Name) {
        errs = append(errs, fmt.Sprintf("SERVER_NAME '%s' is not a valid RFC 3696 subdomain", e.Server.Name))
    }

    interfaceCIDR, cidrErr := cidr.Parse(e.Wireguard.Interface.Address)
    if e.Wireguard.Interface.Address == "" {
        errs = append(errs, "WIREGUARD_ADDRESS is mandatory")
    } else if cidrErr != nil {
        errs = append(errs, fmt.Sprintf("WIREGUARD_ADDRESS '%s' is not a valid CIDR", e.Wireguard.Interface.Address))
    }

    if e.Wireguard.Interface.ListenPortAlpha != "" && e.Wireguard.Interface.ListenPortAlpha != "0" && (e.Wireguard.Interface.ListenPort < 1 || e.Wireguard.Interface.ListenPort > 65535) {
        errs = append(errs, "WIREGUARD_LISTENPORT must be numeric and between 0 and 65535 if specified")
    }

    if e.Wireguard.Interface.PrivateKey == "" {
        errs = append(errs, "WIREGUARD_PRIVATEKEY is mandatory")
    } else if len(e.Wireguard.Interface.PrivateKey) != 44 {
        errs = append(errs, fmt.Sprintf("WIREGUARD_PRIVATEKEY '%s' is not valid", e.Wireguard.Interface.PrivateKey))
    }

    for _, peer := range e.Wireguard.Peers {
        peerCIDR, err := cidr.Parse(peer.AllowedIPs)
        if err != nil {
            errs = append(errs, fmt.Sprintf("WIREGUARD_PEER%d_ALLOWEDIPS '%s' is not a valid CIDR", peer.ID, peer.AllowedIPs))
        } else if !interfaceCIDR.Contains(peerCIDR.IP().String()) {
            errs = append(errs, fmt.Sprintf("WIREGUARD_PEER%d_ALLOWEDIPS '%s' is not within WIREGUARD_ADDRESS '%s' CIDR", peer.ID, peer.AllowedIPs, e.Wireguard.Interface.Address))
        }

        if len(peer.PublicKey) != 44 {
            errs = append(errs, fmt.Sprintf("WIREGUARD_PEER%d_PUBLICKEY '%s' is not valid", peer.ID, peer.PublicKey))
        }
    }

    return errs
}

// ValidateDestroyEnv ensures all the required information is specified before attempting to remove a VPN
func (e Env) ValidateDestroyEnv() []string {
    var errs []string

    if e.CloudServer.ApiKey == "" {
        errs = append(errs, "CLOUDSERVER_APIKEY is mandatory")
    }

    return errs
}

// ValidateServeEnv ensures all the required information is specified for serving an HTTP server
func (e Env) ValidateServeEnv() []string {
    errs := e.ValidateCreateEnv()

    // Ensure port is numeric if specified
    if e.HTTP.PortAlpha != "" && e.HTTP.PortAlpha != "0" && (e.HTTP.Port < 1 || e.HTTP.Port > 65535) {
        errs = append(errs, "HTTP_PORT must be numeric and between 0 and 65535 if specified")
    }

    return errs
}
