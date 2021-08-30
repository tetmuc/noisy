package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tetmuc/noisy/pkg/alert"
)

func NewFeishuRot(webhook string) alert.IAlert {
	feishuRot := FeishuRot{}
	feishuRot.webhook = webhook
	return &feishuRot
}

type FeishuRot struct {
	webhook string
}

type BaseContent struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type RichContent struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string                `json:"title"`
				Content [][]map[string]string `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

type feishuRequest struct {
	BaseContent
}

type feishuRichTextRequest struct {
	OpenId  string `json:"open_id,omitempty"`
	RootId  string `json:"root_id,omitempty"`
	ChatId  string `json:"chat_id,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	MsgType string `json:"msg_type"`
	RichContent
}

type feishuResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (f *FeishuRot) buildAtContents(nomineesList []string) []map[string]string {
	if len(nomineesList) > 0 {
		list := []map[string]string{}
		for _, nominee := range nomineesList {
			content := map[string]string{}
			content["tag"] = "at"
			content["user_id"] = nominee
			list = append(list, content)
		}
		return list
	}
	return nil
}

func (f *FeishuRot) AlertText(keyWord string, title, msg string, nominees ...string) error {
	feiStr := fmt.Sprintf(`%v`, msg)

	contentList := []map[string]string{}
	nominessContent := f.buildAtContents(nominees)

	if len(nominessContent) > 0 {
		contentList = nominessContent
	}
	msgContent := map[string]string{
		"tag":  "text",
		"text": feiStr,
	}

	var req feishuRichTextRequest
	req.MsgType = "post"
	req.Content.Post.ZhCn.Title = fmt.Sprintf("%v [%v]", title, keyWord)
	req.Content.Post.ZhCn.Content = [][]map[string]string{}
	contentList = append(contentList, msgContent)
	req.Content.Post.ZhCn.Content = append(req.Content.Post.ZhCn.Content, contentList)

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpResp, err := http.Post(f.webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	var resp feishuResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Code > 0 {
		logrus.Errorln(fmt.Sprintf("code %v msg %v", resp.Code, resp.Msg))
		return fmt.Errorf("code[%v],msg[%v]", resp.Code, resp.Msg)
	}
	return nil
}

func (f *FeishuRot) AsyncAlertText(keyWord string, title, msg string, nominees ...string) {
	go f.AlertText(keyWord, title, msg, nominees...)
}
