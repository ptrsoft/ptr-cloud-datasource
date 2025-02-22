package metricutil

import (
	"errors"
	"fmt"
	//"github.com/prometheus/client_golang/prometheus"
	"strings"
)

// SanitizeLabelName removes all invalid chars from the label name.
// If the label name is empty or contains only invalid chars, it
// will return an error.
func SanitizeLabelName(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("label name cannot be empty")
	}

	out := strings.Builder{}
	for i, b := range name {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || (b >= '0' && b <= '9' && i > 0) {
			out.WriteRune(b)
		} else if b == ' ' {
			out.WriteRune('_')
		}
	}

	if out.Len() == 0 {
		return "", fmt.Errorf("label name only contains invalid chars: %q", name)
	}

	return out.String(), nil
}

// NewCounterStartingAtZero initializes a new Prometheus counter with an initial
// observation of zero. Used for to guarantee the existence of the specific metric.
//func NewCounterStartingAtZero(opts prometheus.CounterOpts) prometheus.Counter {
//	counter := prometheus.NewCounter(opts)
//	counter.Add(0)
//	return counter
//}

// NewCounterVecStartingAtZero initializes a new Prometheus counter with an initial
// observation of zero for every possible value of each label. Used for the sake of
// consistency among all the possible labels and values.
//func NewCounterVecStartingAtZero(opts prometheus.CounterOpts, labels []string, labelValues map[string][]string) *prometheus.CounterVec {
//	counter := prometheus.NewCounterVec(opts, labels)
//
//	for _, ls := range buildLabelSets(labels, labelValues) {
//		counter.With(ls).Add(0)
//	}
//
//	return counter
//}

//func buildLabelSets(labels []string, labelValues map[string][]string) []prometheus.Labels {
//	var labelSets []prometheus.Labels
//
//	var n func(i int, ls prometheus.Labels)
//	n = func(i int, ls prometheus.Labels) {
//		if i == len(labels) {
//			labelSets = append(labelSets, ls)
//			return
//		}
//
//		label := labels[i]
//		values := labelValues[label]
//
//		for _, v := range values {
//			lsCopy := copyLabelSet(ls)
//			lsCopy[label] = v
//			n(i+1, lsCopy)
//		}
//	}
//
//	n(0, prometheus.Labels{})
//	return labelSets
//}

//func copyLabelSet(ls prometheus.Labels) prometheus.Labels {
//	newLs := make(prometheus.Labels, len(ls))
//	for l, v := range ls {
//		newLs[l] = v
//	}
//	return newLs
//}
