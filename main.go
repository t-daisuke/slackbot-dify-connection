package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	// 環境変数からSlackトークンを取得
	token := os.Getenv("SLACK_BOT_TOKEN")
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
