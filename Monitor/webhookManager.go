package Monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func webhookPageCaptchaStuck(Proxy string, wbKey string) {

	var title = "Task stuck in captcha solving"
	var description = "proxy is stuck in captcha solving" + Proxy
	payloadData := struct {
		Content   interface{} `json:"content"`
		Embeds    []Embed     `json:"embeds"`
		Username  string      `json:"username"`
		AvatarURL string      `json:"avatar_url"`
	}{
		Content: nil,
		Embeds: []Embed{
			{
				Title: title,
				Color: 5814783,
				Fields: []Field{
					{
						Name:  "Task:",
						Value: description,
					},
				},
			},
		},
		Username:  "Monitor",
		AvatarURL: "https://i.imgur.com/gTtPuMp.png",
	}
	WebhookSend(payloadData, wbKey)
}

func WebhookTicketsLive(pageURL string, wbKey string) {
	log.Println("Tickets found, sending webhook")
	var title = "THERE ARE TICKETS AVAILABLE FOR 2"
	var description = "Tickets live here : " + pageURL

	payloadData := struct {
		Content   interface{} `json:"content"`
		Embeds    []Embed     `json:"embeds"`
		Username  string      `json:"username"`
		AvatarURL string      `json:"avatar_url"`
	}{
		Content: nil,
		Embeds: []Embed{
			{
				Title: title,
				Color: 5814783,
				Fields: []Field{
					{
						Name:  "Task",
						Value: description,
					},
				},
			},
		},
		Username:  "Monitor",
		AvatarURL: "https://i.imgur.com/gTtPuMp.png",
	}

	WebhookSend(payloadData, wbKey)
}

func WebhookQueueUp(pageURL string, wbKey string) {
	var title = "Queue is up"
	var description = "Queue is up at " + pageURL
	payloadData := struct {
		Content   interface{} `json:"content"`
		Embeds    []Embed     `json:"embeds"`
		Username  string      `json:"username"`
		AvatarURL string      `json:"avatar_url"`
	}{
		Content: nil,
		Embeds: []Embed{
			{
				Title: title,
				Color: 5814783,
				Fields: []Field{
					{
						Name:  "LINK BELOW",
						Value: description,
					},
				},
			},
		},
		Username:  "Monitor",
		AvatarURL: "https://i.imgur.com/gTtPuMp.png",
	}
	WebhookSend(payloadData, wbKey)
}

func WebhookSend(pld interface{}, wbKey string) {
	payload, err := json.Marshal(pld)
	if err != nil {
		log.Fatal("Encoding json failed")
	}
	req, err := http.NewRequest("POST", wbKey, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("couldn't create webhook")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Coudln't send request")
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Println("Couldn't close body")
	}
}
