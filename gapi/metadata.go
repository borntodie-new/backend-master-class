package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
)

const (
	gRPCUserAgentHeader = "user-agent"
	forwardForHeader    = ":authority"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extraMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %+v\n", md)
		if userAgent := md.Get(gRPCUserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if clientIP := md.Get(forwardForHeader); len(clientIP) > 0 {
			mtdt.ClientIP = clientIP[0]
		}
	}

	// 获取请求的IP信息
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
