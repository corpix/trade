package currencies

const (
	MapperDirectionCommon = "common"
	MapperDirectionMarket = "market"
)

type MapperDirection struct {
	From string
	To   string
}

func (d MapperDirection) String() string {
	return d.From + "->" + d.To
}
