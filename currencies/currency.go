package currencies

type Currency struct {
	Name   string `json:"name"`
	Symbol Symbol `json:"symbol"`
}
