package plugins

import (
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"strconv"
	"k8s.io/apimachinery/pkg/runtime"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const Name = "alice-scheduler"

type Plugin struct{
	handle framework.FrameworkHandle
}

var _ framework.QueueSortPlugin = &Plugin{}
var _ framework.PreFilterPlugin = &Plugin{}

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

func (p *Plugin) PreFilter(_ context.Context, _ *framework.CycleState, pod *v1.Pod) *framework.Status {
	groupName, minAvailable, err := getPodGroupLabels(pod)
	if groupName == "" {
		return framework.NewStatus(framework.Success, "Pass PreFilter, no groupName")
	}
	if (minAvailable == 0) && (err == nil) {
		return framework.NewStatus(framework.Success, "Pass PreFilter, no minAvailable")
	}
	if err != nil {
		return framework.NewStatus(framework.Unschedulable, "Failed PreFilter, fail in reading labels")
	}
	num := p.getPodNumByGroupName(pod.Namespace, groupName)

	if num >= minAvailable {
		return framework.NewStatus(framework.Success, "Pass PreFilter")
	} else {
		return framework.NewStatus(framework.Unschedulable, "Failed PreFilter")
	}
}

func getPodGroupLabels(pod *v1.Pod) (string, int, error) {
	podGroupName, exist := pod.Labels["podGroup"]
	if !exist || len(podGroupName) == 0 {
		return "", 0, nil
	}
	minAvailable, exist := pod.Labels["minAvailable"]
	if !exist || len(minAvailable) == 0 {
		return "", 0, nil
	}
	minNum, err := strconv.Atoi(minAvailable)
	if err != nil {
		return "", 0, err
	}
	if minNum < 1 {
		return "", 0, err
	}
	return podGroupName, minNum, nil
}

func (p *Plugin)getPodNumByGroupName(ns string, pg string) int {
	podList, err := p.handle.ClientSet().CoreV1().Pods(ns).List(metav1.ListOptions{})
	// podList.Items是一個裝很多v1.pod的陣列

	if err != nil {
		return 0
	}

	var num int = 0
	for _, pod := range podList.Items {
		if pod.Labels["podGroup"] == pg {
			num++
		}
	}
	return num
}


func (*Plugin)PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func New(_ *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	return &Plugin{
		handle: handle,
	}, nil
}
