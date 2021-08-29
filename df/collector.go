package diskcollector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DiskCollector interface {
	Collect() PrometheusDiskInfo
}

type CollectHandler struct {
	handle DiskCollector
}

type DiskInfo struct {
	Path       string
	Bsize      float64
	Bused      float64
	Bfree      float64
	Bavailable float64
}

type PrometheusDiskInfo struct {
	DiskInfo
	DescSize      *prometheus.Desc
	DescUsed      *prometheus.Desc
	DescFree      *prometheus.Desc
	DescAvailable *prometheus.Desc
}

func (h *CollectHandler) Register(c DiskCollector) {
	h.handle = c
}

func (h *CollectHandler) Describe(ch chan<- *prometheus.Desc) {
	i := h.handle.Collect()
	ch <- i.DescSize
	ch <- i.DescUsed
	ch <- i.DescFree
	ch <- i.DescAvailable
}

func (h *CollectHandler) Collect(ch chan<- prometheus.Metric) {
	i := h.handle.Collect()
	ch <- prometheus.MustNewConstMetric(i.DescSize, prometheus.GaugeValue, i.Bsize)
	ch <- prometheus.MustNewConstMetric(i.DescUsed, prometheus.GaugeValue, i.Bused)
	ch <- prometheus.MustNewConstMetric(i.DescFree, prometheus.GaugeValue, i.Bfree)
	ch <- prometheus.MustNewConstMetric(i.DescAvailable, prometheus.GaugeValue, i.Bavailable)
}
