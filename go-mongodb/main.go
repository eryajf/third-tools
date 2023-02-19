package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bytedance/sonic/encoder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var ctx = context.Background()

func init() {
	// 初始化数据库
	DB = InitDb().Database("eryajf-test")
}

func InitDb() *mongo.Client {
	uri := "mongodb://root:123456@localhost:27017"
	// uri := "mongodb://root:123456@10.6.6.130:27017"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func main() {
	TestMongo()
	defer func() {
		if err := InitDb().Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func TestMongo() {
	{
		// 查询
		// FindOnd() // 查询单条
		// FindOndAsId() //基于ID进行查询
		// FindMany() // 查询多条
		// FindAsLike() //模糊查询
		FindTest()
	}
	{
		// 插入
		// AddOne()  // 插入单条
		// AddMany() // 插入多条
	}
	{
		// 更新
		// UpdateOne()  // 更新单条
		// UpdateMany() // 更新多条
		// UpdateOneField() // 添加单个字段
		// DeleteOneField() // 删除单个字段
		// UpdateManyField() // 添加多个字段
		// DeleteManyField() // 删除多个字段
	}
	{
		// 替换
		// ReplaceOne() // 替换单条
	}
	{
		// 删除
		// DeleteOne()  // 删除单条
		// DeleteMany() // 删除多条
	}
	{
		// 监听
		// Watch() // 监听功能需要集群配置开启副本集功能才可使用
	}
	{
		// 查询
		// FindOnd() // 查询单条
		// FindMany() // 查询多条
		// Count()
	}
	{
		// 聚合查询
		// Aggregate()
	}
	{
		// 自增ID验证
		// AutoIncrement()
	}

}

func AutoIncrement() {
	// coll := DB.Collection("testdata")
	// filter := bson.D{{"age", bson.D{{"$gt", 3}}}}
	// _, err := coll.DeleteMany(context.TODO(), filter)
	// if err != nil {
	// 	panic(err)
	// }

	res, err := InsertOneEx("ttt", bson.M{"name": "aaa", "age": 1})
	if err != nil {
		fmt.Printf("insert err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res.InsertedID)
	}

}

// exec this before run
// db.message_id.insert({id: 1});

// var mgo *mongo.Database

// InsertOneEx .
func InsertOneEx(collection string, document map[string]interface{}) (*mongo.InsertOneResult, error) {
	err := DB.Collection(collection).FindOneAndUpdate(context.Background(), bson.M{}, bson.M{"$inc": bson.M{"_id": 1}},
		&options.FindOneAndUpdateOptions{Projection: bson.M{"_id": 0}}).Decode(&document)
	if err != nil {
		fmt.Println("aaa")
		return nil, err
	}
	return DB.Collection(collection).InsertOne(context.Background(), document)
}

func Aggregate() {
	query := []bson.M{{
		"$lookup": bson.M{
			"from":         "user",
			"localField":   "identify",
			"foreignField": "groupIdentify",
			"as":           "output",
		}},
		{"$match": bson.M{"identify": "yunweizu"}},
	}
	coll := DB.Collection("group")
	cur, err := coll.Aggregate(context.TODO(), query)
	if err != nil {
		fmt.Printf("aggregate failed:%v\n", err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		// 当数据没有映射到结构体时，可以通过map查询
		one := make(map[string]interface{})
		err := cur.Decode(&one)

		if err != nil {
			fmt.Printf("%v\n", err)
		}
		v, err := encoder.Encode(one, encoder.SortMapKeys)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println(string(v))
	}
}

// Count 查询数量
func Count() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"group_identify", "test"}}

	estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
	if estCountErr != nil {
		panic(estCountErr)
	}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Println(estCount, count)
}

func Watch() {
	coll := DB.Collection("testdata")
	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"operationType", "insert"}}}}}
	cs, err := coll.Watch(context.TODO(), pipeline)
	if err != nil {
		fmt.Printf("watch start failed:%v\n", err)
		// panic(err)
	}
	defer cs.Close(context.TODO())

	fmt.Println("Waiting For Change Events. Insert something in MongoDB!")

	for cs.Next(context.TODO()) {
		var event bson.M
		if err := cs.Decode(&event); err != nil {
			panic(err)
		}
		// output, err := json.MarshalIndent(event["fullDocument"], "", "    ")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("%s\n", output)
		v, err := encoder.Encode(event, encoder.SortMapKeys)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println(string(v))
	}
	if err := cs.Err(); err != nil {
		panic(err)
	}
}

