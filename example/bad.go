package test

import "net/http"

func fetch() error {

	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_ = resp.StatusCode
	return nil

}
