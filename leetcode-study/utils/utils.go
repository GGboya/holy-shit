package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"leetcode/config" // 替换为实际的路径
	"leetcode/entities"

	"gopkg.in/gomail.v2"
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

func FetchLastSubmitTime(ID string) (*time.Time, error) {
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
			"userSlug": ID,
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

func ExtractUserFromURL(url string) (string, error) {
	// 正则表达式模式：匹配以/u/开头的部分并捕获接下来的非/字符
	re := regexp.MustCompile(`https://leetcode\.cn/u/([^/]+)/?`)

	// 使用正则表达式查找匹配的部分
	match := re.FindStringSubmatch(url)

	// 检查是否有匹配的部分
	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("无法从URL提取用户标识符: %s", url)
}

func ConvertLevelToInt(level string) int {
	x, err := strconv.Atoi(level)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return x
}

func ConvrtUserFormatByteToNormal(user []byte) (*entities.User, error) {
	var resp *entities.User
	err := json.Unmarshal(user, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}
	return resp, nil
}

func ConvrtUserFormatByteToSecret(user []byte) (*entities.UserSecret, error) {
	var resp *entities.UserSecret
	err := json.Unmarshal(user, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}
	return resp, nil
}

func SendEmail(qq, title, content string) error {
	// 从环境变量中获取邮箱信息
	qqEmail := os.Getenv("QQ_EMAIL")
	qqAuthCode := os.Getenv("QQ_AUTH_CODE")

	// 验证环境变量是否加载成功
	if qqEmail == "" || qqAuthCode == "" {
		log.Fatal("QQ_EMAIL or QQ_AUTH_CODE is not set in .env file")
	}

	m := gomail.NewMessage()

	// 发件人
	m.SetHeader("From", qqEmail)

	// 收件人
	m.SetHeader("To", qq+"@qq.com")

	// 邮件标题
	m.SetHeader("Subject", title)

	// 邮件内容
	m.SetBody("text/plain", content)

	// QQ邮箱SMTP服务器信息
	d := gomail.NewDialer("smtp.qq.com", 587, qqEmail, qqAuthCode)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}
