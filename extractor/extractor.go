package extractor

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mglnsk/teamscli/feed"

	"github.com/microcosm-cc/bluemonday"
)

type Output struct {
	Date   string   `json:"date,omitempty"`
	Author string   `json:"author,omitempty"`
	Text   string   `json:"text,omitempty"`
	Alt    string   `json:"alt,omitempty"`
	Images []string `json:"images,omitempty"`
	Type   string   `json:"type,omitempty"`
}

func LoadFeed(feedFile string) (feed.FeedJSON, error) {
	var info_json feed.FeedJSON
	channelb, err := os.ReadFile(feedFile)
	if err != nil {
		return info_json, err
	}
	err = json.Unmarshal(channelb, &info_json)
	return info_json, err
}

type MessageParams struct {
	Msg         feed.Message
	Time_filter string
	ShowCalls   bool
	ShowDate    bool
	ShowAuthor  bool
	ShowText    bool
	ShowAlt     bool
	ShowImages  bool
	ShowAll     bool
}

func ParseMessage(params MessageParams) (Output, error) {
	var outjson Output

	if params.Time_filter != "" {
		after_time, err := time.Parse("2006-01-02T15:04:05", params.Time_filter)
		if err != nil {
			return outjson, err
		}
		if params.Msg.ComposeTime == "" {
			return outjson, nil
		}
		message_time, err := time.Parse("2006-01-02T15:04:05Z", params.Msg.ComposeTime)
		if err != nil {
			return outjson, err
		}
		if message_time.Before(after_time) {
			return outjson, nil
		}
	}

	if params.Msg.MessageType == "Event/Call" && params.ShowCalls {
		outjson.Type = params.Msg.MessageType
		outjson.Date = params.Msg.ComposeTime
		b, err := json.Marshal(outjson)
		if err != nil {
			return outjson, err
		}
		fmt.Println(string(b))
		return outjson, err
	}

	p := bluemonday.StripTagsPolicy()
	alt_regexp := regexp.MustCompile(`alt=".+?"`)
	src_regexp := regexp.MustCompile(`src="https.*?"`)

	alt := alt_regexp.FindString(params.Msg.Content)
	if alt != "" {
		alt = alt[5 : len(alt)-1]
	}
	srcs := src_regexp.FindAllString(params.Msg.Content, -1)
	text := strings.TrimSpace(p.Sanitize(params.Msg.Content))

	if params.ShowAll {
		params.ShowDate = true
		params.ShowAuthor = true
		params.ShowText = true
		params.ShowAlt = true
		params.ShowImages = true
	}
	if params.ShowDate {
		outjson.Date = params.Msg.ComposeTime
	}
	if params.ShowAuthor {
		outjson.Author = params.Msg.ImDisplayName
	}
	if params.ShowText && text != "" {
		outjson.Text = text
	}
	if params.ShowAlt {
		outjson.Alt = alt
	}
	if params.ShowImages {
		for i, src := range srcs {
			s := strings.TrimPrefix(src, "src=\"")
			s = strings.TrimSuffix(s, "\"")
			srcs[i] = s
		}
		outjson.Images = srcs

	}
	return outjson, nil
}
