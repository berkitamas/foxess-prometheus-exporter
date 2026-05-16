package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/berkitamas/foxess-prometheus-exporter/foxess"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FoxESSCollector struct {
	apiKey string
	gauge  *prometheus.Desc
}

func NewFoxESSCollector(apiKey string) *FoxESSCollector {
	return &FoxESSCollector{
		apiKey: apiKey,
		gauge: prometheus.NewDesc(
			"foxess_realtime_data",
			"Real-time data from FoxESS",
			[]string{"deviceSN", "variable"},
			nil,
		),
	}
}

func (c *FoxESSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.gauge
}

func (c *FoxESSCollector) Collect(ch chan<- prometheus.Metric) {
	devices, err := foxess.GetDevices(c.apiKey)
	if err != nil {
		log.Printf("Error getting devices: %v", err)
		return
	}

	var snList []string
	for _, d := range devices {
		snList = append(snList, d.DeviceSN)
	}

	data, err := foxess.GetRealTimeData(c.apiKey, snList, nil)
	if err != nil {
		log.Printf("Error getting real-time data: %v", err)
		return
	}

	for _, d := range data {
		for _, v := range d.Datas {
			var val float64
			switch vv := v.Value.(type) {
			case float64:
				val = vv
			case string:
				val, _ = strconv.ParseFloat(vv, 64)
			}
			ch <- prometheus.MustNewConstMetric(
				c.gauge,
				prometheus.GaugeValue,
				val,
				d.DeviceSN,
				v.Variable,
			)
		}
	}
}

func main() {
	apiKey := os.Getenv("FOXESS_API_KEY")
	if apiKey == "" {
		log.Fatal("FOXESS_API_KEY environment variable is required")
	}

	collector := NewFoxESSCollector(apiKey)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting exporter on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
