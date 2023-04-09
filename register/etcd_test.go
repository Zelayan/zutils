package register

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestNewEtcdRegister(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"zeze.com:2379"},
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		t.Fatalf("create etcd client failed: %v", err)
	}
	service := &Service{
		ID:   "0",
		Name: "server01",
		IP:   "127.0.0.1",
		Port: 8080,
	}
	register := NewEtcdRegister(client, "haha", 5)
	err = register.Register(context.Background(), service)
	time.Sleep(time.Duration(10) * time.Second)
	register.UnRegister(context.Background(), service)
	t.Log("unRegister")
	time.Sleep(time.Duration(10) * time.Second)
}
