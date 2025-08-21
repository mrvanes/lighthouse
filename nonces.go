package lighthouse

import (
	"time"
	"github.com/google/uuid"
)

type nonces map[string]time.Time

func (n nonces) addNonce() string {
	n.clear()
	nonce := uuid.New().String()
	time := time.Now()
	n[nonce] = time
	return nonce
}

func (n nonces) deleteNonce(nonce string) {
	delete(n, nonce)
}

func (n nonces) clear() {
	stale := time.Now().Add(time.Minute *-1)
	for nonce, time := range n {
		if time.Before(stale) {
			delete(n, nonce)
		}
	}
}

var Nonces nonces = make(nonces)
