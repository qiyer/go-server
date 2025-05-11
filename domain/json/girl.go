package domain

type GirlLimit struct {
	Level uint   `json:"level"`
	Need  string `json:"need"`
}

type GirlInfo struct {
	Level       uint   `json:"level"`
	Income      uint64 `json:"income"`
	Bonus       string `json:"bonus"`
	UpgradeCost uint64 `json:"upgrade_cost"`
}

type Girl struct {
	GirlId uint        `json:"girl_id"`
	Level  uint        `json:"level"`
	Name   string      `json:"name"`
	Limits []GirlLimit `json:"limits"`
	Infos  []GirlInfo  `json:"infos"`
}
