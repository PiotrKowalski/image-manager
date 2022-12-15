package remote

const (
	KV_STORE_PORT string = "KV_STORE_PORT"
	KV_STORE_HOST string = "KV_STORE_HOST"
)

type SecretRemoteConfigProvider interface {
	LoadStoreConfig() error
	StartRemoteWatch() error
	StopRemoteWatch() error
}
