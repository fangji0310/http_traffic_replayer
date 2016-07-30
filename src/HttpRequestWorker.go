package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequestWorker struct {
	Id int
}

func (worker HttpRequestWorker) Proceed(url string, job HttpJob) {
	job.StartTime = time.Now().UnixNano()
	client := http.Client{}
	request, err := generateRequest(url, job)
	if err != nil {
		log.Print(err)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return
	}
	defer response.Body.Close()
	job.EndTime = time.Now().UnixNano()
	log.Printf("%s %s %d %d", job.Method, job.Path, response.StatusCode, (job.EndTime-job.StartTime)/1000)
}

func generateRequest(url string, job HttpJob) (*http.Request, error) {
	if strings.ToUpper(job.Method) == "GET" {
		return get(url, job)
	}
	return post(url, job)
}
func get(url string, job HttpJob) (request *http.Request, err error) {
	req, reqErr := http.NewRequest("GET", url+job.Path, nil)
	if reqErr != nil {
		log.Print(reqErr)
		err = reqErr
		return
	}
	params := generateParameters(job)
	req.URL.RawQuery = params.Encode()
	request = req
	err = nil
	return
}
func post(url string, job HttpJob) (request *http.Request, err error) {
	params := generateParameters(job)
	req, reqErr := http.NewRequest("POST", url+job.Path, strings.NewReader(params.Encode()))
	if reqErr != nil {
		log.Print(reqErr)
		err = reqErr
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request = req
	err = nil
	return
}
func generateParameters(job HttpJob) url.Values {
	values := url.Values{}
	for _, parameter := range job.Parameters {
		values.Add(parameter.Key, parameter.Value)
	}
	return values
}
