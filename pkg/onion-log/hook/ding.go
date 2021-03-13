package hook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DingMsg struct {
	Msgtype  string           `json:"msgtype"` // enum["text", "link", "markdown"]
	Text     *DingMsgText     `json:"text,omitempty"`
	Markdown *DingMsgMarkdown `json:"markdown,omitempty"`
	Link     *DingMsgLink     `json:"link,omitempty"`
	At       *DingAt          `json:"at,omitempty"`
}

type DingMsgText struct {
	Content string `json:"content"`
}

type DingMsgLink struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicURL     string `json:"picUrl"`
	MessageURL string `json:"messageUrl"`
}

type DingMsgMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Ding struct {
	Sign        string `json:"sign"`
	AccessToken string `json:"access_token"`
}

// addSign 增加验签
// 第一步，把timestamp+"\n"+密钥当做签名字符串
// 使用HmacSHA256算法计算签名，ok
// 然后进行Base64 encode， ok
// 最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集） failed
func (d *Ding) addSign(timestamp int64) string {
	str := fmt.Sprintf("%d\n%s", timestamp, d.Sign)

	h := hmac.New(sha256.New, []byte(d.Sign))
	h.Write([]byte(str))

	baseStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(baseStr)
}

// Send 发送消息
func (d *Ding) Send(msg *DingMsg) error {

	timestamp := time.Now().Unix() * 1000

	method := "POST"
	URL := fmt.Sprintf(
		"https://oapi.dingtalk.com/robot/send?access_token=%s&sign=%s&timestamp=%d",
		d.AccessToken, d.addSign(timestamp), timestamp)

	b, _ := json.Marshal(msg)
	payload := strings.NewReader(string(b))

	client := &http.Client{}
	req, err := http.NewRequest(method, URL, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()
	body, err := ioutil.ReadAll(res.Body)

	result := &DingResponse{}
	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return errors.New(result.ErrMsg)
	}
	return nil
}
