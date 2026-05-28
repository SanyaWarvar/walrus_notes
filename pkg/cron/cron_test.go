package cron

import "testing"

func TestCron_AddFunc_ValidSpec(t *testing.T) {
	c := NewCron()
	err := c.AddFunc("@hourly", func() {})
	if err != nil {
		t.Fatalf("AddFunc() error = %v", err)
	}
}

func TestCron_AddFunc_InvalidSpec(t *testing.T) {
	c := NewCron()
	err := c.AddFunc("not a valid cron spec", func() {})
	if err == nil {
		t.Fatal("expected error for invalid cron spec")
	}
}

func TestCron_StartStop(t *testing.T) {
	c := NewCron()
	c.Start()
	c.Stop()
}
