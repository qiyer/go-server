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
	"time"
)

const (
	BaseData           = "tasks"
	OfflineIncomeBase  = 50  //除以该系数
	SecCoinBase        = 100 //每秒获得的基础金币
	Level1             = 10
	Level2             = 20
	Level3             = 114
	CostBase           = 1.0465
	BaseCaiShen        = 4
	CaiShenGrowth1     = 1.011
	CaiShenGrowth2     = 1.0487
	TimesBonusBaseTime = 300 // 5分钟
)

var InitGirls = []string{"10001:1", "10002:0"}

var Apartments []domain.Apartment

var Girls []domain.Girl

var QuickEarns []domain.QuickEarn

var Vehicles []domain.Vehicle

var Capitals []domain.Capital

var DayBonuses []domain.DayBonus

var OnlineBonuses []domain.OnlineBonus

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

	quickEarn_data, err3 := data.ConfigJsonsFile.ReadFile("quick_earn.json")
	if err3 != nil {
		log.Fatal("读取嵌入文件失败:", err3)
	}

	err3 = json.Unmarshal(quickEarn_data, &QuickEarns)
	if err3 != nil {
		log.Fatalf("解析 JSON 失败: %v", err3)
	}

	fmt.Printf("配置内容 快速搞钱：%+v\n", Girls[0])

	Vehicle_data, err4 := data.ConfigJsonsFile.ReadFile("vehicle.json")
	if err4 != nil {
		log.Fatal("读取嵌入文件失败:", err4)
	}

	err4 = json.Unmarshal(Vehicle_data, &Vehicles)
	if err4 != nil {
		log.Fatalf("解析 JSON 失败: %v", err4)
	}

	fmt.Printf("配置内容 坐骑：%+v\n", Vehicles[0])

	Capital_data, err5 := data.ConfigJsonsFile.ReadFile("capital.json")
	if err5 != nil {
		log.Fatal("读取嵌入文件失败:", err5)
	}

	err5 = json.Unmarshal(Capital_data, &Capitals)
	if err5 != nil {
		log.Fatalf("解析 JSON 失败: %v", err5)
	}

	fmt.Printf("配置内容 资产：%+v\n", Capitals[0])

	DayBonus_data, err6 := data.ConfigJsonsFile.ReadFile("days_bonus.json")
	if err6 != nil {
		log.Fatal("读取嵌入文件失败:", err6)
	}

	err6 = json.Unmarshal(DayBonus_data, &DayBonuses)
	if err6 != nil {
		log.Fatalf("解析 JSON 失败: %v", err6)
	}

	fmt.Printf("配置内容 7天奖励：%+v\n", DayBonuses[0])

	OnlineBonuse_data, err7 := data.ConfigJsonsFile.ReadFile("online_bonus.json")
	if err7 != nil {
		log.Fatal("读取嵌入文件失败:", err7)
	}

	err7 = json.Unmarshal(OnlineBonuse_data, &OnlineBonuses)
	if err7 != nil {
		log.Fatalf("解析 JSON 失败: %v", err7)
	}

	fmt.Printf("配置内容 在线奖励：%+v\n", OnlineBonuses[0])
}

func GetOfflineCoin(secCoin uint64, time uint64) (coin uint64) {
	return secCoin * time / OfflineIncomeBase
}

func GetSecCoin(user User) (coin uint64) {
	var base uint64 = SecCoinBase
	//实际数据需要读表
	for _, gril := range ParseGirls(user.Girls) {
		for _, gl := range Girls {
			if gl.GirlId == gril.GirlId {
				for _, info := range gl.Infos {
					if info.Level == gl.Level {
						base += info.Income
						break
					}
				}
				break
			}
		}
	}

	var vehicleLevel = int(user.Vehicle.Level)
	for i := 0; i < vehicleLevel; i++ {
		for _, vehicle := range Vehicles {
			if vehicle.ID == uint(i+1) {
				base += vehicle.Income
				break
			}
		}
	}

	var index uint64 = 1
	for _, legacy := range user.Legacies {
		index += uint64(legacy.Level) * 10
	}

	return base * index
}

func CheckInDays(user User, daystr int) (days []string) {
	str := fmt.Sprintf("%d", daystr)
	//实际数据需要读表
	for _, day := range user.Days {
		if strings.Contains(day, str) {
			return user.Days
		}
	}

	return append(user.Days, fmt.Sprintf("%d:0", daystr))
}

