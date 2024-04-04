package vps

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"

    "github.com/sjdaws/cloudserver-vpn/env"
)

type Project struct {
    Data        ProjectData `json:"data,omitempty"`
    Description string      `json:"description"`
    Name        string      `json:"name"`
}

type ProjectData struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type ProjectSearch struct {
    Data []ProjectData `json:"data"`
}

type ServerSearch struct {
    Data []ServerData `json:"data"`
}

const projectName = "VPNs"

// createProject creates a new project
func createProject(env env.Env) (int, error) {
    payload, err := json.Marshal(&Project{
        Description: "VPN servers created by cloudserver-vpn",
        Name:        projectName,
    })
    if err != nil {
        return 0, fmt.Errorf("unable to marshal new project configuration: %v", err)
    }

    request, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", apiURL), bytes.NewBuffer(payload))
    if err != nil {
        return 0, fmt.Errorf("unable to create new project request: %v", err)
    }

    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.CloudServer.ApiKey))
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return 0, fmt.Errorf("unable to create new project: %v", err)
    }
    defer closeBody(response.Body)

    body, _ := io.ReadAll(response.Body)

    if response.StatusCode != http.StatusCreated {
        return 0, fmt.Errorf("unable to create new project, invalid status: %s - %s", response.Status, string(body))
    }

    var result Project
    err = json.Unmarshal(body, &result)
    if err != nil {
        return 0, fmt.Errorf("error reading response from server: %v", err)
    }

    return result.Data.ID, nil
}

// findOrCreateProject attempts to find the project to use, creates it if it doesn't exist
func findOrCreateProject(env env.Env) (int, error) {
    // If project is set in env, use it
    if env.CloudServer.Project != 0 {
        return env.CloudServer.Project, nil
    }

    // Find project
    projectID, err := findProject(env)
    if err != nil {
        return 0, err
    }

    if projectID > 0 {
        return projectID, nil
    }

    // Create project
    return createProject(env)
}

// findProject attempts to find the project to use, creates it if it doesn't exist
func findProject(env env.Env) (int, error) {
    request, err := http.NewRequest("GET", fmt.Sprintf("%s/projects?filter[search]=%s", apiURL, projectName), nil)
    if err != nil {
        return 0, fmt.Errorf("unable to create project search request: %v", err)
    }

    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.CloudServer.ApiKey))

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return 0, fmt.Errorf("unable to perform project search: %v", err)
    }
    defer closeBody(response.Body)

    body, _ := io.ReadAll(response.Body)

    if response.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("unable to perform project search, invalid status: %s - %s", response.Status, string(body))
    }

    var result ProjectSearch
    err = json.Unmarshal(body, &result)
    if err != nil {
        return 0, fmt.Errorf("error reading response from server: %v", err)
    }

    for _, project := range result.Data {
        if strings.EqualFold(project.Name, projectName) {
            return project.ID, nil
        }
    }

    return 0, nil
}

// listProjectVPS lists all the servers in a project
func listProjectVPS(env env.Env, projectID int) ([]int, error) {
    request, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%d/servers", apiURL, projectID), nil)
    if err != nil {
        return nil, fmt.Errorf("unable to create server search request: %v", err)
    }

    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.CloudServer.ApiKey))

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("unable to perform server search: %v", err)
    }
    defer closeBody(response.Body)

    body, _ := io.ReadAll(response.Body)

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unable to perform server search, invalid status: %s - %s", response.Status, string(body))
    }

    var result ServerSearch
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("error reading response from server: %v", err)
    }

    var found []int
    for _, server := range result.Data {
        found = append(found, server.ID)
    }

    return found, nil
}
