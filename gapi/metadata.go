package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gatewayUserAgentHeader = "grpcgateway-user-agent"
	gRPCUserAgentHeader    = "user-agent"
	forwardForHeader       = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extraMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(gatewayUserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if userAgent := md.Get(gRPCUserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if clientIPs := md.Get(forwardForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	// 获取请求的IP信息
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
