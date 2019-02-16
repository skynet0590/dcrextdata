package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

const retryDelay = 60
const maxRetryAttempts = 3

// ResponseToType attempts to collect json data from the given url string and decode it into
// the destination
func ResponseToType(client *http.Client, url string, destination interface{}) error {
	resp, err := client.Get(url)
	//fmt.Println("Url: " + url)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		retryAttempts := 1

		ticker := time.NewTicker(retryDelay * time.Second)

		for range ticker.C {
			retryResp, err := client.Get(url)

			if err != nil {
				if retryResp != nil {
					retryResp.Body.Close()
				}

				retryAttempts++
				if retryAttempts > maxRetryAttempts {
					return err
				}
				continue
			}
			resp = retryResp
		}
	}

	err = json.NewDecoder(resp.Body).Decode(destination)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}
