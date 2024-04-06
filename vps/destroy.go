package vps

import (
    "fmt"
    "io"
    "net/http"
    "strings"

    "github.com/sjdaws/cloudserver-vpn/env"
)

// Destroy an existing virtual private server
func Destroy(env env.Env, serverID int) error {
    errs := env.ValidateDestroyEnv()
    if len(errs) > 0 {
        return fmt.Errorf("unable to remove server:\n - %s", strings.Join(errs, "\n - "))
    }

    request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/servers/%d", apiURL, serverID), nil)
    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.CloudServer.ApiKey))

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return fmt.Errorf("unable to remove server: %v", err)
    }
    defer closeBody(response.Body)

    body, _ := io.ReadAll(response.Body)

    if response.StatusCode != http.StatusOK {
        return fmt.Errorf("unable to remove server, invalid status: %s - %s", response.Status, string(body))
    }

    return nil
}

// ListActiveVPS returns the IDs of all the active VPS servers
func ListActiveVPS(env env.Env) ([]ServerData, error) {
    errs := env.ValidateDestroyEnv()
    if len(errs) > 0 {
        return nil, fmt.Errorf("unable to remove servers:\n - %s", strings.Join(errs, "\n - "))
    }

    projectID, err := findOrCreateProject(env)
    if err != nil {
        return nil, err
    }

    return listProjectVPS(env, projectID)
}