// DeleteMany 删除多条
func DeleteMany() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"age", bson.D{{"$gt", 3}}}}
	_, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}

// DeleteOne 删除单条
func DeleteOne() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"identify", "aaa"}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}

// ReplaceOne 替换单条
func ReplaceOne() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"identify", "aaa"}}
	replacement := bson.D{{"name", "小A-r"}, {"identify", "aaa"}, {"content", "this is aaa replace"}}

	_, err := coll.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		panic(err)
	}
}

// UpdateMany 更新多条
func UpdateMany() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"group_identify", "test"}}
	update := bson.D{{"$mul", bson.D{{"age", 3}}}}

	_, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}

// UpdateManyField 添加多个字段数据
func UpdateManyField() {
	coll := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("62159551120b25bd2c801b09")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	filter := bson.M{"_id": objid}
	linkData := []map[string]string{
		{
			"field_identify": "eryajf_guanliandd",
			"model_data_id":  "6215aaf220ea934fb727096c",
		},
		{
			"field_identify": "eryajf_guanliandd",
			"model_data_id":  "6215aaf220ea934fbaaaaaaa",
		},
	}
	updata := bson.M{"$push": bson.M{"link_data": bson.M{"$each": linkData, "$position": 0}}} // $position: 定义插入的位置，0表示插入到最前面，-1 表示数组中最后一个元素之前的位置
	_, err = coll.UpdateOne(ctx, filter, updata)
	if err != nil {
		panic(err)
	}
}

// DeleteManyField 删除多条记录
func DeleteManyField() {
	coll := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("62159551120b25bd2c801b09")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	filter := bson.M{"_id": objid}
	linkData := []map[string]string{
		{
			"field_identify": "eryajf_guanliandd",
			"model_data_id":  "6215aaf220ea934fb727096c",
		},
		{
			"field_identify": "eryajf_guanliandd",
			"model_data_id":  "6215aaf220ea934fbaaaaaaa",
		},
	}
	updata := bson.M{"$pullAll": bson.M{"link_data_testa": linkData}}
	_, err = coll.UpdateOne(ctx, filter, updata)
	if err != nil {
		panic(err)
	}
}

// UpdateOneField 添加单个字段
func UpdateOneField() {
	coll := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("62159551120b25bd2c801b09")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	filter := bson.M{"_id": objid}

	updata := bson.M{"$push": bson.M{"link_data": bson.M{"field_identify": "1", "model_data_id": "5"}}}

	_, err = coll.UpdateOne(ctx, filter, updata)
	if err != nil {
		panic(err)
	}
}

// UpdateOneField 添加单个字段
func UpdateOneFieldbak() {
	coll := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("62159551120b25bd2c801b09")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	filter := bson.M{"_id": objid}

	updata := bson.M{"$push": bson.M{"link_data": bson.M{"field_identify": "1", "model_data_id": "5"}}}

	_, err = coll.UpdateOne(ctx, filter, updata)
	if err != nil {
		panic(err)
	}
}

// DeleteOneField 删除单个字段
func DeleteOneField() {
	coll := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("62159551120b25bd2c801b09")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	filter := bson.M{"_id": objid}

	updata := bson.M{"$pull": bson.M{"link_data": bson.M{"field_identify": "1", "model_data_id": "5"}}}

	_, err = coll.UpdateOne(ctx, filter, updata)
	if err != nil {
		panic(err)
	}
}

// UpdateOne 更新单条
func UpdateOne() {
	coll := DB.Collection("testdata")
	filter := bson.D{{"title", "test"}}
	update := bson.D{{"$set", bson.D{{"avg_rating", 4.5}}}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}

// AddOne 插入单条
func AddOne() {
	coll := DB.Collection("testdata")
	doc := make(map[string]interface{})
	doc["title"] = "test"
	doc["content"] = "this is a test"
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}
}

