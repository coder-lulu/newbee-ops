package config

import (
	"github.com/zeromicro/go-zero/core/conf"
	"testing"
)

func TestWorkerManagerConfiguration(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  bool
	}{
		{`{}`, true},
		{`{"WorkerManager":{"Enabled":false}}`, false},
		{`{"WorkerManager":{"Enabled":true}}`, true},
	} {
		var cfg struct {
			WorkerManager WorkerManagerConf `json:",optional"`
		}
		if err := conf.LoadFromJsonBytes([]byte(tc.input), &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.WorkerManager.IsEnabled() != tc.want {
			t.Fatalf("config %s: enabled=%v want=%v", tc.input, cfg.WorkerManager.IsEnabled(), tc.want)
		}
	}
}
