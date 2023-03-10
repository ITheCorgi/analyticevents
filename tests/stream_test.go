//go:build integration_tests
// +build integration_tests

package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ITheCorgi/analyticevents/internal/entity"
)

func TestAddEvents(t *testing.T) {
	batch := `{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	}`

	w := strings.NewReader(batch)

	client := http.Client{}
	eventStreamUri := fmt.Sprintf("http://%s:%s/analytics/event/streaming", configApp.App.Host, configApp.App.Port)

	req, _ := http.NewRequest(http.MethodPost, eventStreamUri, w)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got not 200 status code: %s", resp.Status)
	}

	m := &entity.AnalyticEvent{}
	inserted, err := db.NewSelect().Model(m).Count(context.Background())
	if err != nil {
		t.Error(err)
	}

	if inserted != 30 {
		t.Errorf("expected inserts %d, got %d", 30, inserted)
	}
}