// AddMany 插入多条
func AddMany() {
	coll := DB.Collection("testdata")
	docs := []interface{}{
		bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}},
		bson.D{{"title", "Showcasing a Blossoming Binary"}, {"text", "Binary data, safely stored with GridFS. Bucket the data"}},
	}
	_, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
}

type Group struct {
	Name     string   `json:"name" bson:"name"`
	NickName string   `json:"nick_name" bson:"nick_name"`
	Users    []string `json:"users" bson:"users"`
}

type User struct {
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}

type Data struct {
	DatasIdentify string `json:"datas_identify"`
	Data          []struct {
		Name       string `json:"name"`
		Identify   string `json:"identify"`
		CreateTime string `json:"create_time"`
	} `json:"data"`
}

func FindTest() {
	// var group Group
	// table := DB.Collection("groups")
	// res := table.FindOne(ctx, bson.M{"name": "ops"})
	// if err := res.Err(); err != nil {
	// 	fmt.Printf("find data failed: %v\n", err)
	// }
	// if err := res.Decode(&group); err != nil {
	// 	fmt.Printf("decode data failed: %v\n", err)
	// }

	// var alreadyLinks []primitive.ObjectID
	// for _, v := range group.Users {
	// 	objid, err := primitive.ObjectIDFromHex(v)
	// 	if err != nil {
	// 		fmt.Printf("%v\n", err)
	// 	}
	// 	alreadyLinks = append(alreadyLinks, objid)
	// }
	// if len(alreadyLinks) == 0 {
	// 	alreadyLinks = append(alreadyLinks, primitive.NilObjectID)
	// }

	// filter := bson.D{}
	// filter = append(filter, bson.E{Key: "_id", Value: bson.M{"$nin": alreadyLinks}})
	// users, err := ListUser(filter, options.FindOptions{})
	// if err != nil {
	// 	fmt.Printf("get data failed: %v\n", err)
	// }
	// for _, v := range users {
	// 	fmt.Printf("用户名: %v 年龄: %v\n", v.Name, v.Age)
	// }

	filters := bson.D{}
	filter := bson.E{Key: "datas_identify", Value: "eryajf"}
	searchFilter := bson.E{Key: "$text", Value: bson.M{"$search": "2022"}}
	filters = append(filters, filter, searchFilter)
	datas, err := ListData(filters, options.FindOptions{})
	if err != nil {
		fmt.Printf("get data failed: %v\n", err)
	}
	for _, v := range datas {
		fmt.Println(v)
	}

}

// ListData 获取数据列表
func ListData(filter bson.D, options options.FindOptions) ([]*Data, error) {
	table := DB.Collection("datas")
	cus, err := table.Find(ctx, filter, &options)
	if err != nil {
		fmt.Printf("find data failed: %v\n", err)
	}
	defer func(cus *mongo.Cursor, ctx context.Context) {
		err := cus.Close(ctx)
		if err != nil {
			return
		}
	}(cus, ctx)

	list := make([]*Data, 0)
	for cus.Next(ctx) {
		data := new(Data)
		err := cus.Decode(&data)
		if err != nil {
			fmt.Printf("decode data failed: %v\n", err)
		}
		list = append(list, data)
	}

	return list, nil
}

// FindOne 查询单条
func FindOndBase() {
	// table := DB.Collection("model_data")

	var result bson.M
	table := DB.Collection("testdata")
	err := table.FindOne(context.TODO(), bson.M{"link_data.model_data_id": "6215aaf220ea934fb727096c"}).Decode(&result)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// fmt.Printf("%v\n", result)
	v, err := encoder.Encode(result, encoder.SortMapKeys)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(string(v))

}

// ModelData 模型数据
type ModelData struct {
	ID            string            `json:"id,omitempty" bson:"_id,omitempty"`
	ModelIdentify string            `json:"modelIdentify" bson:"model_identify"`           // 模型唯一标识
	Data          map[string]string `json:"data" bson:"data"`                              // 数据列表
	IsDel         bool              `json:"isDel,omitempty" bson:"is_del"`                 // 是否删除
	CreateAt      int64             `json:"createAt,omitempty" bson:"create_at,omitempty"` // 创建时间
	ModifyAt      int64             `json:"modifyAt,omitempty" bson:"modify_at,omitempty"` // 更新时间
}

