package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// 環境変数の読み込み
	appToken, botToken := loadEnv()

	// Slack APIクライアントの初期化
	api, botUserID := initSlackAPI(appToken, botToken)

	// Socket Modeクライアントの初期化
	client := initSocketMode(api)

	// イベントの処理
	handleEvents(client, api, botUserID)
}

// 環境変数の読み込み
func loadEnv() (string, string) {
	fmt.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println(".env file loaded successfully.")

	appToken := os.Getenv("SLACK_APP_TOKEN") // xapp-で始まるトークン
	botToken := os.Getenv("SLACK_BOT_TOKEN") // xoxb-で始まるトークン

	if appToken == "" {
		log.Fatalf("SLACK_APP_TOKEN is not set in .env file")
	}
	if botToken == "" {
		log.Fatalf("SLACK_BOT_TOKEN is not set in .env file")
	}
	fmt.Println("Tokens retrieved successfully.")

	return appToken, botToken
}

// Slack APIクライアントの初期化
func initSlackAPI(appToken, botToken string) (*slack.Client, string) {
	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(true),
	)
	fmt.Println("Slack API client initialized.")

	// ボットのユーザーIDを取得
	authTest, err := api.AuthTest()
	if err != nil {
		log.Fatalf("Failed to get bot user ID: %v", err)
	}
	botUserID := authTest.UserID
	fmt.Printf("Bot User ID: %s\n", botUserID)

	return api, botUserID
}

// Socket Modeクライアントの初期化
func initSocketMode(api *slack.Client) *socketmode.Client {
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
	)
	fmt.Println("Socket Mode client initialized.")

	// Socket Modeクライアントを実行
	go client.Run()
	fmt.Println("Socket Mode client started.")

	return client
}

// イベントの処理
func handleEvents(client *socketmode.Client, api *slack.Client, botUserID string) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			handleEventsAPIEvent(evt, client, api, botUserID)
		default:
			// その他のイベントタイプを無視
			// fmt.Printf("Ignored event: %v\n", evt.Type)
		}
	}
}

// Events APIイベントの処理
func handleEventsAPIEvent(evt socketmode.Event, client *socketmode.Client, api *slack.Client, botUserID string) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Could not type cast the event to EventsAPIEvent: %v\n", evt)
		return
	}
	client.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		handleCallbackEvent(eventsAPIEvent.InnerEvent, api, botUserID)
	}
}

// コールバックイベントの処理
func handleCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent, api *slack.Client, botUserID string) {
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		handleAppMentionEvent(ev, api, botUserID)
	}
}

// メンションイベントの処理
func handleAppMentionEvent(ev *slackevents.AppMentionEvent, api *slack.Client, botUserID string) {
	// メンションされたメッセージテキストを取得
	text := ev.Text
	// メンション部分を除去
	text = strings.Replace(text, "<@"+botUserID+">", "", -1)
	text = strings.TrimSpace(text)

	// テキストを"*"に変換
	maskedText := strings.Repeat("*", len(text))

	// メンションしたユーザーにスレッドで返信
	_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(maskedText, false), slack.MsgOptionTS(ev.TimeStamp))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}
