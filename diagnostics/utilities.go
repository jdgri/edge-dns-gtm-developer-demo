package diagnostics

import "net/http"

func setEncodedHeader(req *http.Request) {
	if req.Method == "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return
}
