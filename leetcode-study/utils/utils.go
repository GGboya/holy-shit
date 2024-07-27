package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"leetcode/config" // 替换为实际的路径
)

func SendRequest(payload map[string]interface{}, headers map[string]string) (string, error) {
	url := "https://leetcode.cn/graphql/noj-go/"

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create new HTTP request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func FetchLastSubmitTime() (*time.Time, error) {
	payload := map[string]interface{}{
		"query": `
			query recentAcSubmissions($userSlug: String!) {
				recentACSubmissions(userSlug: $userSlug) {
					submissionId
					submitTime
					question {
						translatedTitle
						titleSlug
						questionFrontendId
					}
				}
			}
		`,
		"variables": map[string]string{
			"userSlug": config.UserSlug,
		},
	}

	submissionResult, err := SendRequest(payload, config.Headers)
	if err != nil {
		return nil, err
	}

	idx := strings.Index(submissionResult, "submitTime")
	if idx == -1 {
		return nil, fmt.Errorf("no submissions found")
	}
	lastTimeStr := submissionResult[idx+12 : idx+22]
	lastTime, err := strconv.ParseInt(lastTimeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	lastSubmitTime := time.Unix(lastTime, 0)
	return &lastSubmitTime, nil
}
