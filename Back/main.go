package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"main.go/constants"
	"main.go/helpers"
)

func main() {
	router := chi.NewRouter()

	router.Get(constants.MyEndpoint, helpers.SearchHandler)
	handler := helpers.SetCorsProtocol(router)

	initText := fmt.Sprintf("Server started on :%s", constants.MyServerPort)
	fmt.Println(initText)

	listenAndServeText := fmt.Sprintf(":%s", constants.MyServerPort)
	http.ListenAndServe(listenAndServeText, handler)
}

// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strings"
// )

// func main() {
// 	query := `{
//         "search_type": "matchphrase",
//         "query":
//         {
//             "term": "hello"
//         },
//         "from": 9183,
//         "max_results": 20,
//         "_source": []
//     }`
// 	req, err := http.NewRequest("POST", "http://localhost:4080/api/email-indexer/_search", strings.NewReader(query))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req.SetBasicAuth("admin", "Complexpass#123")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	log.Println(resp.StatusCode)
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(body))
// }
