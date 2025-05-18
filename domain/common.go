package domain

import (
	"encoding/json"
	"fmt"
	data "go-server/data"
	domain "go-server/domain/json"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	BaseData          = "tasks"
	OfflineIncomeBase = 50  //除以该系数
	SecCoinBase       = 100 //每秒获得的基础金币
	Level1            = 10
	Level2            = 20
	Level3            = 114
	CostBase          = 1.0465
	InitGirls         = "10001:1,10002:0"
)

var Apartments []domain.Apartment

var Girls []domain.Girl

func InitJsons() {

	// 读取嵌入的 JSON 文件
	apartment_data, err := data.ConfigJsonsFile.ReadFile("apartment.json")
	if err != nil {
		log.Fatal("读取嵌入文件失败:", err)
	}

	err = json.Unmarshal(apartment_data, &Apartments)
	if err != nil {
		log.Fatalf("解析 JSON 失败: %v", err)
	}

	fmt.Printf("配置内容：%+v\n", Apartments[0])

	gril_data, err2 := data.ConfigJsonsFile.ReadFile("girls.json")
	if err2 != nil {
		log.Fatal("读取嵌入文件失败:", err2)
	}

	err2 = json.Unmarshal(gril_data, &Girls)
	if err2 != nil {
		log.Fatalf("解析 JSON 失败: %v", err2)
	}

	fmt.Printf("配置内容 秘书：%+v\n", Girls[0])
}

func GetOfflineCoin(secCoin uint64, time uint64) (coin uint64) {
	return secCoin * time / OfflineIncomeBase
}

func GetSecCoin(user User) (coin uint64) {
	var base uint64 = SecCoinBase
	//实际数据需要读表
	for _, gril := range ParseGirls(user.Girls) {
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

func ParseGirls(str string) (grils []MGirl) {
	// 分割并清理数据
	parts := strings.FieldsFunc(str, func(r rune) bool {
		return r == ','
	})

	var gs []MGirl
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		pair := strings.Split(trimmed, ":")
		if len(pair) == 2 {

			gs = append(gs, MGirl{
				GirlId: StrToUint(pair[0]),
				Level:  StrToUint64(pair[1]),
			})
		}
	}
	return gs
}

func StrToUint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 64) // 参数：字符串, 进制（十进制）, 位数（64位）
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0
	}
	return uint(num) // 将 uint64 显式转换为 uint
}

func StrToUint64(str string) uint64 {
	num, err := strconv.ParseUint(str, 10, 64) // 参数：字符串, 进制（十进制）, 位数（64位）
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0
	}
	return num
}
func RoleLevelCost(curLevel int) (coin uint64) {
	var cost uint64 = Level1
	//实际数据需要读表
	if curLevel == 1 {
		cost = Level1
	} else if curLevel == 2 {
		cost = Level2
	} else if curLevel == 3 {
		cost = Level3
	} else {
		var coin = Level3 * math.Pow(CostBase, float64(curLevel-3))
		cost = uint64(coin)
	}

	return cost
}

func GirlLevelCost(roleId uint, curLevel int) (coin uint64) {
	var cost uint64 = 0

	for _, gril := range Girls {
		if gril.GirlId == roleId {
			for _, info := range gril.Infos {
				if info.Level == uint(curLevel) {
					cost = info.UpgradeCost
					break
				}
			}
			break
		}
	}
	return cost
}

func GirlLevelUpCheckNeeds(roleId uint, user User) (success bool) {
	return true
}
