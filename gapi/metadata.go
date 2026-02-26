// simple-bank/gapi/metadata.go
package gapi

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	userAgentHeader      = "user-agent"
	xForwardedFor        = "x-forwarded-for"
	xRealIP              = "x-real-ip"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetatada(ctx context.Context) *Metadata {
	md := &Metadata{}

	if grpcMD, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("==== gRPC Metadata ====")
		for key, values := range grpcMD {
			log.Printf("%s: %v\n", key, values)
		}
		log.Println("=======================")

		// 1️⃣ User-Agent (Gateway ưu tiên trước)
		if ua := grpcMD.Get(grpcGatewayUserAgent); len(ua) > 0 {
			md.UserAgent = ua[0]
		} else if ua := grpcMD.Get(userAgentHeader); len(ua) > 0 {
			md.UserAgent = ua[0]
		}

		// 2️⃣ Client IP (Gateway ưu tiên)
		if ip := grpcMD.Get(xForwardedFor); len(ip) > 0 {
			// x-forwarded-for có thể là "ip1, ip2"
			md.ClientIP = strings.TrimSpace(strings.Split(ip[0], ",")[0])
			return md
		}

		if ip := grpcMD.Get(xRealIP); len(ip) > 0 {
			md.ClientIP = ip[0]
			return md
		}
	}

	// 3️⃣ Fallback cho gRPC thuần
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		addr := p.Addr.String()
		if strings.Contains(addr, ":") {
			md.ClientIP = strings.Split(addr, ":")[0]
		} else {
			md.ClientIP = addr
		}
	}

	return md
}
