package integrationtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/msproject/relive/util"
)

var accountCreateURL = reliveTestCfg.reliveServerURL + "/api/accounts/create"

func createDummyAccount() ([]byte, error) {
	data := &util.CreateAccountReq{
		UserName:  "testact",
		Email:     "test@relive.com",
		FirstName: "testFirst",
		LastName:  "testLast",
		PWD:       "test001",
		Role:      0,
	}

	jsonObj, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//str := string(jsonObj)

	//return []byte(str), nil
	return jsonObj, nil
}

func testCreateAccount() error {
	jsonData, createErr := createDummyAccount()

	if createErr != nil {
		return fmt.Errorf("failed to create dummy recording! err: %v", createErr)
	}

	req, err := http.NewRequest("POST", accountCreateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create recording http request! err: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to create recording from API! err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to get 2XX response for create account %s", resp.Status)
	}

	return nil
}

func testCreateAccountAdmin() error {
	return nil
}

func testCreateAccountCustomer() error {
	return nil
}

func testCreateSubscription() error {
	return nil
}

func testUploadMedia() error {
	return nil
}
