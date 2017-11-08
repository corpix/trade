package currencies

type Mapping map[Currency]Currency

func NewMappingFromIntersection(keys Currencies, values Currencies) Mapping {
	var (
		cm = Mapping{}
	)

	for _, k := range keys {
		for _, v := range values {
			if k.Name == v.Name {
				cm[k] = v
			}
		}
	}

	return cm
}
