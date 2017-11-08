package currencies

const (
	DefaultSymbolPairDelimiter = "-"
)

type SymbolPair struct {
	Left  Symbol
	Right Symbol
}

func NewSymbolPair(left Symbol, right Symbol) SymbolPair {
	return SymbolPair{
		Left:  left,
		Right: right,
	}
}
