package plugins

import (
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"strconv"
)

const Name = "alice-plugin"

type Plugin struct{}

var _ framework.QueueSortPlugin = &Plugin{}

func (p *Plugin) Name() string {
	return Name
}

func (*Plugin) Less(p1, p2 *framework.PodInfo) bool {
	gp1 := getGroupPriority(p1.Pod)
	gp2 := getGroupPriority(p2.Pod)

	if gp1 == gp2 {
		return compQOS(p1.Pod, p2.Pod)
	} else {
		return gp1 > gp2
	}
}

func getGroupPriority(p *v1.Pod) int {
	value, ok := p.Labels["groupPriority"]
	if ok == true {
		num, err := strconv.Atoi(value)
		if err == nil {
			return num
		} else {
			return -10
		}
	} else {
		return -10
	}
}

func compQOS(p1, p2 *v1.Pod) bool {
	p1QOS, p2QOS := v1qos.GetPodQOS(p1), v1qos.GetPodQOS(p2)
	if p1QOS == v1.PodQOSGuaranteed {
		return true
	}
	if p1QOS == v1.PodQOSBurstable {
		return p2QOS != v1.PodQOSGuaranteed
	}
	return p2QOS == v1.PodQOSBestEffort
}