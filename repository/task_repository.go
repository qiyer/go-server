package repository

import (
	"context"
	"errors"
	"time"

	"go-server/domain"
	"go-server/redis"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateUserCoins(c context.Context, id primitive.ObjectID, coin uint64) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 创建原子操作管道
	pipeline := []bson.M{
		{
			"$set": bson.M{
				"coins":     bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
				"updatedAt": time.Now(),
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

func CheckIn(c context.Context, id primitive.ObjectID) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	// collection := (*DB).Collection(domain.CollectionUser)

	// // 创建原子操作管道
	// pipeline := []bson.M{
	// 	{
	// 		"$set": bson.M{
	// 			"coins":     bson.M{"$add": bson.A{"$coins", coin}}, // 原子性+3
	// 			"updatedAt": time.Now(),
	// 		},
	// 	},
	// }

	// // 执行findAndModify操作
	// opts := options.FindOneAndUpdate().
	// 	SetReturnDocument(options.After). // 返回更新后的文档
	// 	SetUpsert(false)                  // 禁止自动创建文档

	// var updatedUser domain.User
	// err := collection.FindOneAndUpdate(
	// 	context.TODO(),
	// 	bson.M{"_id": id},
	// 	pipeline,
	// 	opts,
	// ).Decode(&updatedUser)

	updatedUser := domain.User{}
	err := errors.New("等表格设计 和公式设计完成后再来实现")
	return updatedUser, err
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
				"level":     bson.M{"$add": bson.A{"$level", level}}, // 原子性+3
				"coins":     bson.M{"$subtract": bson.A{"$coins", costCoin}},
				"updatedAt": time.Now(),
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

func PassChapter(c context.Context, id primitive.ObjectID, chapter int) (domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	// collection := (*DB).Collection(domain.CollectionUser)

	updatedUser := domain.User{}
	err := errors.New("等表格设计 和公式设计完成后再来实现")
	return updatedUser, err
}

func Ranking(c context.Context) ([]domain.User, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 3. 构建查询选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "level", Value: -1}})                               // 按等级降序排序
	findOptions.SetLimit(10)                                                             // 限制10条结果
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

func UpgradeApartment(c context.Context, id primitive.ObjectID) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 构建过滤条件
	filter := bson.M{"_id": id}
	var user *domain.User
	// 需要查等级，扣除金币
	user, err1 := redis.GetUserFromCache(id)

	if err1 != nil {
		user2, err2 := GetByID(c, id)
		if err2 != nil {
			return err2
		}
		user = &user2
	}
	var level = user.Build.Level
	var coin = user.Coins

	if int(level) >= len(domain.Apartments) {
		return errors.New("小区已满级")
	}

	for _, apartment := range domain.Apartments {
		if apartment.Level == level {
			if coin < apartment.UpgradeCost {
				return errors.New("金币不足")
			}
			coin = coin - apartment.UpgradeCost
			break
		}
	}
	// 定义更新操作（使用 $set 精确更新字段）
	update := bson.M{
		"$set": bson.M{
			"build.level":   bson.M{"$add": bson.A{"$level", 1}},
			"build.updated": time.Now(), // 可添加更新时间戳
			"coins":         coin,
		},
	}

	// 执行更新
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
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
