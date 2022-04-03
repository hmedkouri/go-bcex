package bcex

type Fees struct {
	MakerRate float64 `json:"makerRate"`
	TakerRate float64 `json:"takerRate"`
	VolumeInUSD float64 `json:"volumeInUSD"`
}