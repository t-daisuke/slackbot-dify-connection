# 断念

# slackbot-dify-connection
difyのアプリをslackbotとして使いたいと思って作りました。

# Needs
- goを入れてください
- slackの管理者に問い合わせは必要です。


# How to connect Slack

1. Slack APIの管理ページ(https://api.slack.com/apps?new_app=1)にアクセスし、新しいアプリを作成します。
2. Createを選ぶ
3. Socket ModeをONにする。名前を入れないとtoken作れない。
4. Basic InformationでApp-Level Tokenを作る。これをSLACK_APP_TOKENに設定する。
5. OAuth & Permissionsで必要な権限をつけてBOT_USER_OAUTH_TOKENを取得する。これをSLACK_BOT_TOKENに設定する。
6. それぞれ.env_exampleをコピーして.envを作成して、tokenを設定する。
7. Event SubscriptionでSubscribe to bot eventsを選択する。

# Notion
2024/09/26時点ではSocket Modeを推奨していたので、それに対応

