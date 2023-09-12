package auth

import (
	"context"
	"fmt"
)

// kAuthIpContextKeyName is the name of the ip address field
// in the context values.
const kAuthIpContextKeyName = "ip_address"

// SetIpAddress adds the given ip address to the provided context
// returning a new derived context.
func SetIpAddress(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, kAuthIpContextKeyName, ip)
}

// GetIpOrErr extracts the active ip address from the provided context
// and returns an error if the ip address can not be extracted.
func GetIpOrErr(ctx context.Context) (string, error) {
	value := ctx.Value(kAuthIpContextKeyName)
	if value == nil {
		return "", fmt.Errorf("no ip address")
	}
	ip, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("ip address type mismatch")
	}
	return ip, nil
}
