package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mhf-dev-proxy/config"
	"mhf-dev-proxy/network"
	"net"
	"sync"
	"time"
)

const initialReqs = 77

type Proxy struct {
	*sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	src    *network.CryptConn
	dst    *network.CryptConn
	resend bool
}

func NewProxy(ctx context.Context, src net.Conn) (Proxy, error) {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	proxy := Proxy{
		Mutex:  &sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
		src:    network.NewCryptConn(src),
	}
	err = proxy.connectDst()
	if err != nil {
		return proxy, err
	}
	return proxy, nil
}

func (p *Proxy) close() {
	p.cancel()
	p.src.Close()
	p.dst.Close()
}

func (p *Proxy) connectDst() error {
	log.Println("connecting to server...")
	dst, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.ProxyConfig.ServerHost, config.ProxyConfig.ServerPort))
	if err != nil {
		return err
	}
	log.Println("connected!")
	p.dst = network.NewCryptConn(dst)
	return nil
}

func (p *Proxy) attemptReconnect() bool {
	log.Println("disconnected, will attempt to reconnect")
	var err error
	ticker := time.Tick(time.Second * 2)
	for {
		select {
		case <-ticker:
			err = p.connectDst()
			if err != nil {
				continue
			}
			return true
		case <-p.ctx.Done():
			return false
		}
	}
}

func (p *Proxy) handleSrc2Dst() {
	defer func() {
		log.Println("stopping from src->dst")
		p.close()
	}()
	var initialPackets [][]byte
	for {
		packet, err := p.src.ReadPacket()
		if err != nil {
			return
		}
		if len(initialPackets) < initialReqs {
			initialPackets = append(initialPackets, bytes.Clone(packet))
		}
		p.Lock()
		if p.resend {
			for _, initialPacket := range initialPackets {
				if err := p.dst.SendPacket(initialPacket); err != nil {
					log.Println("failed to resync, closing")
					return
				}
			}
			p.resend = false
		}
		p.dst.SendPacket(packet)
		p.Unlock()
	}
}

func (p *Proxy) handleDst2Src() {
	defer func() {
		log.Println("stopping from dst->src")
		p.close()
	}()
	var err error
	for {
		var packet []byte
		for {
			packet, err = p.dst.ReadPacket()
			if err == nil {
				break
			} else {
				p.Lock()
				p.resend = true
				result := p.attemptReconnect()
				p.Unlock()
				if !result {
					return
				}
			}
		}
		err = p.src.SendPacket(packet)
		if err != nil {
			return
		}
	}
}
