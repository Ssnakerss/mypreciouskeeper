package grpcClient

import "context"

func (c *GRPCClient) Ping(ctx context.Context) (int64, error) {
	resp, err := c.PingClient.Ping(ctx, nil)
	if err != nil {
		return 0, err
	}
	return resp.Resp, err
}
