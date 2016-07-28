package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)
const time_format = "2006/01/02 15:04:05"

var maxWorkers *int
var maxQueueSize *int
var targetUrl *string
var filePath *string
var timeFormat *string

func main() {
	var fp *os.File
	var err error
	maxWorkers   = flag.Int("max_workers", 5, "the number of workers")
	maxQueueSize = flag.Int("max_queue_size", 100, "the size of job queue")
	targetUrl    = flag.String("target_url", "http://127.0.0.1/target", "target url")
	filePath	= flag.String("replay_log_path", "replay_log.txt", "replay_log_path")
	timeFormat	= flag.String("time_format", time_format, "time_format")
	flag.Parse()
	log.Printf("max_workers %d max_queue_size %d target_url %s replay_log_path %s ", *maxWorkers, *maxQueueSize, *targetUrl, *filePath)
	fp, err = os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	jobChannel := make(chan HttpJob, *maxQueueSize)
	initWorker(jobChannel, *maxWorkers, *targetUrl)
	replayLog(jobChannel, fp)
	time.Sleep(10 * time.Second)
}
func initWorker(jobChannel chan HttpJob, maxWorkers int, targetUrl string) {
	for i := 0; i < maxWorkers; i++ {
		worker := HttpRequestWorker{i}
		go func(worker HttpRequestWorker) {
			for job := range jobChannel {
				worker.Proceed(targetUrl, job)
			}
		}(worker)
	}
}
func replayLog(jobChannel chan HttpJob, fp *os.File) {
	baseTimeSeconds := time.Now().Unix()
	var originalBaseTimeSeconds int64 = 0
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		requestTimeSeconds, job, err := parse([]byte(scanner.Text()))
		if err != nil {
			log.Fatal(err)
			return
		}
		if (originalBaseTimeSeconds == 0) {
			originalBaseTimeSeconds = requestTimeSeconds
		}
		waitTime := (requestTimeSeconds - originalBaseTimeSeconds) - (time.Now().Unix() - baseTimeSeconds)
		if waitTime > 0 {
			time.Sleep(time.Second * time.Duration(waitTime))
		}
		jobChannel <- job
	}
}

func parse(line []byte) (requestTimeSeconds int64, job HttpJob, err error) {
	unmarshalError := json.Unmarshal(line, &job)
	if unmarshalError != nil {
		log.Fatal(unmarshalError)
		err = unmarshalError
		return
	}
	job.EnqueueTime = time.Now().UnixNano()
	time, parseErr := time.Parse(*timeFormat, job.RequestTime)
	if parseErr != nil {
		log.Fatal(parseErr)
		err = parseErr
		return
	}
	requestTimeSeconds = time.Unix()
	return
}
