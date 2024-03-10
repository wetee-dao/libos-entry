package libos

import (
	"testing"
)

func TestGetLocalReport(t *testing.T) {
	appId := "testAppId"
	fs := &MockFs{}

	_, _, _, err := GetLocalReport(appId, fs)
	if err != nil {
		t.Errorf("GetLocalReport failed: %v", err)
	}
}
