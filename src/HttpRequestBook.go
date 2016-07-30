package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
)

type HttpRequestBook struct {
	JobChannel  chan HttpJob
	Timeformat  string
	FilePointer *os.File
}

func (book HttpRequestBook) replay() {
	baseTime := time.Now().Unix()
	var originalFirstRequestTime int64 = 0
	scanner := bufio.NewScanner(book.FilePointer)
	for scanner.Scan() {
		requestTime, job, err := book.parse([]byte(scanner.Text()))
		if err != nil {
			log.Fatal(err)
			return
		}
		if originalFirstRequestTime == 0 {
			originalFirstRequestTime = requestTime
		}
		waitTime := (requestTime - originalFirstRequestTime) - (time.Now().Unix() - baseTime)
		if waitTime > 0 {
			time.Sleep(time.Second * time.Duration(waitTime))
		}
		job.EnqueueTime = time.Now().UnixNano()
		book.JobChannel <- job
	}
}

func (book HttpRequestBook) parse(line []byte) (requestTime int64, job HttpJob, err error) {
	unmarshalError := json.Unmarshal(line, &job)
	if unmarshalError != nil {
		log.Fatal(unmarshalError)
		err = unmarshalError
		return
	}
	time, parseErr := time.Parse(book.Timeformat, job.RequestTime)
	if parseErr != nil {
		log.Fatal(parseErr)
		err = parseErr
		return
	}
	requestTime = time.Unix()
	return
}
