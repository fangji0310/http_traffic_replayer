package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	var fp *os.File
	var err error
	var (
		maxWorkers   = flag.Int("max_workers", 5, "the number of workers")
		maxQueueSize = flag.Int("max_queue_size", 100, "the size of job queue")
		targetUrl    = flag.String("target_url", "http://127.0.0.1/target", "target url")
		filePath	= flag.String("replay_log_path", "replay_log.txt", "replay_log_path")
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
	requestNo := 1
	baseTimeSeconds := time.Now().Unix()
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	originalBaseTimeSeconds, firstJob, err := parse([]byte(scanner.Text()))
	if err != nil {
		log.Fatal(err)
		return
	}
	firstJob.Line = requestNo
	firstJob.EnqueueTime = time.Now().UnixNano()
	jobChannel <- firstJob
	for scanner.Scan() {
		requestNo++
		requestTimeSeconds, job, err := parse([]byte(scanner.Text()))
		if err != nil {
			log.Fatal(err)
			return
		}
		waitTime := (requestTimeSeconds - originalBaseTimeSeconds) - (time.Now().Unix() - baseTimeSeconds)
		if waitTime > 0 {
			time.Sleep(time.Second * time.Duration(waitTime))
		}
		job.Line = requestNo
		job.EnqueueTime = time.Now().UnixNano()
		jobChannel <- job
	}
}
func parse(line []byte) (requestTimeSeconds int64, job HttpJob, err error) {
	unmarShalErr := json.Unmarshal(line, &job)
	if unmarShalErr != nil {
		log.Fatal(unmarShalErr)
		err = unmarShalErr
		return
	}
	time, parseErr := time.Parse("2006/01/02 15:04:05", job.RequestTime)
	if parseErr != nil {
		log.Fatal(parseErr)
		err = parseErr
		return
	}
	requestTimeSeconds = time.Unix()
	return
}
