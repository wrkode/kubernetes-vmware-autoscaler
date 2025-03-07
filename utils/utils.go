package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Fred78290/kubernetes-vmware-autoscaler/types"

	"github.com/Fred78290/kubernetes-vmware-autoscaler/context"
	"gopkg.in/yaml.v2"
	apiv1 "k8s.io/api/core/v1"
)

// MergeKubernetesLabel merge kubernetes map in one
func MergeKubernetesLabel(labels ...types.KubernetesLabel) types.KubernetesLabel {
	merged := types.KubernetesLabel{}

	for _, label := range labels {
		for k, v := range label {
			merged[k] = v
		}
	}

	return merged
}

// NodeFromJSON deserialize a string to apiv1.Node
func NodeFromJSON(s string) (*apiv1.Node, error) {
	data := &apiv1.Node{}

	err := json.Unmarshal([]byte(s), &data)

	return data, err
}

// ToYAML serialize interface to yaml
func ToYAML(v interface{}) string {
	if v == nil {
		return ""
	}

	b, _ := yaml.Marshal(v)

	return string(b)
}

// ToJSON serialize interface to json
func ToJSON(v interface{}) string {
	if v == nil {
		return ""
	}

	b, _ := json.Marshal(v)

	return string(b)
}

// FileExists Check if FileExists
func FileExists(name string) bool {
	if len(name) == 0 {
		return false
	}

	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// MinInt min(a,b)
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// MaxInt max(a,b)
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt64 min(a,b)
func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// MaxInt64 max(a,b)
func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func NewRequestContext(requestTimeout time.Duration) *context.Context {
	return context.NewContext(time.Duration(requestTimeout.Seconds()))
}
