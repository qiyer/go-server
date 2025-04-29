package domain

const (
	BaseData          = "tasks"
	OfflineIncomeBase = 50  //除以该系数
	SecCoinBase       = 100 //每秒获得的基础金币
)

func GetOfflineCoin(secCoin uint64, time uint64) (coin uint64) {
	return secCoin * time / OfflineIncomeBase
}

func GetSecCoin(user User) (coin uint64) {
	var base uint64 = SecCoinBase
	//实际数据需要读表
	for _, gril := range user.Grils {
		base += gril.Level * 10
	}
	for _, island := range user.Islands {
		base += uint64(island.Level) * 10
	}
	var index uint64 = 1
	for _, legacy := range user.Legacies {
		index += uint64(legacy.Level) * 10
	}

	return base * index
}
