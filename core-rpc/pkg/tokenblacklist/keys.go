package tokenblacklist

import (
	"crypto/sha256"
	"encoding/hex"
)

// DefaultTTLSeconds 与 refresh token 最长有效期一致，用于 Setex。
const DefaultTTLSeconds = 604800 // 7 天

func Key(kind string, token string) string {
	sum := sha256.Sum256([]byte(token))
	return "auth:blk:" + kind + ":" + hex.EncodeToString(sum[:])
}

func AccessKey(token string) string {
	return Key("access", token)
}

func RefreshKey(token string) string {
	return Key("refresh", token)
}
