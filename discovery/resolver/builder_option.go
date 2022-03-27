package resolver

import "google.golang.org/grpc/grpclog"

type BuilderOptApplyFn func(builder *builder)

func WithLogger(logger grpclog.LoggerV2) BuilderOptApplyFn {
	return func(builder *builder) {
		builder.logger = logger
	}
}

func WithIPv4Whitelist(subnets []string) BuilderOptApplyFn {
	return func(builder *builder) {
		builder.whitelistSubnets = subnets
	}
}