func CheckIn(user User, daystr string) (isCheck bool, days []string) {
	if len(user.Days) > 7 {
		return false, user.Days
	}

	var _days []string

	//实际数据需要读表
	for _, day := range user.Days {
		if strings.Contains(day, daystr) {
			parts := strings.FieldsFunc(day, func(r rune) bool {
				return r == ':'
			})
			if parts[1] == "0" {
				return true, append(_days, fmt.Sprintf("%s:1", daystr))
			}
			if parts[1] == "1" {
				return true, append(_days, fmt.Sprintf("%s:2", daystr))
			}
		} else {
			_days = append(_days, day)
		}
	}

	return false, user.Days
}

func ParseGirls(parts []string) (grils []MGirl) {
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

func GirlUnlockCheckNeeds(roleId uint, user User) (success bool) {
	var unlockStr string = ""
	for _, gril := range Girls {
		if gril.GirlId == roleId {
			unlockStr = gril.Unlock
			break
		}
	}
	fmt.Println("GirlUnlockCheckNeeds unlockStr:", unlockStr)

	if unlockStr == "" {
		return true
	} else {
		parts := strings.FieldsFunc(unlockStr, func(r rune) bool {
			return r == ','
		})

		var roleLevel uint = 0
		var girlId uint = 0
		var girlLevel uint = 0

		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			pair := strings.Split(trimmed, "=")
			if len(pair) == 2 {
				if pair[0] == "10000" {
					roleLevel = StrToUint(pair[1])
				} else {
					girlId = StrToUint(pair[0])
					girlLevel = StrToUint(pair[1])
				}
			}
		}

		fmt.Println("GirlUnlockCheckNeeds roleLevel:", roleLevel)
		fmt.Println("GirlUnlockCheckNeeds girlId:", girlId)
		fmt.Println("GirlUnlockCheckNeeds girlLevel:", girlLevel)

		if user.Level < int(roleLevel) {
			return false
		}

		for _, part := range user.Girls {
			trimmed := strings.TrimSpace(part)
			pair := strings.Split(trimmed, ":")
			if len(pair) == 2 {
				if StrToUint(pair[0]) == girlId {
					fmt.Println("GirlUnlockCheckNeeds pair[0]:", pair[0])
					if StrToUint(pair[1]) < girlLevel {
						return false
					} else {
						return true
					}
				} else {
					continue
				}
			}
		}

		return false
	}
}

func VehicleUnlockCheckNeeds(roleId uint, user User) (success bool, coin uint64) {
	var need_level int = 0
	var upgrade_cost uint64 = 0
	for _, vehicle := range Vehicles {
		if vehicle.ID == roleId {
			need_level = vehicle.NeedLevel
			upgrade_cost = vehicle.UpgradeCost
			break
		}
	}

	if condition := user.Level < int(need_level); condition {
		return false, 0
	}

	if condition := user.Coins < upgrade_cost; condition {
		return false, 0
	}
	return true, upgrade_cost
}

func CapitalUnlockCheckNeeds(roleId uint, user User) (success bool, coin uint64) {
	var cost uint64 = 0
	for _, capital := range Capitals {
		if capital.ID == roleId {
			cost = capital.Price
			break
		}
	}

	if condition := user.Coins < cost; condition {
		return false, 0
	}
	return true, cost
}

func SellCapital(capitalId uint, user User) (coin uint64, capitals []string) {

	_capitals := []string{}
	coin = 0
	for _, part := range user.Capitals {
		trimmed := strings.TrimSpace(part)
		pair := strings.Split(trimmed, ":")
		if len(pair) == 2 {
			if StrToUint(pair[0]) == capitalId {
				fmt.Println("SellCapital pair[0]:", pair[0])
				var capital domain.Capital = GetCapital(capitalId)
				coin = uint64(capital.Price) + uint64(capital.Price)*uint64(capital.Bonus)*(uint64(time.Now().Unix())-StrToUint64(pair[1]))/(1000*100)
			} else {
				_capitals = append(_capitals, part)
			}
		}
	}
	return coin, _capitals
}

func GetCapital(capitalId uint) (capital domain.Capital) {
	for _, _capital := range Capitals {
		if _capital.ID == capitalId {
			return _capital
		}
	}
	return domain.Capital{}
}

func GetCapitalIncome(user User) (coin uint64, capitals []string) {

	_capitals := []string{}
	coin = 0
	for _, part := range user.Capitals {
		trimmed := strings.TrimSpace(part)
		pair := strings.Split(trimmed, ":")
		if len(pair) == 2 {
			var capitalId uint = StrToUint(pair[0])
			var capital domain.Capital = GetCapital(capitalId)
			coin = coin + uint64(capital.Price) + uint64(capital.Price)*uint64(capital.Bonus)*(uint64(time.Now().Unix())-StrToUint64(pair[1]))/(1000*100)
			_capitals = append(_capitals, fmt.Sprintf("%s:%d", pair[0], time.Now().Unix()))
		}
	}
	return coin, _capitals
}
