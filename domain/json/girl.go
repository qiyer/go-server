package domain

type UnlockBonus struct {
	Level uint `json:"level"`
	Bonus uint `json:"bonus"`
}

type GirlInfo struct {
	Level       uint   `json:"level"`
	Income      uint64 `json:"income"`
	UpgradeCost uint64 `json:"upgrade_cost"`
}

type Girl struct {
	GirlId      uint        `json:"girl_id"`
	Level       uint        `json:"level"`
	Name        string      `json:"name"`
	Unlock      string      `json:"unlock"`
	UnlockBonus UnlockBonus `json:"unlock_bonus"`
	Infos       []GirlInfo  `json:"infos"`
}
