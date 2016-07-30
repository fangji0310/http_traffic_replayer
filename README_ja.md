# http_traffic_replayer
GO言語の勉強のために書いてみました。
jmeterでApacheのGetリクエストをリプレイできるのですが、POSTでリクエストするようなAPIサーバを作成したあと、
リニューアル前に本番のアクセスログから想定リクエスト実行したい場合があります。
1行1リクエストでjsonファイルを用意すれば、このツールでログファイルのタイムスタンプを元に同じような負荷をかけてくれます。


###### usage
    http_traffic_replayer --max_workers n --max_queue_size n --target_url xxx --replay_log_path <file_path>

例. http_traffic_replayer --max_workers 10 --max_queue_size 10 --target_url http://127.0.0.1/ --replay_log_path ../json.txt

プロパティ

|プロパティ名|デフォルト|詳細|
|---|---|---|
|max_workers|5|リクエストを実行するWORKER数|
|max_queue_size|100|キューに挿入する最大サイズ|
|target_url|http://127.0.0.1/ |リプレイ時に付与するURLの接頭辞(http://xxx/)|
|replay_log_path|replay_log.txt|リプレイする内容が書かれたJSONファイルのファイルパス|
|time_format|2006/01/02 15:04:05|JSONに記載された時間のタイムフォーマット|

###### build
    go build -o http_traffic_replayer *.go
    
###### reference
http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/


