package main

type PrometheusData struct {
	Status string
	Data   PrometheusResult
}

type PrometheusResult struct {
	ResultType string
	Result     []PrometheusMetric
}

type PrometheusMetric struct {
	Metric map[string]string
	Values [][]interface{}
}
