package fetchzinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mail-indexer/constants"
	"mail-indexer/models"
	"net/http"
	"time"
)

// Creates the Index for email-indexer (THIS ALSO COULD BE DONE MANUALLY :/)
func FetchCreateZincIndex() {
	username := constants.USER_NAME
	password := constants.PASSWORD

	body := constants.EMAIL_INDEXER_INDEX

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	url := constants.SERVER + constants.ENDPOINT_INDEX
	req, err := http.NewRequest(constants.METHOD_POST, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}
}

// receives an array of emails and then send them to zincSearch
func FetchZing(allEmails []models.Email) {
	username := constants.USER_NAME
	password := constants.PASSWORD

	var body models.ZincRequest
	body.Index = constants.INDEXER_NAME
	body.Records = allEmails

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	url := constants.SERVER + constants.ENDPOINT
	req, err := http.NewRequest(constants.METHOD_POST, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		return
	}

	fmt.Println("start fetch: ", start)
	fmt.Println("end fetch: ", time.Since(start))
}
