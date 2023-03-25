package znet

import (
	"net"
	"testing"
	"time"
)

func TestPinger_Ping(t *testing.T) {
	type fields struct {
		ip      net.IP
		timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "127.0.0.1",
			fields: fields{
				ip: net.ParseIP("127.0.0.1"),
				timeout: time.Duration(1) * time.Second,
			},
			wantErr: false,
		},
		{
			name: "42.193.97.239",
			fields: fields{
				ip: net.ParseIP("42.193.97.239"),
				timeout: time.Duration(1) * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pinger{
				ip:      tt.fields.ip,
				timeout: tt.fields.timeout,
			}
			if err := p.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Pinger.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
