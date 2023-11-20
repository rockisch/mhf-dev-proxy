package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mhf-dev-proxy/config"
	"net"
	"os/signal"
	"syscall"
)

func handleConn(ctx context.Context, src net.Conn) {
	var err error

	proxy, err := NewProxy(ctx, src)
	if err != nil {
		log.Println("failed to connect to remote server:", err)
		src.Close()
		return
	}
	go proxy.handleSrc2Dst()
	go proxy.handleDst2Src()
}

func main() {
	srcPort := flag.Int("src", 8090, "port the client will connect to")
	dstPort := flag.Int("dst", 54001, "port the proxy will connect to")
	dstHost := flag.String("host", "127.0.0.1", "host the proxy will connect to")
	mode := flag.Int("mode", int(config.ZZ), "server version")
	flag.Parse()

	config.ProxyConfig.RealClientMode = config.Mode(*mode)
	config.ProxyConfig.ServerPort = *dstPort
	config.ProxyConfig.ServerHost = *dstHost

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	lc := net.ListenConfig{
		Control: reusePort,
	}
	listen, err := lc.Listen(ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", *srcPort))
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		<-ctx.Done()
		listen.Close()
	}()
	log.Println("starting dev proxy")
	for {
		conn, err := listen.Accept()
		if ctx.Err() != nil {
			break
		} else if err != nil {
			log.Println("failed to accept socket", err)
			continue
		}
		log.Println("connection received")
		handleConn(ctx, conn)
	}
}
