package middlewares

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/kotche/url-shortening-service/internal/app/config"
)

const accessProhibited = "Access to the internal network is prohibited"

type TrustedNetwork struct {
	TrustedSubnet net.IP
}

func NewTrustedNetwork(cfg *config.Config) *TrustedNetwork {
	trustedSubnet := net.ParseIP(cfg.TrustedSubnet)
	return &TrustedNetwork{
		TrustedSubnet: trustedSubnet,
	}
}

// TrustedNetworkHandler checks whether the client's IP address is included in the trusted subnet
func (t *TrustedNetwork) TrustedNetworkHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t.TrustedSubnet == nil {
			log.Print("middlewares TrustedNetworkHandler: empty TrustedSubnet")
			http.Error(w, accessProhibited, http.StatusForbidden)
			return
		}
		ipStr := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipStr)
		if ip == nil {
			ipStr = r.Header.Get("X-Forwarded-For")
			ipStrs := strings.Split(ipStr, ",")
			ipStr = ipStrs[0]
			ip = net.ParseIP(ipStr)
		}
		if !bytes.Equal(t.TrustedSubnet, ip) {
			log.Printf("middlewares TrustedNetworkHandler: TrustedSubnet - %s, ip - %s", t.TrustedSubnet, ip)
			http.Error(w, accessProhibited, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
