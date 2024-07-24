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
	ConnectionTimeout           int
	FirstAttemptWaitingInterval int
	NormalWaitingInterval       int
	IsVerbose                   bool
	TaskID                      int

	SoftId            int
	HcaptchaUserAgent string
	HcaptchaRespKey   string
	Cookies           []string
}

type ImageSettings struct {
	Phrase        bool
	CaseSensitive bool
	Numeric       int
	MathOperation bool
	MinLength     int
	MaxLength     int
	LanguagePool  string
	Comment       string
	WebsiteURL    string
}

type Proxy struct {
	Type      string
	IPAddress string
	Port      int
	Login     string
	Password  string
}

type RecaptchaV2 struct {
	WebsiteURL    string
	WebsiteKey    string
	WebsiteSToken string
	IsInvisible   bool
	DataSValue    string
	UserAgent     string
	Proxy         *Proxy
}

type RecaptchaV3 struct {
	WebsiteURL   string
	WebsiteKey   string
	MinScore     float64
	PageAction   string
	IsEnterprise bool
	APIDomain    string
}

type Hcaptcha struct {
	WebsiteURL        string
	WebsiteKey        string
	IsInvisible       bool
	IsEnterprise      bool
	EnterprisePayload map[string]interface{}
	Proxy             *Proxy
}

type FunCaptcha struct {
	WebsiteURL       string
	WebsitePublicKey string
	ApiSubdomain     string
	DataBlob         string
	Proxy            *Proxy
}

type Turnstile struct {
	WebsiteURL string
	WebsiteKey string
	Action     string
	CData      string
	Proxy      *Proxy
}

type GeeTest struct {
	WebsiteURL     string
	Gt             string
	Challenge      string
	ApiSubdomain   string
	Version        int
	InitParameters map[string]interface{}
	Proxy          *Proxy
}

type AntiGate struct {
	WebsiteURL        string
	TemplateName      string
	Variables         map[string]interface{}
	DomainsOfInterest []string
	Proxy             *Proxy
}

type ImageToCoordinates struct {
	Mode       string
	Comment    string
	WebsiteURL string
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
	imageData, err := ac.ReadImageFile(path)
	if err != nil {
		return "", err
	}
	return ac.SolveImage(base64.StdEncoding.EncodeToString(imageData), settings)
}

