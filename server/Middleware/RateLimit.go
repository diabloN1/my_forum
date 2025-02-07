package middleware

import (
	Cruds "forum/Api"
	"net/http"
	"sync"
	"time"
)

var (
	requests     = make(map[string][]time.Time) 
	mu           sync.Mutex                     
	maxRequests  = 5                            
	timeInterval = time.Second                  
)


func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr 

		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		timestamps := requests[ip]

		var newTimestamps []time.Time
		for _, t := range timestamps {
			if now.Sub(t) < timeInterval {
				newTimestamps = append(newTimestamps, t)
			}
		}

		newTimestamps = append(newTimestamps, now)
		requests[ip] = newTimestamps

		if len(newTimestamps) > maxRequests {
			Cruds.ShowError(w,"Too many requests.", 429)
			return
		}

		next(w, r)
	}
}
