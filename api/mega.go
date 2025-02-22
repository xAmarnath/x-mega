package handler

import (
	"io"
	"net/http"
	"time"
)

const baseURL = "https://megacloud.tv"

func ProxyHandle(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	reqURL := baseURL + r.URL.Path
	if query := r.URL.RawQuery; query != "" {
		reqURL += "?" + query
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://megacloud.tv/embed-1/e-1/muqImFVgS273?z='")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Copy cookies if present
	cookie := r.Header.Get("Cookie")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
