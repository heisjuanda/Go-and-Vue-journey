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

	router.Get("/"+constants.SEARCH_ENDPOINT, helpers.SearchHandler)
	router.Get("/"+constants.SINGLE_MAIL_ENDPOINT, helpers.GetSingleEmailHandler)
	handler := helpers.SetCorsProtocol(router)

	initText := fmt.Sprintf("Server started on :%s", constants.SEARCH_SERVER_PORT)
	fmt.Println(initText)

	listenAndServeText := fmt.Sprintf(":%s", constants.SEARCH_SERVER_PORT)
	http.ListenAndServe(listenAndServeText, handler)
}
