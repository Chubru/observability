package metrics

import (
	metrics "github.com/VictoriaMetrics/metrics"
	"net/http"
	"strings"
)

type VMetricsParameters struct {
	Namespace string
	Subsystem string
}

func (v *VMetrics) buildLabel(label string) string {
	if v.baseLabels == "" && label == "" {
		return ""
	}

	if v.baseLabels == "" {
		return "{" + label + "}"
	}
	if label == "" {
		return "{" + v.baseLabels + "}"
	}
	return "{" + v.baseLabels + "," + label + "}"
}

func (v VMetricsParameters) buildMetricName() string {
	var metricName strings.Builder
	if v.Namespace != "" {
		metricName.WriteString(v.Namespace + "_")
	}
	if v.Subsystem != "" {
		metricName.WriteString(v.Subsystem + "_")
	}

	return metricName.String()
}

type VMetrics struct {
	prefix     string
	baseLabels string
}

func NewVMetrics(opts VMetricsParameters, baseLabels string) *VMetrics {
	return &VMetrics{opts.buildMetricName(), baseLabels}
}

func (v *VMetrics) NewCounter(name string, label string) *vmetrics.Counter {
	return metrics.NewCounter(v.prefix + name + v.buildLabel(label))
}

func (v *VMetrics) GetOrCreateCounter(name string, label string) *vmetrics.Counter {
	return metrics.GetOrCreateCounter(v.prefix + name + v.buildLabel(label))
}

func (v *VMetrics) NewGauge(name string, label string, f func() float64) *vmetrics.Gauge {
	return metrics.NewGauge(v.prefix+name+v.buildLabel(label), f)
}

func (v *VMetrics) GetOrCreateGauge(name string, label string, f func() float64) *vmetrics.Gauge {
	return metrics.GetOrCreateGauge(v.prefix+name+v.buildLabel(label), f)
}

func (v *VMetrics) NewHistogram(name string, label string) *vmetrics.Histogram {
	return metrics.NewHistogram(v.prefix + name + v.buildLabel(label))
}

func (v *VMetrics) GetOrCreateHistogram(name string, label string) *vmetrics.Histogram {
	return metrics.GetOrCreateHistogram(v.prefix + name + v.buildLabel(label))
}

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	return mux
}
