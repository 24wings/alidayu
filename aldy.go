package alidayu

import (
	// "strings"
	"fmt"
	"crypto/hmac"
	"encoding/base64"

	"crypto/sha1"
	"encoding/json"

	"io/ioutil"

	"net/http"
	"net/url"
	"sort"
	// "strings"
	"time"

	// "github.com/satori/go.uuid"
)

// sendSmsResponse
type sendSmsResponse struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

const (
	dyURL = "http://dysmsapi.aliyuncs.com"
)

// signHMAC 获取签名
func signHMAC(params url.Values, appSecret string) (signature string) {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	str := ""
	sort.Strings(keys)
	for _, k := range keys {
		str += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(params.Get(k))
	}
	signstr := "GET&%2F&" + url.QueryEscape(str[1:])
	mac := hmac.New(sha1.New, []byte(appSecret+"&"))
	mac.Write([]byte(signstr))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type SMSSendDetail struct{
	SendDate time.Time
	SendStatus int
	ReceiveDate time.Time
	ErrorCode string
	TemplateCode string
	Content string
	PhoneNum string
	OutId string
}

type QuerySendDetailsSuccessResponse struct{
	TotalCount int
	Message string
	RequestId string
	SmsSendDetailDTOs struct{SmsSendDetailDTO []SMSSendDetail}
	Code string
}


// SendSMS
func SendSMS(mobileNo, signName, templateCode, paramString, appKey, appSecret string) (bool, string, error) {
	params := url.Values{}

	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureVersion", "1.0")
	params.Set("SignatureNonce", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("AccessKeyId", appKey)
	params.Add("Format", "JSON")
	params.Set("RegionId", "cn-hangzhou")
	// var code TemplateCode;
	// json.Unmarshal([]byte(templateCode), )
	// 存储用户的真实验证码
	params.Set("OutId", paramString)
	paramString =`{"code":"`+paramString+`"}`
	
	params.Set("SignName", signName)
	params.Set("TemplateCode", templateCode)
	params.Set("TemplateParam", paramString)

	params.Set("Action", "SendSms")
	params.Set("PhoneNumbers", mobileNo)
	params.Set("Version", "2017-05-25")

	signstr := signHMAC(params, appSecret)
	params.Set("Signature", signstr)
	req, err := http.NewRequest(http.MethodGet, dyURL+"/?"+params.Encode(), nil)

	req.Header.Set("x-sdk-client", "Java/2.0.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Java/1.6.0_45")

	c := new(http.Client)
	resp, err := c.Do(req)

	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}
	var result QuerySendDetailsSuccessResponse
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return false, "", err
	}
	return result.Code == "OK", result.Message, nil
}
func QueryDetail(PhoneNumber string,signName string,sendDate string ,templateCode, appKey, appSecret string)(bool,SMSSendDetail,string){
	params := url.Values{}

	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureVersion", "1.0")
	params.Set("SignatureNonce", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("AccessKeyId", appKey)
	params.Add("Format", "JSON")
	params.Set("RegionId", "cn-hangzhou")


	params.Set("SignName", signName)
	params.Set("TemplateCode", templateCode)
	params.Set("sendStatus", "3")
	// params.Set("TemplateParam", paramString)
	// params.Set("OutId", "")
	params.Set("SendDate", sendDate)
	params.Set("ReceiveDate", "2018-03-22")
	params.Set("Action", "QuerySendDetails")
	params.Set("PhoneNumber", PhoneNumber)
	params.Set("Version", "2017-05-25")
	params.Set("PageSize", "10")
	params.Set("CurrentPage","1")
	

	signstr := signHMAC(params, appSecret)
	params.Set("Signature", signstr)
	req, err := http.NewRequest(http.MethodGet, dyURL+"/?"+params.Encode(), nil)

	req.Header.Set("x-sdk-client", "Java/2.0.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Java/1.6.0_45")

	c := new(http.Client)
	resp, err := c.Do(req)

	if err != nil {
		return false, SMSSendDetail{}, err.Error()
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, SMSSendDetail{}, err.Error()
	}
	var result QuerySendDetailsSuccessResponse
	err = json.Unmarshal(bs, &result)
	fmt.Println(result,result.Code,string(bs) );
	if(len(result.SmsSendDetailDTOs.SmsSendDetailDTO)>0){
	
		return true, result.SmsSendDetailDTOs.SmsSendDetailDTO[0], "发送成功"
	}else{
		return false,SMSSendDetail{},"近期未发送手机消息"
	}

	
	
	
}