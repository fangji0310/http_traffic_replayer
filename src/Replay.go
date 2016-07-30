package main

import (
	"flag"
	"log"
	"os"
	"time"
)

const time_format = "2006/01/02 15:04:05"

func main() {
	var fp *os.File
	var err error
	var (
		maxWorkers = flag.Int("max_workers", 5, "the number of workers")
		maxQueueSize = flag.Int("max_queue_size", 100, "the size of job queue")
		targetUrl = flag.String("target_url", "http://127.0.0.1/target", "target url")
		filePath = flag.String("replay_log_path", "replay_log.txt", "replay_log_path")
		timeFormat = flag.String("time_format", time_format, "time_format")
	)
	flag.Parse()
	log.Printf("max_workers %d max_queue_size %d target_url %s replay_log_path %s ", *maxWorkers, *maxQueueSize, *targetUrl, *filePath)
	fp, err = os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	jobChannel := make(chan HttpJob, *maxQueueSize)
	initWorker(jobChannel, *maxWorkers, *targetUrl)
	replayLog(jobChannel, *timeFormat, fp)
	time.Sleep(5 * time.Second)
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
func replayLog(jobChannel chan HttpJob, timeFormat string, fp *os.File) {
	httpRequestBook := HttpRequestBook{JobChannel:jobChannel, Timeformat:timeFormat, FilePointer:fp}
	httpRequestBook.replay()
}
