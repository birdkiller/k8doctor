package cleaner

import (
	"testing"
)

func TestClean_OOMCase(t *testing.T) {
	input := "Pod一直重启，日志报OOM killed，容器是java应用"
	s := Clean(input)

	if len(s.States) == 0 {
		t.Error("Expected to detect OOMKilled state")
	}

	if len(s.Resources) == 0 {
		t.Error("Expected to detect Pod resource")
	}

	if s.Context == "" {
		t.Error("Expected Context to be non-empty")
	}

	t.Logf("Input: %s", input)
	t.Logf("States: %v", s.States)
	t.Logf("Resources: %v", s.Resources)
	t.Logf("Keywords: %v", s.Keywords)
	t.Logf("Context: %s", s.Context)
}

func TestClean_CrashLoopCase(t *testing.T) {
	input := "Pod处于CrashLoopBackOff状态，一直重启"
	s := Clean(input)

	if len(s.States) == 0 {
		t.Error("Expected to detect CrashLoopBackOff state")
	}

	t.Logf("Input: %s", input)
	t.Logf("States: %v", s.States)
	t.Logf("Context: %s", s.Context)
}

func TestClean_NodeCase(t *testing.T) {
	input := "Node变成NotReady状态了，集群有个节点挂了"
	s := Clean(input)

	if len(s.States) == 0 {
		t.Error("Expected to detect NotReady state")
	}

	if len(s.Resources) == 0 {
		t.Error("Expected to detect Node resource")
	}

	t.Logf("Input: %s", input)
	t.Logf("States: %v", s.States)
	t.Logf("Resources: %v", s.Resources)
	t.Logf("Context: %s", s.Context)
}

func TestClean_ErrorCode(t *testing.T) {
	input := "容器退出码是137，一直崩溃"
	s := Clean(input)

	if len(s.ErrorCodes) == 0 {
		t.Error("Expected to detect error code 137")
	}

	t.Logf("Input: %s", input)
	t.Logf("ErrorCodes: %v", s.ErrorCodes)
}

func TestClean_NetworkCase(t *testing.T) {
	input := "Service访问超时，502错误，Ingress返回网关超时"
	s := Clean(input)

	t.Logf("Input: %s", input)
	t.Logf("States: %v", s.States)
	t.Logf("Keywords: %v", s.Keywords)
}
