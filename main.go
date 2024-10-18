package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// グローバル変数の定義
var (
	// conversationID = ""
	apiKey     string
	difyAPIURL string
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
	apiKey = os.Getenv("DIFY_API_KEY")       // DifyのAPIキー
	difyAPIURL = os.Getenv("DIFY_API_URL")   // DifyのAPI URL

	if appToken == "" {
		log.Fatalf("SLACK_APP_TOKEN is not set in .env file")
	}
	if botToken == "" {
		log.Fatalf("SLACK_BOT_TOKEN is not set in .env file")
	}
	if apiKey == "" {
		log.Fatalf("API_KEY is not set in .env file")
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
	fmt.Printf("Received AppMentionEvent: %v\n", ev)
	// メンションされたメッセージテキストを取得
	text := ev.Text
	// メンション部分を除去
	text = strings.Replace(text, "<@"+botUserID+">", "", -1)
	text = strings.TrimSpace(text)

	// Dify APIにリクエストを送信
	answer, err := callDifyAPI(text, ev.User)
	if err != nil {
		// エラーが発生した場合、「アプリに問題が発生しました」と返信
		_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("アプリに問題が発生しました", false), slack.MsgOptionTS(ev.TimeStamp))
		if err != nil {
			fmt.Printf("Failed to send error message: %v\n", err)
		}
		return
	}

	// Dify APIからの回答を返信
	_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(answer, false), slack.MsgOptionTS(ev.TimeStamp))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}

// Dify APIへのリクエスト
func callDifyAPI(query string, userID string) (string, error) {
	// リクエストボディを作成
	requestBody := map[string]interface{}{
		"inputs":        map[string]interface{}{}, // 空のオブジェクト
		"query":         query,
		"response_mode": "blocking",
		"user":          userID,
	}

	// JSONにシリアライズ
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return "", err
	}

	// リクエストボディを表示
	fmt.Printf("Request body: %s\n", string(jsonData))

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", difyAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		return "", err
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// HTTPクライアントでリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return "", err
	}

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		// エラーメッセージを表示
		fmt.Printf("Received non-OK HTTP status: %s\nResponse body: %s\n", resp.Status, string(body))
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	// レスポンスをパース
	var responseData struct {
		Answer         string                 `json:"answer"`
		MessageID      string                 `json:"message_id"`
		ConversationID string                 `json:"conversation_id"`
		Mode           string                 `json:"mode"`
		Metadata       map[string]interface{} `json:"metadata"`
		CreatedAt      int64                  `json:"created_at"`
	}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return "", err
	}

	// 会話IDを更新（必要に応じて）
	// conversationID = responseData.ConversationID

	return responseData.Answer, nil
}
