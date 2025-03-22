package rest

import "sync"

var (
	refreshTokenStore   = make(map[string]int64)
	refreshTokenStoreMu sync.RWMutex
)

// StoreRefreshToken сохраняет refresh токен, связывая его с userID.
func StoreRefreshToken(token string, userID int64) {
	refreshTokenStoreMu.Lock()
	defer refreshTokenStoreMu.Unlock()
	refreshTokenStore[token] = userID
}

// GetUserIDByRefreshToken возвращает userID для данного refresh токена.
func GetUserIDByRefreshToken(token string) (int64, bool) {
	refreshTokenStoreMu.RLock()
	defer refreshTokenStoreMu.RUnlock()
	uid, ok := refreshTokenStore[token]
	return uid, ok
}

// RemoveRefreshToken удаляет refresh токен из хранилища.
func RemoveRefreshToken(token string) {
	refreshTokenStoreMu.Lock()
	defer refreshTokenStoreMu.Unlock()
	delete(refreshTokenStore, token)
}
