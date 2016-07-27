# http_traffic_replayer
http_traffic_replayer make you able to replay http_traffic

###### usage
    http_traffic_replayer --max_workers n --max_queue_size n --target_url xxx --replay_log_path <file_path>

ex. http_traffic_replayer --max_workers 10 --max_queue_size 10 --target_url http://127.0.0.1/ --replay_log_path ../json.txt

properties

|name|default|detail|
|---|---|---|
|name|default|detail|
|max_workers|5|the number of workers|
|max_queue_size|100|the size of job queue|
|target_url|http://127.0.0.1/ |target_url|
|replay_log_path|replay_log.txt|filepath to replay|

###### build
    go build -o http_traffic_replayer *.go


