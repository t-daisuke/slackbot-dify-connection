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
	// .envファイルを読み込む
	fmt.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println(".env file loaded successfully.")

	// .envファイルからトークンを取得
	appToken := os.Getenv("SLACK_APP_TOKEN") // xapp-で始まるトークン
	botToken := os.Getenv("SLACK_BOT_TOKEN") // xoxb-で始まるトークン

	if appToken == "" {
		log.Fatalf("SLACK_APP_TOKEN is not set in .env file")
	}
	if botToken == "" {
		log.Fatalf("SLACK_BOT_TOKEN is not set in .env file")
	}
	fmt.Println("Tokens retrieved successfully.")

	// Slack APIクライアントを作成
	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(true),
	)
	fmt.Println("Slack API client initialized.")

	// Socket Modeクライアントを作成
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
	)
	fmt.Println("Socket Mode client initialized.")

	// Socket Modeクライアントを実行
	go client.Run()
	fmt.Println("Socket Mode client started.")

	// ボットのユーザーIDを取得
	authTest, err := api.AuthTest()
	if err != nil {
		log.Fatalf("Failed to get bot user ID: %v", err)
	}
	botUserID := authTest.UserID
	fmt.Printf("Bot User ID: %s\n", botUserID)

	// イベントを処理するループ
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				fmt.Printf("Could not type cast the event to EventsAPIEvent: %v\n", evt)
				continue
			}
			client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
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
			}
		default:
			// その他のイベントタイプを無視
			// fmt.Printf("Ignored event: %v\n", evt.Type)
		}
	}
}
