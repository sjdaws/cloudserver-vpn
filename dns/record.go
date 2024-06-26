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

    rc, err := getZoneResourceContainer(api, ctx, env.Cloudflare.Zone)
    if err != nil {
        return err
    }

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
    proxied := false
    if recordID == "" {
        _, err = api.CreateDNSRecord(ctx, rc, cloudflare.CreateDNSRecordParams{Content: vps.IP, Name: env.Server.FQDN, Proxied: &proxied, TTL: 60, Type: "A"})
    } else {
        _, err = api.UpdateDNSRecord(ctx, rc, cloudflare.UpdateDNSRecordParams{Content: vps.IP, ID: recordID, Proxied: &proxied, TTL: 60})
    }

    if err != nil {
        return fmt.Errorf("unable to set dns record: %v", err)
    }

    return nil
}

// Retrieve the IP address for a DNS record=
func Retrieve(env env.Env, fqdn string) (string, error) {
    api, err := cloudflare.NewWithAPIToken(env.Cloudflare.ApiKey)
    if err != nil {
        return "", fmt.Errorf("unable to connect to cloudflare api: %v", err)
    }

    ctx := context.Background()

    rc, err := getZoneResourceContainer(api, ctx, env.Cloudflare.Zone)
    if err != nil {
        return "", err
    }

    records, _, err := api.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{Name: fqdn})
    if err != nil {
        return "", fmt.Errorf("unable to list dns records for %s: %v", env.Cloudflare.Zone, err)
    }

    for _, record := range records {
        if strings.EqualFold(record.Name, fqdn) {
            return record.Content, nil
        }
    }

    return "", fmt.Errorf("unable to find dns record for %s", fqdn)
}

// getZoneResourceContainers finds the resource container for a zone name
func getZoneResourceContainer(api *cloudflare.API, ctx context.Context, fqzn string) (*cloudflare.ResourceContainer, error) {
    zones, err := api.ListZones(ctx, fqzn)
    if err != nil {
        return nil, fmt.Errorf("unable to list dns zones: %v", err)
    }

    var zoneID string
    for _, zone := range zones {
        if strings.EqualFold(zone.Name, fqzn) {
            zoneID = zone.ID
        }
    }

    if zoneID == "" {
        return nil, fmt.Errorf("unable to determine zone id for %s", fqzn)
    }

    return cloudflare.ZoneIdentifier(zoneID), nil
}
