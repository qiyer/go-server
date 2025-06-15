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

type GirlBaseInfo struct {
	LLevel     uint   `json:"l_level"`
	RLevel     uint   `json:"r_level"`
	LevelIndex uint64 `json:"level_index"`
	LevelBase  uint64 `json:"level_base"`
	CoinIndex  uint64 `json:"coin_index"`
	CoinBase   uint64 `json:"coin_base"`
}

type Girl struct {
	GirlId      uint           `json:"girl_id"`
	Level       uint           `json:"level"`
	Name        string         `json:"name"`
	Unlock      string         `json:"unlock"`
	UnlockBonus UnlockBonus    `json:"unlock_bonus"`
	Infos       []GirlInfo     `json:"infos"`
	BaseInfos   []GirlBaseInfo `json:"binfos"`
}

type QuickEarn struct {
	Level uint   `json:"level"`
	Bonus uint64 `json:"bonus"`
}

type Vehicle struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Income      uint64 `json:"income"`
	UpgradeCost uint64 `json:"upgrade_cost"`
	NeedLevel   int    `json:"need_level"`
}

type Capital struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint64 `json:"price"` // 购买价格
	Bonus uint   `json:"bonus"` // 每秒收益
}

type DayBonus struct {
	Day   int    `json:"day"`
	Type  string `json:"type"`
	Bonus uint   `json:"bonus"`
}

type OnlineBonus struct {
	Id    int    `json:"id"`
	Min   int    `json:"day"`
	Type  string `json:"type"`
	Bonus uint   `json:"bonus"`
}

type Boss struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Damage uint64 `json:"damage"` // 伤害
	Bonus  uint   `json:"bonus"`  // 奖励
}
