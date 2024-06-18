package anticaptcha

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Client struct {
	ClientKey                   string
	WebsiteSToken               *string
	RecaptchaDataSValue         *string
	ConnectionTimeout           int
	FirstAttemptWaitingInterval int
	NormalWaitingInterval       int
	IsVerbose                   bool
	TaskID                      int
	FuncaptchaApiJSSubdomain    *string
	FuncaptchaDataBlob          *string
	SoftId                      int
	OSTronAddress               string
	HcaptchaUserAgent           *string
	HcaptchaRespKey             *string
	Cookies                     *[]string
}

type ImageSettings struct {
	Phrase        bool
	CaseSensitive bool
	Numeric       bool
	MathOperation bool
	MinLength     int
	MaxLength     int
	LanguagePool  string
	Comment       string
}

type RecaptchaV2 struct {
	WebsiteURL    string
	WebsiteKey    string
	WebsiteSToken string
	IsInvisible   bool
	DataSValue    string
}

func NewClient(apiKey string) *Client {
	ac := Client{
		ClientKey:                   apiKey,
		ConnectionTimeout:           120,
		FirstAttemptWaitingInterval: 5,
		NormalWaitingInterval:       5,
		IsVerbose:                   true,
		SoftId:                      1187,
	}
	return &ac
}

func (ac *Client) SetAPIKey(key string) {
	ac.ClientKey = key
}

func (ac *Client) ShutUp() {
	ac.IsVerbose = false
}

func (ac *Client) SetSoftId(softId int) {
	ac.SoftId = softId
}

func (ac *Client) GetBalance() (float64, error) {
	response, err := ac.JSONRequest("getBalance", map[string]interface{}{"clientKey": ac.ClientKey})
	if err != nil {
		return 0, err
	}
	return response["balance"].(float64), nil
}

func (ac *Client) GetCreditsBalance() (float64, error) {
	response, err := ac.JSONRequest("getBalance", map[string]interface{}{"clientKey": ac.ClientKey})
	if err != nil {
		return 0, err
	}
	if credits, ok := response["captchaCredits"].(float64); ok {
		return credits, nil
	}
	return 0, nil
}

func (ac *Client) SolveImageFile(path string, settings ImageSettings) (string, error) {
	imageData, err := readImageFile(path)
	if err != nil {
		return "", err
	}
	return ac.SolveImage(base64.StdEncoding.EncodeToString(imageData), settings)
}

func (ac *Client) SolveImage(body string, settings ImageSettings) (string, error) { //, phrase bool, caseSensitive bool, isNumeric bool
	task := map[string]interface{}{
		"type":         "ImageToTextTask",
		"body":         body,
		"phrase":       settings.Phrase,
		"case":         settings.CaseSensitive,
		"numeric":      settings.Numeric,
		"comment":      settings.Comment,
		"math":         settings.MathOperation,
		"minLength":    settings.MinLength,
		"maxLength":    settings.MaxLength,
		"languagePool": settings.LanguagePool,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["text"].(string), nil
}

func (ac *Client) ReportIncorrectImageCaptcha() error {
	_, err := ac.JSONRequest("reportIncorrectImageCaptcha", map[string]interface{}{
		"clientKey": ac.ClientKey,
		"taskId":    ac.TaskID,
	})
	return err
}

func (ac *Client) SolveRecaptchaV2Proxyless(recaptcha RecaptchaV2) (string, error) {
	task := map[string]interface{}{
		"type":                "RecaptchaV2TaskProxyless",
		"websiteURL":          recaptcha.WebsiteURL,
		"websiteKey":          recaptcha.WebsiteKey,
		"websiteSToken":       recaptcha.WebsiteSToken,
		"recaptchaDataSValue": recaptcha.DataSValue,
	}
	if recaptcha.IsInvisible {
		task["isInvisible"] = true
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	if cookies, ok := solution["cookies"].([]string); ok {
		ac.Cookies = &cookies
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func CreateTaskAndWaitForResult(ac *Client, task map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"clientKey": ac.ClientKey,
		"task":      task,
		"softId":    ac.SoftId,
	}
	if task["languagePool"] != nil {
		payload["languagePool"] = task["languagePool"]
	}
	taskCreateResult, err := ac.JSONRequest("createTask", payload)
	if err != nil {
		return nil, err
	}
	if taskID, ok := taskCreateResult["taskId"].(float64); ok {
		ac.TaskID = int(taskID)
		solution, err := ac.WaitForResult(ac.TaskID)
		if err != nil {
			return nil, err
		}
		return solution, nil
	}
	return nil, errors.New(taskCreateResult["errorCode"].(string))
}

func (ac *Client) GetCookies() *[]string {
	return ac.Cookies
}

func (ac *Client) ReportIncorrectRecaptcha() error {
	_, err := ac.JSONRequest("reportIncorrectRecaptcha", map[string]interface{}{
		"clientKey": ac.ClientKey,
		"taskId":    ac.TaskID,
	})
	return err
}

func (ac *Client) ReportCorrectRecaptcha() error {
	_, err := ac.JSONRequest("reportCorrectRecaptcha", map[string]interface{}{
		"clientKey": ac.ClientKey,
		"taskId":    ac.TaskID,
	})
	return err
}

func (ac *Client) WaitForResult(taskId int) (map[string]interface{}, error) {
	if ac.IsVerbose {
		fmt.Println("created task with ID", taskId)
		fmt.Println("waiting", ac.FirstAttemptWaitingInterval, "seconds")
	}
	time.Sleep(time.Duration(ac.FirstAttemptWaitingInterval) * time.Second)

	for taskId > 0 {
		checkResult, err := ac.JSONRequest("getTaskResult", map[string]interface{}{
			"clientKey": ac.ClientKey,
			"taskId":    taskId,
		})
		if err != nil {
			return nil, err
		}
		if status, ok := checkResult["status"].(string); ok && status == "ready" {
			return checkResult["solution"].(map[string]interface{}), nil
		}
		if status, ok := checkResult["status"].(string); ok && status == "processing" && ac.IsVerbose {
			fmt.Println("captcha result is not yet ready")
		}
		if ac.IsVerbose {
			fmt.Println("waiting", ac.NormalWaitingInterval, "seconds")
		}
		time.Sleep(time.Duration(ac.NormalWaitingInterval) * time.Second)
	}
	return nil, errors.New("ERROR_NO_SLOT_AVAILABLE")
}

func (ac *Client) JSONRequest(methodName string, payload map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.anti-captcha.com/" + methodName

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Duration(ac.ConnectionTimeout) * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response["errorId"] == nil {
		return nil, errors.New("Incorrect API response, something is wrong")
	}
	if response["errorId"].(float64) > 0 {
		if ac.IsVerbose {
			fmt.Println("Received API error", response["errorCode"], ":", response["errorDescription"])
		}
		return nil, errors.New(response["errorCode"].(string))
	}
	return response, nil
}

func readImageFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize < 100 {
		return nil, errors.New("Captcha file is too small")

	}
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
