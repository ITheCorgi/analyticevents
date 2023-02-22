package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ITheCorgi/analyticevents/internal/entity"
)

func Test_AddEvents(t *testing.T) {
	defer func() {
		m := []entity.AnalyticEvent{}
		_, err := db.NewTruncateTable().Model(&m).Exec(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	client := http.Client{}
	input := `{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}
	{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0, "param_str":"some text"}`

	w := strings.NewReader(input)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/analytics/event/streaming", configApp.App.Host, configApp.App.Port), w)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	_, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	var (
		cnt   int
		model []entity.AnalyticEvent
	)

	cnt, err = db.NewSelect().Model(&model).Count(context.Background())
	if err != nil {
		t.Error(err)
	}

	if cnt != 6 {
		t.Errorf("expected %d, got %d", 6, cnt)
	}
}
