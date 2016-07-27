package main

import "fmt"

type parameters struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type HttpJob struct {
	RequestTime string `json:"request_time"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Parameters  []parameters
	EnqueueTime int64
	StartTime   int64
	EndTime     int64
	Line        int
}

func (job HttpJob) String() string {
	param := ""
	for _, p := range job.Parameters {
		param += "  " + p.String()
	}
	str := fmt.Sprintf("%s %s %s %s\n", job.RequestTime, job.Method, job.Path, param)
	return str
}
func (p parameters) String() string {
	return fmt.Sprintf("%s:%s", p.Key, p.Value)
}
