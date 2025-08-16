package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"go-server/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserByCacheOrDB(c context.Context, id primitive.ObjectID) (domain.User, error) {
	user, err := GetUserCache(id.Hex())
	if err != nil {
		user, err = GetByID(c, id)
		if err != nil {
			return domain.User{}, err
		}
	}

	return user, nil
}

func UpdateUserCoinsWithTime(c context.Context, id primitive.ObjectID, coin uint64, online int, bnousTime int64, autoClickerTime int64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":                bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
				"updatedAt":            time.Now(),
				"onlineTime":           online,
				"timesBonusTimeStamp":  bnousTime,
				"autoClickerTimeStamp": autoClickerTime,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateUserCoinsWithClick(c context.Context, id primitive.ObjectID, coin uint64, clickTime int64, num int, clicker int) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":              bson.M{"$add": bson.A{"$coins", coin}},
				"lastClickTimeStamp": clickTime,
				"boxNum":             num,
				"boxClicker":         clicker,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateUserCoins(c context.Context, id primitive.ObjectID, coin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins": bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateUserCoinsByBox(c context.Context, id primitive.ObjectID, coin uint64, num int) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":  bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
				"boxNum": bson.M{"$subtract": bson.A{"$boxNum", num}},
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateOnlineRewards(c context.Context, id primitive.ObjectID, rewards []int) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 构建过滤条件
	filter := bson.M{"_id": id}

	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"onlineRewards": rewards,
			},
		},
	}

	// 执行更新
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func ClaimOnlineRewards(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	// collection := (*DB).Collection(domain.CollectionUser)

	updatedUser := domain.User{}
	err := errors.New("等表格设计 和公式设计完成后再来实现")
	return updatedUser, err
}

func LevelUp(c context.Context, id primitive.ObjectID, level int, costCoin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"level": bson.M{"$add": bson.A{"$level", level}}, // 原子性+3
				"coins": bson.M{"$subtract": bson.A{"$coins", costCoin}},
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func RoleLevelUp(c context.Context, id primitive.ObjectID, girls []string, costCoin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("RoleLevelUp costCoin：%+v\n", costCoin)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins": bson.M{"$subtract": bson.A{"$coins", costCoin}},
				"girls": girls,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateBoxNum(c context.Context, id primitive.ObjectID, num int, clicker int) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("UpdateBoxNum num%+v\n", num)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"boxNum":     num,
				"boxClicker": clicker,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UnLockVehicle(c context.Context, id primitive.ObjectID, vehicles string, costCoin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("RoleLevelUp costCoin：%+v\n", costCoin)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":    bson.M{"$subtract": bson.A{"$coins", costCoin}},
				"vehicles": vehicles,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UnLockCapital(c context.Context, id primitive.ObjectID, capitals []string, costCoin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("UnLockCapital costCoin：%+v\n", costCoin)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":    bson.M{"$subtract": bson.A{"$coins", costCoin}},
				"capitals": capitals,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func SellCapital(c context.Context, id primitive.ObjectID, capitals []string, coin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("SellCapital Coin：%+v\n", coin)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":    bson.M{"$add": bson.A{"$coins", coin}},
				"capitals": capitals,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func PassChapter(c context.Context, id primitive.ObjectID, chapter int) (domain.User, int, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// user, _ := GetUserByCacheOrDB(c, id)
	// if (user.Chapter + 1) != chapter {
	// 	return domain.User{}, 0, errors.New("章节不连续")
	// }

	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"chapter": chapter,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, chapter, err
}

func Ranking(c context.Context) ([]domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 3. 构建查询选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "level", Value: -1}})                               // 按等级降序排序
	findOptions.SetLimit(100)                                                            // 限制10条结果
	findOptions.SetProjection(bson.D{{Key: "name", Value: 1}, {Key: "level", Value: 1}}) // 排除_id字段 {Key: "_id", Value: 0},

	// 4. 执行查询
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	// 5. 处理结果
	var results []domain.User
	err = cur.All(context.TODO(), &results)
	return results, err
}

func VehicleRanking(c context.Context) ([]domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 3. 构建查询选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "vehicle.level", Value: -1}})                                                   // 按等级降序排序
	findOptions.SetLimit(50)                                                                                         // 限制10条结果
	findOptions.SetProjection(bson.D{{Key: "name", Value: 1}, {Key: "level", Value: 1}, {Key: "vehicle", Value: 1}}) // 排除_id字段 {Key: "_id", Value: 0},

	// 4. 执行查询
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	// 5. 处理结果
	var results []domain.User
	err = cur.All(context.TODO(), &results)
	return results, err
}

