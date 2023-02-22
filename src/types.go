// PrometheusData is the top level struct of the prometheus API response
// It contains the Status of the response, and the Data, which is an array of PrometheusResult
// The Status should be set to Success if the query completed successfully
// The Data should contain the results of the query

package main

// PrometheusData is the top level struct of the prometheus API response
type PrometheusData struct {
	Status string
	Data   PrometheusResult
}

// PrometheusResult contains the result type and the result
type PrometheusResult struct {
	ResultType string
	Result     []PrometheusMetric
}

// PrometheusMetric contains the metric, and the values
type PrometheusMetric struct {
	Metric map[string]string
	Values [][]interface{}
}
