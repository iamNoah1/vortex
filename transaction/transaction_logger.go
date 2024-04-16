package transaction

type TransactionLogger interface {
	Put(key string, value string) error
	Delete(key string) error
	ReplayEvents() error
}
