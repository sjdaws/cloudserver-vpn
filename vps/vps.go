package vps

import (
    "io"
)

const apiURL = "https://cloudserver.nz/api/v1"

// closeBody closes a ReadCloser ignoring errors
func closeBody(body io.ReadCloser) {
    _ = body.Close()
}
