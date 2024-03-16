package helpers

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"main.go/constants"
)

func SetCorsProtocol(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)

		next.ServeHTTP(w, r)
	})
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	searchTerm := r.URL.Query().Get("term")
	if searchTerm == "" {
		http.Error(w, "Search term is empty", http.StatusBadRequest)
		log.Println("Search term is empty")
		return
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}
	order := r.URL.Query().Get("order")
	if len(order) > 2 {
		order = ""
	}

	query := fmt.Sprintf(
		`{
        	"search_type": "matchphrase",
        	"query": {
            	"term": "%s"
        	},
			"sort_fields": ["%sdate"],
        	"from": %s,
        	"max_results": 20,
        	"_source": []
    	}`,
		searchTerm, order, page)

	reqUrl := constants.Server + constants.Endpoint
	req, err := http.NewRequest(constants.MethodPost, reqUrl, strings.NewReader(query))
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			log.Println("Request timed out:", err)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr != nil {
			log.Println("URL error:", urlErr)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		} else {
			log.Println("Other error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	req.SetBasicAuth(constants.UseName, constants.Password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", constants.InfoBrowsers)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("Request timed out:", err)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			log.Println("Request timed out:", err)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr != nil {
			log.Println("URL error:", urlErr)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		} else {
			log.Println("Other error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
