package domain

type Apartment struct {
	Level       uint   `json:"level"`
	Name        string `json:"name"`
	Bonus       string `json:"bonus"`
	UpgradeCost uint64 `json:"upgrade_cost"`
}
