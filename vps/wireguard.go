package vps

import (
    "encoding/base64"
    "fmt"

    "github.com/sjdaws/cloudserver-vpn/env"
)

const interfaceTemplate = `[Interface]
Address = %s
ListenPort = %d
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PrivateKey = %s

`
const listenPort = 51820
const peerTemplate = `[Peer]
AllowedIPs = %s
PublicKey = %s

`

func generateEncodedWireguardConfiguration(env env.Env) string {
    port := env.Wireguard.Interface.ListenPort
    if port == 0 && env.Wireguard.Interface.ListenPortAlpha != "0" {
        port = listenPort
    }

    config := fmt.Sprintf(interfaceTemplate, env.Wireguard.Interface.Address, port, env.Wireguard.Interface.PrivateKey)

    for _, peer := range env.Wireguard.Peers {
        config += fmt.Sprintf(peerTemplate, peer.AllowedIPs, peer.PublicKey)
    }

    return base64.StdEncoding.EncodeToString([]byte(config))
}
