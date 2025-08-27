package opcua

import (
	"context"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/uasc"
)

// Connect creates an OPC UA client.
func Connect(ctx context.Context, url, securityMode, securityPolicy, username, password string) (*opcua.Client, error) {
	opts := []opcua.Option{
		opcua.SecurityModeString(securityMode),
		opcua.SecurityPolicy(securityPolicy),
		opcua.AutoReconnect(true),
	}
	if username != "" {
		auth := uasc.UserNameIdentityToken{UserName: username, Password: password}
		opts = append(opts, opcua.Authenticate(&auth))
	}
	c := opcua.NewClient(url, opts...)
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}
	return c, nil
}