func (ac *Client) SolveImage(body string, settings ImageSettings) (string, error) {
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

func (ac *Client) SolveRecaptchaV2(recaptcha RecaptchaV2) (string, error) {
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
		ac.Cookies = cookies
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func (ac *Client) SolveRecaptchaV2ProxyOn(recaptcha RecaptchaV2) (string, error) {
	task := map[string]interface{}{
		"type":                "RecaptchaV2Task",
		"websiteURL":          recaptcha.WebsiteURL,
		"websiteKey":          recaptcha.WebsiteKey,
		"websiteSToken":       recaptcha.WebsiteSToken,
		"recaptchaDataSValue": recaptcha.DataSValue,
		"userAgent":           recaptcha.UserAgent,
		"proxyType":           recaptcha.Proxy.Type,
		"proxyAddress":        recaptcha.Proxy.IPAddress,
		"proxyPort":           recaptcha.Proxy.Port,
		"proxyLogin":          recaptcha.Proxy.Login,
		"proxyPassword":       recaptcha.Proxy.Password,
	}
	if recaptcha.IsInvisible {
		task["isInvisible"] = true
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	if cookies, ok := solution["cookies"].([]string); ok {
		ac.Cookies = cookies
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func (ac *Client) SolveRecaptchaV3(recaptcha RecaptchaV3) (string, error) {
	task := map[string]interface{}{
		"type":         "RecaptchaV3TaskProxyless",
		"websiteURL":   recaptcha.WebsiteURL,
		"websiteKey":   recaptcha.WebsiteKey,
		"minScore":     recaptcha.MinScore,
		"pageAction":   recaptcha.PageAction,
		"isEnterprise": recaptcha.IsEnterprise,
		"apiDomain":    recaptcha.APIDomain,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func (ac *Client) SolveHcaptcha(hcaptcha Hcaptcha) (string, error) {
	task := map[string]interface{}{
		"type":              "HCaptchaTaskProxyless",
		"websiteURL":        hcaptcha.WebsiteURL,
		"websiteKey":        hcaptcha.WebsiteKey,
		"isEnterprise":      hcaptcha.IsEnterprise,
		"enterprisePayload": hcaptcha.EnterprisePayload,
	}
	if hcaptcha.IsInvisible {
		task["isInvisible"] = true
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	if userAgent, ok := solution["userAgent"].(string); ok {
		ac.HcaptchaUserAgent = userAgent
	}
	if respKey, ok := solution["respKey"].(string); ok {
		ac.HcaptchaRespKey = respKey
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func (ac *Client) SolveHcaptchaProxyOn(hcaptcha Hcaptcha) (string, error) {
	task := map[string]interface{}{
		"type":              "HCaptchaTask",
		"websiteURL":        hcaptcha.WebsiteURL,
		"websiteKey":        hcaptcha.WebsiteKey,
		"isEnterprise":      hcaptcha.IsEnterprise,
		"enterprisePayload": hcaptcha.EnterprisePayload,
		"proxyType":         hcaptcha.Proxy.Type,
		"proxyAddress":      hcaptcha.Proxy.IPAddress,
		"proxyPort":         hcaptcha.Proxy.Port,
		"proxyLogin":        hcaptcha.Proxy.Login,
		"proxyPassword":     hcaptcha.Proxy.Password,
	}
	if hcaptcha.IsInvisible {
		task["isInvisible"] = true
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	if userAgent, ok := solution["userAgent"].(string); ok {
		ac.HcaptchaUserAgent = userAgent
	}
	if respKey, ok := solution["respKey"].(string); ok {
		ac.HcaptchaRespKey = respKey
	}
	return solution["gRecaptchaResponse"].(string), nil
}

func (ac *Client) SolveFunCaptcha(funcaptcha FunCaptcha) (string, error) {
	task := map[string]interface{}{
		"type":                     "FunCaptchaTaskProxyless",
		"websiteURL":               funcaptcha.WebsiteURL,
		"websitePublicKey":         funcaptcha.WebsitePublicKey,
		"funcaptchaApiJSSubdomain": funcaptcha.ApiSubdomain,
		"data":                     funcaptcha.DataBlob,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["token"].(string), nil
}

func (ac *Client) SolveFunCaptchaProxyOn(funcaptcha FunCaptcha) (string, error) {
	task := map[string]interface{}{
		"type":                     "FunCaptchaTask",
		"websiteURL":               funcaptcha.WebsiteURL,
		"websitePublicKey":         funcaptcha.WebsitePublicKey,
		"funcaptchaApiJSSubdomain": funcaptcha.ApiSubdomain,
		"data":                     funcaptcha.DataBlob,
		"proxyType":                funcaptcha.Proxy.Type,
		"proxyAddress":             funcaptcha.Proxy.IPAddress,
		"proxyPort":                funcaptcha.Proxy.Port,
		"proxyLogin":               funcaptcha.Proxy.Login,
		"proxyPassword":            funcaptcha.Proxy.Password,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["token"].(string), nil
}

func (ac *Client) SolveTurnstile(turnstile Turnstile) (string, error) {
	task := map[string]interface{}{
		"type":           "TurnstileTaskProxyless",
		"websiteURL":     turnstile.WebsiteURL,
		"websiteKey":     turnstile.WebsiteKey,
		"action":         turnstile.Action,
		"turnstileCData": turnstile.CData,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["token"].(string), nil
}

func (ac *Client) SolveTurnstileProxyOn(turnstile Turnstile) (string, error) {
	task := map[string]interface{}{
		"type":           "TurnstileTask",
		"websiteURL":     turnstile.WebsiteURL,
		"websiteKey":     turnstile.WebsiteKey,
		"action":         turnstile.Action,
		"turnstileCData": turnstile.CData,
		"proxyType":      turnstile.Proxy.Type,
		"proxyAddress":   turnstile.Proxy.IPAddress,
		"proxyPort":      turnstile.Proxy.Port,
		"proxyLogin":     turnstile.Proxy.Login,
		"proxyPassword":  turnstile.Proxy.Password,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return "", err
	}
	return solution["token"].(string), nil
}

func (ac *Client) SolveGeeTest(geetest GeeTest) (map[string]interface{}, error) {
	task := map[string]interface{}{
		"type":                      "GeeTestTaskProxyless",
		"websiteURL":                geetest.WebsiteURL,
		"gt":                        geetest.Gt,
		"challenge":                 geetest.Challenge,
		"geetestApiServerSubdomain": geetest.ApiSubdomain,
		"version":                   geetest.Version,
		"initParameters":            geetest.InitParameters,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return solution, nil
}

func (ac *Client) SolveGeeTestProxyOn(geetest GeeTest) (map[string]interface{}, error) {
	task := map[string]interface{}{
		"type":                      "GeeTestTask",
		"websiteURL":                geetest.WebsiteURL,
		"gt":                        geetest.Gt,
		"challenge":                 geetest.Challenge,
		"geetestApiServerSubdomain": geetest.ApiSubdomain,
		"version":                   geetest.Version,
		"initParameters":            geetest.InitParameters,
		"proxyType":                 geetest.Proxy.Type,
		"proxyAddress":              geetest.Proxy.IPAddress,
		"proxyPort":                 geetest.Proxy.Port,
		"proxyLogin":                geetest.Proxy.Login,
		"proxyPassword":             geetest.Proxy.Password,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return solution, nil
}

func (ac *Client) SolveAntiGate(antigate AntiGate) (map[string]interface{}, error) {
	task := map[string]interface{}{
		"type":              "AntiGateTask",
		"websiteURL":        antigate.WebsiteURL,
		"templateName":      antigate.TemplateName,
		"variables":         antigate.Variables,
		"domainsOfInterest": antigate.DomainsOfInterest,
		"proxyType":         antigate.Proxy.Type,
		"proxyAddress":      antigate.Proxy.IPAddress,
		"proxyPort":         antigate.Proxy.Port,
		"proxyLogin":        antigate.Proxy.Login,
		"proxyPassword":     antigate.Proxy.Password,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return solution, nil
}

func (ac *Client) SolveImageToCoordinates(body string, settings ImageToCoordinates) ([]interface{}, error) { //, phrase bool, caseSensitive bool, isNumeric bool
	task := map[string]interface{}{
		"type":       "ImageToCoordinatesTask",
		"body":       body,
		"comment":    settings.Comment,
		"mode":       settings.Mode,
		"websiteURL": settings.WebsiteURL,
	}
	solution, err := CreateTaskAndWaitForResult(ac, task)
	if err != nil {
		return []interface{}{}, err
	}
	return solution["coordinates"].([]interface{}), nil
}

func CreateTaskAndWaitForResult(ac *Client, task map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"clientKey": ac.ClientKey,
		"task":      task,
		"softId":    ac.SoftId,
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

func (ac *Client) GetCookies() []string {
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

func (ac *Client) ReadImageFile(filePath string) ([]byte, error) {
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
