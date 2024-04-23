package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rate-limiter/api/types"
	"rate-limiter/pkg/constants"
	"time"
)

func CallConcurrentReqs(ctx context.Context, endpoint string , users []string) {
	client := &http.Client{}
	reqUrl := constants.BaseURL + endpoint

loop:
	for {
		for i := 0; i < GetRandomReqCount(); i++ {
			select {
			case <-ctx.Done():
				break loop
			default:
			}

			go func(uid string) {
				reqBody, err := json.Marshal(types.Request{UserID: uid})

				if err != nil {
					log.Printf("Error marshalling request: %v", err)
					return
				}

				req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(reqBody))
				if err != nil {
					log.Printf("Error creating request: %v", err)
					return
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error sending request for UserID %s: %v", uid, err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading response body: %v", err)
					return
				}

				var userData types.Response
				err = json.Unmarshal(body, &userData)
				if err != nil {
					log.Printf("Error unmarshalling response: %v", err)
					return
				}
			}(GetRandomUserID(users))
			time.Sleep(10 * time.Nanosecond)
		}
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			break loop
		default:

		}
	}
}
