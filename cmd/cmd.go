package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mglnsk/teamscli/auth"
	"github.com/mglnsk/teamscli/extractor"
	"github.com/mglnsk/teamscli/feed"

	"github.com/alecthomas/kong"
)

type UpdateCmd struct {
	TenantID    string `help:"ID of your Teams provider e.g aaaaaaaa-d714-4a1f-8101-eeeeeeeeeeee" required:""`
	ID          string `help:"Channel id from teams e.g. 11:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"`
	RefreshFile string `help:"A file used for caching refresh tokens as they are issued" required:""`
	OutFile     string `help:"Output filename for channel json" required:""`
}

type RefreshCmd struct {
	TenantID    string `help:"ID of your Teams provider e.g aaaaaaaa-d714-4a1f-8101-eeeeeeeeeeee" required:""`
	RefreshFile string `help:"A file used for caching refresh tokens as they are issued" required:""`
}

type ExtractCmd struct {
	ShowText   bool   `help:"Show text messages"`
	ShowAlt    bool   `help:"Show alt tags"`
	ShowImages bool   `help:"Show links to images"`
	ShowDate   bool   `help:"Show timestamps"`
	ShowAuthor bool   `help:"Show senders"`
	ShowCalls  bool   `help:"Show calls"`
	ShowAll    bool   `help:"Show everything"`
	InFile     string `help:"Channel info file to load from required:" required:""`
	After      string `help:"Only include entries which happend after a certain time"`
}

type DownloadCmd struct {
	TenantID      string `help:"ID of your Teams provider e.g aaaaaaaa-d714-4a1f-8101-eeeeeeeeeeee" required:""`
	RefreshFile   string `help:"A file used for caching refresh tokens as they are issued" required:""`
	DownloadImage string `help:"Link for an image to download"`
	OutFile       string `help:"Output filename for channel json" required:""`
}

func (dcmd *DownloadCmd) Run(ctx *kong.Context) error {
	if dcmd.DownloadImage != "" {
		refb, err := os.ReadFile(dcmd.RefreshFile)
		if err != nil {
			return err
		}
		refresh := string(refb)
		skype_token := auth.GetSkypeToken(dcmd.TenantID, refresh)
		picture_b, err := feed.DownloadImage(dcmd.DownloadImage, skype_token)
		if err != nil {
			return err
		}
		os.WriteFile(dcmd.OutFile, picture_b, 0644)
	}
	return nil
}

func (ucmd *UpdateCmd) Run(ctx *kong.Context) error {
	refb, err := os.ReadFile(ucmd.RefreshFile)
	if err != nil {
		return err
	}
	refresh := string(refb)

	bearer, new_refresh := auth.GetTeamsToken(ucmd.TenantID, refresh)
	err = os.WriteFile(ucmd.RefreshFile, []byte(new_refresh), 0644)
	if err != nil {
		return err
	}

	//? Channel info
	info, err := feed.GetChannelInfo(ucmd.ID, bearer)
	if err != nil {
		return err
	}
	err = os.WriteFile(ucmd.OutFile, info, 0644)
	return err
}

func (rcmd *RefreshCmd) Run(ctx *kong.Context) error {
	refb, err := os.ReadFile(rcmd.RefreshFile)
	if err != nil {
		return err
	}
	refresh := string(refb)
	_, new_refresh := auth.GetTeamsToken(rcmd.TenantID, refresh)
	err = os.WriteFile(rcmd.RefreshFile, []byte(new_refresh), 0644)
	return err
}

func (ext *ExtractCmd) Run(ctx *kong.Context) error {
	info_json, err := extractor.LoadFeed(ext.InFile)
	if err != nil {
		return err
	}

	params := extractor.MessageParams{
		Time_filter: ext.After,
		ShowCalls:   ext.ShowCalls,
		ShowDate:    ext.ShowDate,
		ShowAuthor:  ext.ShowAuthor,
		ShowText:    ext.ShowText,
		ShowAlt:     ext.ShowAlt,
		ShowImages:  ext.ShowImages,
		ShowAll:     ext.ShowAll,
	}

	for _, r := range info_json.ReplyChains {
		for _, m := range r.Messages {
			params.Msg = m
			outjson, err := extractor.ParseMessage(params)
			if err != nil {
				return err
			}
			b, err := json.Marshal(outjson)
			if err != nil {
				return err
			}
			if len(b) > 2 { // {} has the length of 2
				fmt.Println(string(b))
			}
		}
	}
	return nil
}
