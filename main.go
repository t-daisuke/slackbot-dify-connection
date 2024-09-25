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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// .envファイルからSLACK_BOT_TOKENを取得
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatalf("SLACK_BOT_TOKEN is not set in .env file")
	}

	api := slack.New(token)

	// WebSocketを使ってリアルタイムでイベントを処理
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// イベントを処理するループ
	for msg := range rtm.IncomingEvents {
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
		}
	}
}
