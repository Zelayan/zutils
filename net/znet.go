package znet

import (
	"fmt"
	"net"
	"time"
)

type Pinger struct {
    ip      net.IP
    timeout time.Duration
}

func NewPinger(ip string, timeout time.Duration) (*Pinger, error) {
    parsedIP := net.ParseIP(ip)
    if parsedIP == nil {
        return nil, fmt.Errorf("invalid IP address")
    }
    return &Pinger{
        ip:      parsedIP,
        timeout: timeout,
    }, nil
}

func (p *Pinger) Ping() error {
    conn, err := net.DialTimeout(p.getNetwork(), p.ip.String(), p.timeout)
    if err != nil {
        return fmt.Errorf("Ping failed: %s", err.Error())
    }
    conn.Close()
    return nil
}

func (p *Pinger) getNetwork() string {
    if p.ip.To4() != nil {
        return "ip4:icmp"
    }
    return "ip6:ipv6-icmp"
}