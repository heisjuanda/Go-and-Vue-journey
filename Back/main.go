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
