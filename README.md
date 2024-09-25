# 断念

# slackbot-dify-connection
difyのアプリをslackbotとして使いたいと思って作りました。

# Needs
- pythonでします。
- slackの管理者に問い合わせは必要です。

# Ref
https://zenn.dev/replica/articles/34b7fa8d667ec6


# How to connect Slack

1. Slack APIの管理ページ(https://api.slack.com/apps?new_app=1)にアクセスし、新しいアプリを作成します。
2. Createを選んでbotを選びます。

4. Slackアプリの設定

Goのコードに移る前に、Slack側でボットの設定を行います。

	1.	(Slack APIの管理ページ)[https://api.slack.com/apps?new_app=1]にアクセスし、新しいアプリを作成します。
	2.	「Bot User」を作成し、必要な権限を設定します。たとえば、chat:write権限を追加します。
	3.	アプリをワークスペースにインストールし、OAuthトークンを取得します。
	4.	取得したOAuthトークンを環境変数に設定します。

```
export SLACK_BOT_TOKEN=your_slack_token_here
```

5. go run main.go
これで、Goで作成したSlackボットが起動し、メッセージを受け取った際に「Hello, you said:」と応答するようになります。

ボットの追加機能

GoのSlackボットには、さらに以下のような機能を追加できます。

	•	特定のキーワードに反応する
	•	スラッシュコマンドを処理する
	•	外部APIを呼び出してデータを取得する
	•	定期的に通知を送る

Goの並行処理（goroutines）やチャンネルを活用して、効率的で強力なボットを作ることができます。