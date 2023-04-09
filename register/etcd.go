package register

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegister struct {
	client  *clientv3.Client
	Service string
	lease   int64
}

func NewEtcdRegister(client *clientv3.Client, servicesName string, lease int64) *EtcdRegister {
	return &EtcdRegister{
		client:  client,
		Service: servicesName,
		lease:   lease,
	}
}

func (e *EtcdRegister) Register(ctx context.Context, service *Service) error {
	grant, err := e.client.Lease.Grant(ctx, e.lease)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("/%s/%s/%s", e.Service, service.Name, service.ID)
	value, err := json.Marshal(service)
	if err != nil {
		return err
	}
	e.client.Put(ctx, key, string(value), clientv3.WithLease(grant.ID))
	go e.heartBeat(ctx, grant.ID, key, string(value))
	return nil
}

func (e *EtcdRegister) UnRegister(ctx context.Context, service *Service) error {
	defer func() {
		// 关闭自动续约机制
		e.client.Lease.Close()
	}()
	key := fmt.Sprintf("/%s/%s/%s", e.Service, service.Name, service.ID)
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *EtcdRegister) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := e.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	for {
		if curLeaseID == 0 {
			// TODO 重新注册
		}
		select {
		case k, ok := <-kac:
			// 如果自动续约机制关闭，则退出
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
			fmt.Printf("new TTL:%v\n", k.TTL)
		case <-ctx.Done():
			return
		}
	}
}
