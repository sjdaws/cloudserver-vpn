package dns

import (
    "context"
    "fmt"
    "strings"

    "github.com/cloudflare/cloudflare-go"
    "github.com/sjdaws/cloudserver-vpn/env"
    "github.com/sjdaws/cloudserver-vpn/vps"
)

// Configure a DNS record for the server
func Configure(env env.Env, vps *vps.VPS) error {
    api, err := cloudflare.NewWithAPIToken(env.Cloudflare.ApiKey)
    if err != nil {
        return fmt.Errorf("unable to connect to cloudflare api: %v", err)
    }

    ctx := context.Background()

    zones, err := api.ListZones(ctx, env.Cloudflare.Zone)
    if err != nil {
        return fmt.Errorf("unable to list dns zones: %v", err)
    }

    var zoneID string
    for _, zone := range zones {
        if strings.EqualFold(zone.Name, env.Cloudflare.Zone) {
            zoneID = zone.ID
        }
    }

    if zoneID == "" {
        return fmt.Errorf("unable to determine zone id for %s", env.Cloudflare.Zone)
    }

    rc := cloudflare.ZoneIdentifier(zoneID)

    records, _, err := api.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{Name: env.Server.FQDN})
    if err != nil {
        return fmt.Errorf("unable to list dns records for %s: %v", env.Cloudflare.Zone, err)
    }

    var recordID string
    for _, record := range records {
        if strings.EqualFold(record.Name, env.Server.FQDN) {
            recordID = record.ID
        }
    }

    // Create record if it doesn't exist, otherwise update
    if recordID == "" {
        proxied := false
        _, err = api.CreateDNSRecord(ctx, rc, cloudflare.CreateDNSRecordParams{Content: vps.IP, Name: env.Server.FQDN, Proxied: &proxied, Type: "A"})
    } else {
        _, err = api.UpdateDNSRecord(ctx, rc, cloudflare.UpdateDNSRecordParams{Content: vps.IP, ID: recordID})
    }

    if err != nil {
        return fmt.Errorf("unable to set dns record: %v", err)
    }

    return nil
}
