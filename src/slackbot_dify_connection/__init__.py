import os
from dotenv import load_dotenv
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler

# 環境変数の読み込み
load_dotenv()

# Slack Boltアプリの初期化
app = App(token=os.environ["SLACK_BOT_TOKEN"])

# メッセージイベントのリスナー
@app.message("hello")
def message_hello(message, say):
    say(f"こんにちは、<@{message['user']}>さん！")

# メンションイベントのリスナー
@app.event("app_mention")
def handle_mention(event, say):
    say(f"はい、<@{event['user']}>さん。何かお手伝いできますか？")

# アプリの起動
if __name__ == "__main__":
    handler = SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"])
    handler.start()