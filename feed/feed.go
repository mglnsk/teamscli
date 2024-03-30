package feed

import (
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	MessageType   string `json:"messageType"`
	Content       string `json:"content"`
	ImDisplayName string `json:"imDisplayName"`
	Type          string `json:"type"`
	ComposeTime   string `json:"composeTime"`
}

type Reply struct {
	Messages []Message `json:"messages"`
}

type FeedJSON struct {
	ReplyChains     []Reply `json:"replyChains"`
	HasMore         bool    `json:"hasMore"`
	ChannelTenantId string  `json:"channelTenantId"`
}

func GetChannelInfo(channelID, token string) ([]byte, error) {
	channel_format := "https://teams.microsoft.com/api/csa/emea/api/v2/teams/%s@thread.tacv2/channels/%s@thread.tacv2?filterSystemMessage=true&pageSize=50"
	channel := fmt.Sprintf(channel_format, channelID, channelID)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", channel, nil)

	req.Header.Set("content-length", "0")
	req.Header.Set("authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyB, err := io.ReadAll(resp.Body)
	return bodyB, err
}

func DownloadImage(imageURL, skype_token string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", imageURL, nil)

	req.Header.Set("authorization", "skype_token "+skype_token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyB, err := io.ReadAll(resp.Body)
	return bodyB, err
}
