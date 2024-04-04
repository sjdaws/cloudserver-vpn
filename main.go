package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/sjdaws/cloudserver-vpn/dns"
    "github.com/sjdaws/cloudserver-vpn/env"
    "github.com/sjdaws/cloudserver-vpn/helpers"
    "github.com/sjdaws/cloudserver-vpn/vps"
)

const usageText = `
usage: %s [OPTIONS]

Run a New Zealand based VPN on demand for only 1.5c per hour

Options:

  --create       Create a VPN server
  --remove       Remove all created VPN servers
  --remove id    Remove a single VPN server

`

func main() {
    if len(os.Args) == 1 {
        fmt.Printf(usageText, os.Args[0])
        os.Exit(0)
    }

    config := env.Read()

    switch strings.ToLower(os.Args[1]) {
    case "--create":
        log.Print("Creating and configuring vps")

        server, err := vps.Create(config)
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("VPS created, ID: %d, IP: %s", server.ID, server.IP)

        if config.Cloudflare.Zone != "" {
            log.Printf("Configuring DNS record %s", config.Server.FQDN)

            err = dns.Configure(config, server)
            if err != nil {
                log.Fatal(err)
            }

            log.Print("DNS configured")
        }

        log.Print("Completed successfully")

    case "--remove":
        var active []int
        var err error

        if len(os.Args) == 3 {
            active = []int{helpers.AtoI(os.Args[2])}
        } else {
            active, err = vps.ListActiveVPS(config)
            if err != nil {
                log.Fatal(err)
            }

            log.Printf("Found %d server(s) to clean up", len(active))
        }

        for _, serverID := range active {
            log.Printf("Removing server %d", serverID)

            err = vps.Destroy(config, serverID)
            if err != nil {
                log.Fatal(err)
            }

            log.Printf("Server %d removed", serverID)
        }

        log.Print("Completed successfully")
    }
}
