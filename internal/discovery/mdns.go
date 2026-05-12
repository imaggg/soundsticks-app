package discovery

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grandcat/zeroconf"
)

type Device struct {
	IP       string
	Hostname string
	Port     int
}

// Browse discovers SoundSticks devices via mDNS (_jbl-product._tcp.local).
// Sends confirmed devices to the returned channel until ctx is done or closed.
func Browse(ctx context.Context) (<-chan Device, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("mdns resolver: %w", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	out := make(chan Device, 4)

	go func() {
		defer close(out)
		if err := resolver.Browse(ctx, "_jbl-product._tcp", "local.", entries); err != nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-entries:
				if !ok {
					return
				}
				if len(entry.AddrIPv4) == 0 {
					continue
				}
				ip := entry.AddrIPv4[0].String()
				if confirmDevice(ip) {
					out <- Device{
						IP:       ip,
						Hostname: entry.HostName,
						Port:     entry.Port,
					}
				}
			}
		}
	}()

	return out, nil
}

// confirmDevice verifies the device by fetching its UPnP description (no auth required).
func confirmDevice(ip string) bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s:59152/description.xml", ip))
	if err != nil {
		return false
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
