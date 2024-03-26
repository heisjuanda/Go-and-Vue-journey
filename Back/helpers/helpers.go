package helpers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"main.go/constants"
)

func SetCorsProtocol(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)

		handler.ServeHTTP(w, r)
	})
}

func SearchHandler(w http.ResponseWriter, request *http.Request) {
	searchType := request.URL.Query().Get("type")
	if searchType == "" {
		http.Error(w, "Search term is empty", http.StatusBadRequest)
		fmt.Println("Search term is empty")
		return
	}

	searchTerm := request.URL.Query().Get("term")

	page := request.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}
	order := request.URL.Query().Get("order")
	if len(order) > 2 {
		order = ""
	}

	query := fmt.Sprintf(
		`{
        	"search_type": "%s",
        	"query": {
            	"term": "%s"
        	},
			"sort_fields": ["%sdate"],
        	"from": %s,
        	"max_results": 20,
        	"_source": []
    	}`,
		searchType, searchTerm, order, page)

	requestUrl := constants.ZINCSEARCH_SERVER + "/" + constants.ZINCSEARCH_ENDPOINT
	zincRequest, zincError := http.NewRequest(constants.METHOD_POST, requestUrl, strings.NewReader(query))
	if zincError != nil {
		if urlErr, ok := zincError.(*url.Error); ok && urlErr.Timeout() {
			fmt.Println("Request timed out:", zincError)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr != nil {
			fmt.Println("URL error:", urlErr)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		} else {
			fmt.Println("Other error:", zincError)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	zincRequest.SetBasicAuth(constants.ZINCSEARCH_USER_NAME, constants.ZINCSEARCH_PASSWORD)

	zincRequest.Header.Set("Content-Type", "application/json")
	zincRequest.Header.Set("User-Agent", constants.INFO_BROWSERS)

	zincResponse, zincResponseError := http.DefaultClient.Do(zincRequest)
	if zincResponseError != nil {
		if netErr, ok := zincResponseError.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Request timed out:", zincResponseError)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr, ok := zincResponseError.(*url.Error); ok && urlErr.Timeout() {
			fmt.Println("Request timed out:", zincResponseError)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if urlErr != nil {
			fmt.Println("URL error:", urlErr)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		} else {
			fmt.Println("Other error:", zincResponseError)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	defer zincResponse.Body.Close()

	body, bodyReadError := io.ReadAll(zincResponse.Body)
	if bodyReadError != nil {
		fmt.Println(bodyReadError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(zincResponse.StatusCode)
	_, bodyReadError = w.Write(body)
	if bodyReadError != nil {
		fmt.Println(bodyReadError)
	}
}
