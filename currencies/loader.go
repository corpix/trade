package currencies

type Loader interface {
	Common() (Currencies, error)
	Market(string) (Currencies, error)
}
