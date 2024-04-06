package env

import (
    "fmt"
    "os"
    "strings"
    "syscall"

    "github.com/sjdaws/cloudserver-vpn/helpers"
)

type Cloudflare struct {
    ApiKey string
    Zone   string
}

type CloudServer struct {
    ApiKey  string
    Project int
}

type Env struct {
    Cloudflare  Cloudflare
    CloudServer CloudServer
    HTTP        HTTP
    Server      Server
    Wireguard   Wireguard
}

type HTTP struct {
    Port      int
    PortAlpha string
}

type Interface struct {
    Address         string
    ListenPort      int
    ListenPortAlpha string
    PrivateKey      string
}

type Peer struct {
    AllowedIPs string
    ID         int
    PublicKey  string
}

type Server struct {
    FQDN string
    Name string
}

type Wireguard struct {
    Interface Interface
    Peers     []Peer
}

// Read environment variables into struct
func Read() Env {
    var env Env

    // Cloudflare
    env.Cloudflare.ApiKey = os.Getenv("CLOUDFLARE_APIKEY")
    env.Cloudflare.Zone = os.Getenv("CLOUDFLARE_ZONE")

    // Voyager
    env.CloudServer.ApiKey = os.Getenv("CLOUDSERVER_APIKEY")
    env.CloudServer.Project = helpers.AtoI(os.Getenv("CLOUDSERVER_PROJECT"))

    // HTTP server
    env.HTTP.Port = helpers.AtoI(os.Getenv("HTTP_PORT"))
    env.HTTP.PortAlpha = os.Getenv("HTTP_PORT")

    // Server
    env.Server.Name = os.Getenv("SERVER_NAME")
    env.Server.FQDN = env.Server.Name

    // Wireguard interface
    env.Wireguard.Interface.Address = os.Getenv("WIREGUARD_ADDRESS")
    env.Wireguard.Interface.ListenPort = helpers.AtoI(os.Getenv("WIREGUARD_LISTENPORT"))
    env.Wireguard.Interface.ListenPortAlpha = os.Getenv("WIREGUARD_LISTENPORT")
    env.Wireguard.Interface.PrivateKey = os.Getenv("WIREGUARD_PRIVATEKEY")

    // Wireguard peers
    for i := 0; i <= 254; i++ {
        allowedIPs, ipsFound := syscall.Getenv(fmt.Sprintf("WIREGUARD_PEER%d_ALLOWEDIPS", i))
        publicKey, pkFound := syscall.Getenv(fmt.Sprintf("WIREGUARD_PEER%d_PUBLICKEY", i))

        if ipsFound && pkFound {
            env.Wireguard.Peers = append(env.Wireguard.Peers, Peer{AllowedIPs: allowedIPs, ID: i, PublicKey: publicKey})
        }
    }

    // Calculated settings
    if env.Cloudflare.Zone != "" {
        env.Server.FQDN = strings.ToLower(fmt.Sprintf("%s.%s", env.Server.Name, env.Cloudflare.Zone))
    }

    return env
}