func UpgradeApartment(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.Build.Level
	var coin = user.Coins

	if int(level) >= len(domain.Apartments) {
		return domain.User{}, errors.New("小区已满级")
	}

	for _, apartment := range domain.Apartments {
		if apartment.Level == level {
			if coin < apartment.UpgradeCost {
				return domain.User{}, errors.New("金币不足")
			}
			coin = coin - apartment.UpgradeCost
			break
		}
	}
	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"build.level": bson.M{"$add": bson.A{"$build.level", 1}},
				"coins":       coin,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpgradeVehicle(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.Vehicle.Level
	var coin = user.Coins

	if int(level) >= len(domain.Vehicles) {
		return domain.User{}, errors.New("坐骑已满级")
	}

	for _, vehicle := range domain.Vehicles {
		if vehicle.ID == level {
			if coin < vehicle.UpgradeCost {
				return domain.User{}, errors.New("您的金币不足")
			}
			if vehicle.NeedLevel > user.Level {
				return domain.User{}, errors.New("您的主角色等级不足")
			}
			coin = coin - vehicle.UpgradeCost
			break
		}
	}
	// 定义更新操作（使用 $set 精确更新字段）
	update := []bson.M{
		{
			"$set": bson.M{
				"vehicle.level": bson.M{"$add": bson.A{"$vehicle.level", 1}},
				"coins":         coin,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func ChangeVehicleVehicle(c context.Context, id primitive.ObjectID, displayLevel int) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	if displayLevel >= len(domain.Vehicles) {
		return domain.User{}, errors.New("小区已满级")
	}

	if displayLevel < 1 {
		return domain.User{}, errors.New("等级有误")
	}

	// 定义更新操作（使用 $set 精确更新字段）
	update := []bson.M{
		{
			"$set": bson.M{
				"vehicle.displayLevel": displayLevel,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func CaiShen(c context.Context, id primitive.ObjectID) (domain.User, uint64, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 需要查等级，扣除金币
	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.Level
	var addCoin uint64 = 0
	//需要check广告是否播放成功
	if level < 97 {
		addCoin = uint64((domain.BaseCaiShen + level) * 10000)
	} else if level < 201 {
		var base = (domain.BaseCaiShen + 97) * 10000
		var bcoin = float64(base) * math.Pow(domain.CaiShenGrowth1, float64(level-96))
		addCoin = uint64(bcoin)
	} else {
		var base = 3119753
		var bcoin = float64(base) * math.Pow(domain.CaiShenGrowth2, float64(level-200))
		addCoin = uint64(bcoin)
	}
	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"build.updated": time.Now(), // 可添加更新时间戳
				"coins":         bson.M{"$add": bson.A{"$coins", addCoin}},
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, addCoin, err
}

func Training(c context.Context, id primitive.ObjectID) (domain.User, int, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 需要查等级，扣除金币
	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.GirlTrainLevel
	//需要check广告是否播放成功
	if level >= 13 {
		return user, 0, errors.New("秘书培训等级已满")
	}
	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"grilTrainLevel": bson.M{"$add": bson.A{"$grilTrainLevel", 1}},
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, updatedUser.GirlTrainLevel, err
}

func QuickEarn(c context.Context, id primitive.ObjectID) (domain.User, uint64, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.QuickEarn
	var addCoin uint64 = 0

	for _, earn := range domain.QuickEarns {
		if int(earn.Level) == level {
			addCoin = earn.Bonus
			break
		}
	}

	var coin = user.Coins + addCoin
	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"build.updated": time.Now(), // 可添加更新时间戳
				"coins":         coin,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, addCoin, err
}

func ContinuousClick(c context.Context, id primitive.ObjectID, add int) (domain.User, int, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	user, _ := GetUserByCacheOrDB(c, id)

	var level = user.ContinuousClick + add
	if level == 3 {
		level = 2
	}

	if level > 40 {
		level = 40
	}

	// 定义更新操作（使用 $set 精确更新字段）

	update := []bson.M{
		{
			"$set": bson.M{
				"continuousClick": level,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		update,
		opts,
	).Decode(&updatedUser)

	return updatedUser, level, err
}

func TimesBonus(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	user, _ := GetUserByCacheOrDB(c, id)
	var level = user.TimesBonusRatio
	var bonusTime = user.TimesBonusTimeStamp
	timestamp := time.Now().Unix()
	if user.TimesBonusTimeStamp > timestamp {
		bonusTime = bonusTime + domain.TimesBonusBaseTime
	} else {
		bonusTime = timestamp + domain.TimesBonusBaseTime
		//如果用户的临时倍数奖励小于2，则随机生成一个1-8的整数作为倍数奖励
		if user.TimesBonusRatio < 2 {
			src := rand.NewSource(time.Now().UnixNano())
			r := rand.New(src)
			// 生成 1-8 的随机整数
			randomNumber := r.Intn(99) + 1
			level = 2
			if randomNumber < 15 {
				level = 2
			} else if randomNumber < 50 {
				level = 3
			} else if randomNumber < 80 {
				level = 4
			} else if randomNumber < 95 {
				level = 5
			} else {
				level = 10
			}
		}
	}
	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"timesBonus":          level,
				"timesBonusTimeStamp": bonusTime,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func AutoClickerTime(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 构建过滤条件
	user, _ := GetUserByCacheOrDB(c, id)
	var autoClickerTime = user.AutoClickerTimeStamp
	timestamp := time.Now().Unix()
	if user.AutoClickerTimeStamp > timestamp {
		autoClickerTime = autoClickerTime + domain.TimesBonusBaseTime
	} else {
		autoClickerTime = timestamp + domain.TimesBonusBaseTime
	}

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"autoClickerTimeStamp": autoClickerTime,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateBossInfo(c context.Context, id primitive.ObjectID, bosses []string, coin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	fmt.Printf("UpdateBossInfo Coin：%+v\n", coin)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":  bson.M{"$add": bson.A{"$coins", coin}},
				"bosses": bosses,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func CreateTask(c context.Context, task *domain.Task) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionTask)

	_, err := collection.InsertOne(c, task)

	return err
}

func FetchTaskByUserID(c context.Context, userID string) ([]domain.Task, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionTask)

	var tasks []domain.Task

	idHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return tasks, err
	}

	cursor, err := collection.Find(c, bson.M{"userID": idHex})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &tasks)
	if tasks == nil {
		return []domain.Task{}, err
	}

	return tasks, err
}

func UpdateRankingLikeRecord(c context.Context, id primitive.ObjectID, coin uint64, rankingLikeRecord []string) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":             bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
				"rankingLikeRecord": rankingLikeRecord,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}

func UpdateRankingRewards(c context.Context, id primitive.ObjectID, coin uint64, rankingRecord []string) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":         bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
				"rankingRecord": rankingRecord,
			},
		},
	}

	// 执行findAndModify操作
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After). // 返回更新后的文档
		SetUpsert(false)                  // 禁止自动创建文档

	var updatedUser domain.User
	err := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		pipeline,
		opts,
	).Decode(&updatedUser)

	return updatedUser, err
}
