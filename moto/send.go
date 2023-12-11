package moto

import (
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func Send(user string, fallback string, title string, texts []string, messageLink string, image string, imageLink string) string {
	decoded := new(string)

	fields := []string{}
	if image != "" {
		if imageLink != "" {
			fields = append(fields, `{ "short": false, "value": "[![embedded image](`+image+`)](`+imageLink+`)" }`)
		} else {
			fields = append(fields, `{ "short": false, "value": "![embedded image](`+image+`)" }`)
		}
	}

	for _, text := range texts {
		fields = append(fields, `{ "short": false, "value": "`+text+`" }`)
	}

	bot := "huoji"
	emoji := "seedling"
	channel := ""
	hook := ""
	switch user {
	case "hang":
		channel = "hm"
		hook = MattermostWebhookHang
	case "fang":
		channel = "fm"
		hook = ""
	}
	link := "https://mattermost.qifnle.com/hooks/" + hook
	if fallback == "" {
		fallback = title
	}
	form := `{ "username": "` + bot + `", "icon_emoji": "` + emoji + `", "channel": "` + channel + `", "attachments": [ { "color": "#202C51", "fallback": "` + fallback + `", "title": "` + title + `", "title_link": "` + messageLink + `", "fields": [` + strings.Join(fields, ",") + `] } ] }`

	client := resty.New()

	client.SetTimeout(5 * time.Second)
	res, err := client.R().
		SetHeaders(map[string]string{
			"content-type": "application/json",
		}).
		SetBody(form).
		SetResult(decoded).
		Post(link)
	if err != nil {
		panic(err)
	}

	return res.String()
}
