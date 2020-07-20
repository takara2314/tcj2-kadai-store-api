# 🎒 TCJ2 Kadai Store API
[Devoirs](https://github.com/approvers/devoirs) を使用して [MS Teams](https://www.microsoft.com/ja-jp/microsoft-365/microsoft-teams/group-chat-software/)(Microsoft 365 Education) で出されている課題の情報を取得し、APIを提供します。
この場をお借りして Devoirs の開発に関わった方々に感謝申し上げます。

## 👨‍💻 API利用方法
### 課題一覧を取得 — GET /get
#### Curl Example
```Bash
$ curl https://example.com/get \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```
#### Response Example
```JSON
{
  "acquisition": "2020-04-01T12:34:56.000Z",
  "homeworks": [
    {
      "subject": "情報工学",
      "omitted": "情報",
      "name": "第1回課題",
      "id": "1234",
      "due": "2020-04-01T12:34:56.000Z"
	}
  ]
}
```

## ✔ 前提 (提供側)
- Node.js 13+
- npm
- Go

## 🛠 環境構築 (提供側)
### 1. WebサーバーとGUIアプリの実行環境を構築します。
Devoirs(v0.3.0以前)はCLIアプリですが、認証時にGUIを用いるので、RDPやX転送などでGUIアプリを実行できる環境を構築してください。

### 2. Devoirs **v0.3.0** のソースコードをインストールします。
実行ファイルを実行するときに出るログを使用するので、**GUI版のDevoirs(v1.0.0以降)では動作しません。**

### 3. このレポジトリをGOPATHの中にクローンします。

### 4. GOPATHの中に以下のファイルを加えます。
- kadai-store-api.token
**APIで許可するトークン**を記述します。
改行区切りで複数のトークンを指定することができます。

#### 任意
- kadai-store-api_discord-alarm.token
- kadai-store-api_admin-discord-ID.id
それぞれ**DiscordBotのトークン**、**API管理者のDiscordID**を記述します。
複数のトークンやIDを入れることはできません。
これらのファイルを加えることによって、Devoirsの実行エラーが生じたときに、DiscordのDMで通知を受け取ることができます。

### 5. それぞれのファイルやフォルダを以下のディレクトリ管理下に配置します。
```
./
├─ deviors/ ................................. Deviors v0.3.0
└─ go/ ...................................... $GOPATH
   ├─ kadai-store-api.token ................. APIで許可するトークン
   ├─ kadai-store-api_discord-alarm.token ... DiscordBOTのトークン
   ├─ kadai-store-api_admin-discord-ID.id ... API管理者のDiscordID
   └─ kadai-store-api/ ...................... このレポジトリ
```

### 6. deviors/src/main.ts の20~26行目の次の構文を変更します。
```TypeScript:main.ts
for (const c of await client.getClasses()) {
	console.log(`-`, c.name);

	for (const a of await client.getAssignments(c.id)) {
		console.log('\t', a['isCompleted'] ? '✔' : '❗', a.displayName);
	}
}
```
を
```TypeScript:main.ts
for (const c of await client.getClasses()) {
	console.log(`-`, c.name);

	for (const a of await client.getAssignments(c.id)) {
		console.log("・" + a.displayName + "\t" + a.id + "\t" + a.dueDateTime);
	}
}
```
に書き換えます。
これをすることによって、課題のIDや提出期限も取得できるようになります。

### 7. go/tcj2-kadai-store-api/ で**go build**を実行し、実行ファイルを生成します。
実行できない場合は、[Releases](https://github.com/takara2314/tcj2-kadai-store-api/releases)から実行ファイルをダウンロードしてください。

### 8. 生成された実行ファイルを実行します。
8080ポートでAPIが提供されます。