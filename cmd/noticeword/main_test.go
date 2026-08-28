package main

import "testing"

func TestEntrypointPackageBuilds(t *testing.T) {
	if "noticeword" == "" {
		t.Fatal("service name missing")
	}
}
