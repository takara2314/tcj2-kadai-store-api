# 🎒 TCJ2 Kadai Store API
[Devoirs](https://github.com/approvers/devoirs) を使用して [Microsoft Teams](https://www.microsoft.com/ja-jp/microsoft-365/microsoft-teams/group-chat-software/)(Microsoft 365 Education) で出されている課題の情報を取得し、APIを提供します。
この場をお借りして Devoirs の開発に関わった方々に感謝申し上げます。

## 👨‍💻 API利用方法
### 課題一覧を取得 — GET /get
#### Parameters
| Parameter |                                                           |
|:----------|:---------------------------------------------------------:|
| due       | 提出期限が過ぎているものを除く場合は`future`を入れてください。 |
| timezone  | `Asia/Tokyo`のみサポートしています。それ以外はUTCとなります。 |
#### Curl Example
```Bash
$ curl https://example.com/get?due=future&timezone=Asia/Tokyo \
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
`acquisition`はDevoirsから取得した時刻です。

## ✔ 前提 (提供側)
- Node.js 13+
- npm
- Go

## 🛠 環境構築 (提供側)
### 1. WebサーバーとGUIアプリの実行環境を構築します。
デフォルトでApacheやNginxなどを用いて**FastCGIとしてサーバーを開く設定になっています。** 使用しない場合は、`config.yaml`の`fcgi-server`を`false`にしてください。
Devoirs(v0.3.0以前)はCLIアプリですが、アカウントの認証時にGUIを用いるので、RDPやX転送などの**GUIアプリを実行できる環境を構築してください。**

### 2. Devoirs v0.3.0 のソースコードをインストールします。
実行ファイルを実行するときに出るログを使用するので、**GUI版のDevoirs(v1.0.0以降)では動作しません。**

### 3. このレポジトリをGOPATHの中にインストールします。
[Releases](https://github.com/takara2314/tcj2-kadai-store-api/releases)のから最新版のソースコードをダウンロードして展開してください。
また、GOPATH以外では正しく動作しない場合があります。

### 4. config.yamlで基本的な設定を行います。
| Keys                   |                                                             |
|:-----------------------|:-----------------------------------------------------------:|
| update-times           | Devoirsを実行して情報を更新する時刻(分)                        |
| server-port            | API提供サーバーのポート番号                                   |
| fcgi-server            | FastCGIとして動かすかどうか                                  |
| get-limit              | 10分間にこのAPIにGETできる回数 (無制限なら`-1`を入れてください) |
| subjects/teams         | その教科のTeamsのチーム名                                    |
| subjects/syllabus      | シラバスでの教科名                                           |
| subjects/omitted       | 省略された教科名                                             |
| discord/alarm          | Discordでエラーを通知するかどうか                             |
| discord/admin-id       | エラー通知するユーザーのID (`alarm`を`true`にした方のみ)           |
| discord/message-format | メッセージフォーマット (`alarm`を`true`にした方のみ)               |
| discord/command-prefix | ボットを呼び出すコマンドの接頭辞 (`alarm`を`true`にした方のみ)      |

### 5. token.yamlでトークンについての設定を行います。
| Keys           |                                                                              |
|:---------------|:----------------------------------------------------------------------------:|
| allowed-tokens | 当APIにアクセスを許可するトークン (ここで定めたものを利用者に教えてください)       |
| discord-token  | Discordボットのトークン (`config.yaml`の`discord`の`alarm`を`true`にした方のみ) |

### 6. それぞれのファイルやフォルダを以下のディレクトリ管理下に配置します。
```
./
├─ devoirs/ ................................. devoirs v0.3.0
│  ├─ src ................................... devoirsのソースコード
│
└─ go/ ...................................... $GOPATH
   └─ kadai-store-api/ ...................... このレポジトリ
      ├─ main.go ............................ 主にWebアプリの処理
      ├─ config.yaml ........................ 基本的な設定
      ├─ token.yaml ......................... トークンの設定
```
``├─`` で終わってるものは、まだ続きがあることを示します。

### 7. devoirs/src/models/assignment.ts の内容を以下の内容に変更します。
```TypeScript:main.ts
import { ClassId } from './class';

export type AssignmentId = string;
export type DateTime = string;

export interface Assignment {
  id: AssignmentId;
  classId: ClassId;
  dueDateTime: DateTime;
  displayName: string;
  isCompleted: boolean;
}

export function compare(a: Assignment, b: Assignment): number {
  return a.dueDateTime.localeCompare(b.dueDateTime);
}
```
これをすることによって、課題の提出期限のデータを扱うことができるようになります。

### 8. devoirs/src/main.ts の20~26行目の次の構文を変更します。
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

### 9. GUI環境で devoirs/ でDevoirsを実行します。
```Bash
$ npm install
$ npm start
```
Devoirsを初めて起動すると、学校で提供されているMicrosoftアカウントのログインを要求するためにGUIウィンドウが表示されます。
ログインした後は、**自動認証が行われる1時間後まで**GUIが起動することは基本的にありません。

### 10. go/tcj2-kadai-store-api/ でgo buildを実行し、実行ファイルを生成します。
```Bash
$ go build
```
**tcj2-kadai-store-api**という名前で生成されます。Windowsだと拡張子がexeになります。

### 11. 生成された実行ファイルを実行します。
8080ポートでAPIが提供されます。

## ⌨️ Discordコマンド (オプションでDiscord関連の設定を行った方のみ)
DMやグループ、サーバーのどこでも実行することができます。
### ::kadai-store ping
サーバーが落ちていないかを確かめるコマンドです。
API管理者ユーザー以外でも実行できます。
### ::kadai-store version
`TCJ2 Kadai Store API`のバージョンを確認するコマンドです。
API管理者ユーザー以外でも実行できます。
### ::kadai-store stop
サーバーを強制終了させるコマンドです。
API管理者ユーザーのみ実行可能です。