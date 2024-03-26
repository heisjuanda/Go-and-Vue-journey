package constants

import "os"

var ZINCSEARCH_USER_NAME string = os.Getenv("ZINCSEARCH_USER_NAME")

var ZINCSEARCH_PASSWORD string = os.Getenv("ZINCSEARCH_PASSWORD")

var INDEXER_NAME string = os.Getenv("INDEX_NAME")

const INFO_BROWSERS string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"

const METHOD_POST string = "POST"

var ZINCSEARCH_SERVER string = os.Getenv("ZINCSEARCH_SERVER")
var ZINCSEARCH_ENDPOINT string = os.Getenv("ZINCSEARCH_ENDPOINT")

var SEARCH_SERVER_PORT string = os.Getenv("SEARCH_SERVER_PORT")
var SEARCH_ENDPOINT string = os.Getenv("SEARCH_ENDPOINT")
