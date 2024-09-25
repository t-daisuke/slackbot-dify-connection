package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	// .envファイルを読み込む
	fmt.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println(".env file loaded successfully.")

	// .envファイルからSLACK_BOT_TOKENを取得
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatalf("SLACK_BOT_TOKEN is not set in .env file")
	}
	fmt.Println("SLACK_BOT_TOKEN retrieved successfully.")

	// Slack APIクライアントを作成
	api := slack.New(token)
	fmt.Println("Slack API client initialized.")

	// WebSocketを使ってリアルタイムでイベントを処理
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	fmt.Println("WebSocket connection started.")

	// イベントを処理するループ
	for msg := range rtm.IncomingEvents {
		fmt.Println("Event received: ", msg)

		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// メッセージが投稿された場合の処理
			fmt.Printf("Message: %v\n", ev)

			// 応答するメッセージを作成
			_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Hello, you said: "+ev.Text, false))
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			}
		default:
			// 他のイベントは無視
			fmt.Println("Unhandled event type")
		}
	}
}
