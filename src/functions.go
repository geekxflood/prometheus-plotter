package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func parseTime(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func getData(url string, query string, start int64, end int64) ([]byte, error) {
	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", fmt.Sprintf("%d", start))
	q.Add("end", fmt.Sprintf("%d", end))
	q.Add("step", "1m")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func extractData(data []byte) (plotter.XYs, error) {
	var pd PrometheusData
	err := json.Unmarshal(data, &pd)
	if err != nil {
		return nil, err
	}

	if pd.Status != "success" {
		return nil, fmt.Errorf("prometheus returned status '%s'", pd.Status)
	}

	if len(pd.Data.Result) == 0 {
		return nil, fmt.Errorf("no data returned by prometheus")
	}

	result := pd.Data.Result[0]
	values := result.Values

	pts := make(plotter.XYs, len(values))
	for i, v := range values {
		t, ok := v[0].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid timestamp '%v'", v[0])
		}
		y, ok := v[1].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid value '%v'", v[1])
		}
		pts[i].X = t
		pts[i].Y = y
	}

	return pts, nil
}

func createPlot(data plotter.XYs, title string, xlabel string, ylabel string, filename string) error {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	s, err := plotter.NewLine(data)
	if err != nil {
		return err
	}

	s.LineStyle.Width = vg.Points(1)
	s.LineStyle.Color = plotutil.Color(0)

	p.Add(s)

	err = p.Save(8*vg.Inch, 4*vg.Inch, filename)
	if err != nil {
		return err
	}

	return nil
}
