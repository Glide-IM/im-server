package rpc

import (
	"context"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

type selector struct {
	services map[string]string
	round    client.Selector
	tags     map[string]string
}

func newSelector() *selector {
	s := map[string]string{}
	return &selector{
		services: s,
		round:    NewRoundRobinSelector(),
		tags:     map[string]string{},
	}
}

func (r *selector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)

	if target, ok := m["ExtraTarget"]; ok {
		if _, ok := r.services[target]; ok {
			return target
		}
		logger.E("unknown service addr, ExtraTarget:", target)
	}

	if tag, ok := m["ExtraTag"]; ok {
		if path, ok := r.tags[tag]; ok {
			if _, ok := r.services[path]; ok {
				logger.D("route by tag: %s=%s", tag, path)
				return path
			}
		}
	}
	return r.round.Select(ctx, servicePath, serviceMethod, args)
}

func (r *selector) UpdateServer(servers map[string]string) {
	r.round.UpdateServer(servers)
	for k, v := range servers {
		r.services[k] = v
	}
}