// // FindOne 查询单条
// func FindOnd() {
// 	table := DB.Collection("movies")
// 	// err := table.FindOne(context.TODO(), bson.M{"title": "aa"}).Decode(&result)
// 	// if err != nil {
// 	// 	fmt.Printf("Error: %v\n", err)
// 	// }
// 	// if _, ok := result["testa"]; ok {
// 	// 	fmt.Println("存在")
// 	// } else {
// 	// 	fmt.Println("不存在")
// 	// }

// 	var result bson.M

// 	var findoptions options.FindOptions
// 	findoptions.SetLimit(1)
// 	findoptions.SetSort(bson.D{{"_id", -1}})
// 	a, err := table.Find(ctx, bson.M{}, &findoptions)
// 	if err != nil {
// 		fmt.Printf("find failed:%v\n", err)
// 	}
// 	err = a.Decode(&result)
// 	if err != nil {
// 		fmt.Printf("err")
// 	}
// 	if _, ok := result["testa"]; ok {
// 		fmt.Println("存在")
// 	} else {
// 		fmt.Println("不存在")
// 	}
// }

type Article struct {
	Title string
	Text  string
}

// FindAsLike 模糊查询
func FindAsLike() {
	table := DB.Collection("testdata")
	findOptions := options.Find()

	filter := bson.D{}
	filter = append(filter, bson.E{
		Key:   "title",
		Value: bson.M{"$regex": primitive.Regex{Pattern: ".*" + "a" + ".*", Options: "i"}}}) //i 表示不区分大小写

	cus, err := table.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Printf("find failed: %v\n", err)
	}

	defer func(cus *mongo.Cursor, ctx context.Context) {
		err := cus.Close(ctx)
		if err != nil {
			return
		}
	}(cus, ctx)

	list := make([]*Article, 0)
	for cus.Next(ctx) {
		article := new(Article)
		err := cus.Decode(&article)
		if err != nil {
			fmt.Printf("decode failed: %v\n", err)
			// return nil, tools.NewMongoError(err)
		}
		list = append(list, article)
	}

	fmt.Println("results: ", list)
	for _, v := range list {
		fmt.Println(v)
	}
}

// FindOndAsId 基于ID查询单条
func FindOndAsId() {
	// var result bson.M
	table := DB.Collection("testdata")

	objid, err := primitive.ObjectIDFromHex("61f24a4dc8d32bc297dbed02")
	if err != nil {
		fmt.Printf("obj id failed: %v\n", err)
	}

	filter := bson.M{"_id": objid}

	res := table.FindOne(ctx, filter)
	fmt.Println(errors.Is(mongo.ErrNoDocuments, res.Err()))
	// err = table.FindOne(context.TODO(), filter).Decode(&result)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }
	// fmt.Printf("%v\n", result)
	// v, err := encoder.Encode(result, encoder.SortMapKeys)
	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// }
	// fmt.Println(string(v))
}

// FindMany 查询多条
func FindMany() {
	/*
		插入测试数据
		db.testdata.insert([
			{ "name" : "小A", "identify": "aaa","age":1},
			{ "name" : "小B", "identify": "bbb","age":2},
			{ "name" : "小C", "identify": "ccc","age":3},
			{ "name" : "小D", "identify": "ddd","age":4},
			{ "name" : "小E", "identify": "eee","age":5},
			{ "name" : "小F", "identify": "fff","age":6},
		])
	*/
	table := DB.Collection("testdata")
	filter := bson.D{{"age", bson.D{{"$lte", 3}}}} // 查询年龄小于等于3的，这里特别有意思，能够使用$lte这种方法，类似这样的，MongoDB还提供了很多其他的查询方法，比如$gt等等
	cursor, err := table.Find(ctx, filter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		v, err := encoder.Encode(result, encoder.SortMapKeys)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println(string(v))
	}

}
