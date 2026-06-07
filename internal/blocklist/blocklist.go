package blocklist

type Blocklist interface {
	IsBlocked(key string) (bool, error)
}
