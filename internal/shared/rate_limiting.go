package shared

import (
	"sync"
	"time"

	beecontext "github.com/beego/beego/v2/server/web/context"
	"golang.org/x/time/rate"
)

const visitorTTL = 10 * time.Minute

type IPStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var store = &IPStore{
	visitors: make(map[string]*visitor),
}

func (i *IPStore) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	now := time.Now()
	record, exists := i.visitors[ip]
	if !exists {
		record = &visitor{
			// limit 5 RPS(request per second), burst capacity: 10 requests
			limiter: rate.NewLimiter(rate.Limit(5), 10),
		}
		i.visitors[ip] = record
	}
	record.lastSeen = now

	i.cleanup(now)

	return record.limiter
}

func (i *IPStore) cleanup(now time.Time) {
	for ip, record := range i.visitors {
		if now.Sub(record.lastSeen) > visitorTTL {
			delete(i.visitors, ip)
		}
	}
}

func IPLimiterFilter(ctx *beecontext.Context) {
	ip := ctx.Input.IP()

	limiter := store.GetLimiter(ip)

	if !limiter.Allow() {
		WriteError(ctx, NewManyRequestError("many_request", "too many requests"))
		return
	}
}
