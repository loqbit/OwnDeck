package discovery

import "testing"

func TestClientStatus_Detected(t *testing.T) {
	s := ClientStatus(true)
	if s != "detected" {
		t.Errorf("ClientStatus(true) = %q, want %q", s, "detected")
	}
}

func TestClientStatus_NotFound(t *testing.T) {
	s := ClientStatus(false)
	if s != "not found" {
		t.Errorf("ClientStatus(false) = %q, want %q", s, "not found")
	}
}
