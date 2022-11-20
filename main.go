package main

import (
	"context"
	"log"
	"time"

	gohttpclient "github.com/bozd4g/go-http-client"
)

//* DOC: https://prometheus.io/docs/prometheus/latest/querying/api/

// Prometheus rules data struct

type RulesResponseObject struct {
	Status string `json:"status"`
	Data   struct {
		Groups []RulesGroupObject `json:"groups"`
	} `json:"data"`
}

type RulesGroupObject struct {
	Name           string        `json:"name"`
	File           string        `json:"file"`
	Rules          []RulesObject `json:"rules"`
	Interval       int           `json:"interval"`
	Limit          int           `json:"limit"`
	EvaluationTime float64       `json:"evaluationTime"`
	LastEvaluation time.Time     `json:"lastEvaluation"`
}

type RulesObject struct {
	State          string            `json:"state"`
	Name           string            `json:"name"`
	Query          string            `json:"query"`
	Duration       int               `json:"duration"`
	Labels         LabelsObject      `json:"labels"`
	Annotations    AnnotationsObject `json:"annotations"`
	Alerts         []interface{}     `json:"alerts"`
	Health         string            `json:"health"`
	EvaluationTime float64           `json:"evaluationTime"`
	LastEvaluation time.Time         `json:"lastEvaluation"`
	Type           string            `json:"type"`
}

type LabelsObject struct {
	Severity string `json:"severity"`
}

type AnnotationsObject struct {
	Summary string `json:"summary"`
}

type RulesInProm struct {
	rules []struct {
		ruleName  string
		ruleQuery string
	}
}

func PromApiCaller(endpoint string) *gohttpclient.Response {
	ctx := context.Background()
	client := gohttpclient.New("http://localhost:9090")

	res, err := client.Get(ctx, endpoint)
	if err != nil {
		log.Fatalln(err)
	}
	return res
}

// plotAlerts plots the alerts from the Prometheus API
func GetRules() RulesInProm {
	log.Println("Getting current rules from Prometheus API")
	var rList RulesResponseObject
	var rInProm RulesInProm

	res := PromApiCaller("/api/v1/rules")

	// log.Println(res)

	err := res.Json(&rList)
	if err != nil {
		log.Fatalln(err)
	}

	for _, group := range rList.Data.Groups {
		for _, rule := range group.Rules {
			// log.Println(rule.Query)
			rInProm.rules = append(rInProm.rules, struct {
				ruleName  string
				ruleQuery string
			}{ruleName: rule.Name, ruleQuery: rule.Query})

		}
	}

	return rInProm
}

func main() {
	log.Println("Prometheus Graph Plotter")
	rProm := GetRules()
	log.Println(rProm.rules)
}
