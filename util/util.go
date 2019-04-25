package util

import (
	"github.com/multiformats/go-multiaddr"
)

// StringListToMaddrs 字符串地址列表赚maddrs
func StringListToMaddrs(addrs []string) ([]multiaddr.Multiaddr, error) {
	maddrs := make([]multiaddr.Multiaddr, len(addrs))
	for k, addr := range addrs {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return maddrs, err
		}
		maddrs[k] = maddr
	}
	return maddrs, nil
}
