package main

import (
	"context"
	"fmt"
	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	client *redis.Client
)

//连接redis服务端
func goRedisConnect() {
	client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0, //redis默认有0-15共16个数据库，这里设置操作索引为0的数据库
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()

	if err != nil {

		log.Fatal(err)
	}

	if pong != "PONG" {

		log.Fatal("客户端连接redis服务端失败")
	} else {

		fmt.Println("客户端已成功连接至redis服务端")
	}
}

//string类型数据操作
//redis命令：set key val
func set(key, val string) {
	var ctx = context.Background()
	//有效期为0表示不设置有效期，非0表示经过该时间后键值对失效
	result, err := client.Set(ctx, key, val, 0).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(result)
}

//redis命令：get key
func get(key string) {
	var ctx = context.Background()
	val, err := client.Get(ctx, key).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：mset key1 val1 key2 val2 key3 val3 ...
func mset(key1, val1, key2, val2, key3, val3 string) {
	var ctx = context.Background()
	//以下三种方式都可以，习惯于对象操作的我优先选择第三种
	//result,err := client.MSet(key1,val1,key2,val2,key3,val3).Result()
	//result,err := client.MSet([]string{key1,val1,key2,val2,key3,val3}).Result()
	result, err := client.MSet(ctx, map[string]interface {
	}{
		key1: val1, key2: val2, key3: val3}).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(result)
}

//redis命令：mget key1 key2 key3 ...
func mget(key1, key2, key3 string) {
	var ctx = context.Background()
	vals, err := client.MGet(ctx, key1, key2, key3).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：del key1 key2 key3 ...
func del(key1, key2, key3 string) {
	var ctx = context.Background()
	result, err := client.Del(ctx, key1, key2, key3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(result)
}

//redis命令：getrange key start end
func getrange(key string, start, end int64) {
	var ctx = context.Background()
	val, err := client.GetRange(ctx, key, start, end).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：strlen key
func strlen(key string) {

	len, err := client.StrLen(context.Background(), key).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//redis命令：setex key time val
func setex(key, val string, expire int) {

	//time.Duration其实也是int64，不过是int64的别名罢了，但这里如果expire使用int64也无法与time.Second运算，
	//因为int64和Duration虽然本质一样，但表面上属于不同类型，go语言中不同类型是无法做运算的
	result, err := client.Set(context.Background(), key, val, time.Duration(expire)*time.Second).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(result)
}

//redis命令：append key val
func append(key, val string) {

	//将val插入key对应值的末尾，并返回新串长度
	len, err := client.Append(context.Background(), key, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//redis命令：exists key
func exists(key string) {

	//返回1表示存在，0表示不存在
	isExists, err := client.Exists(context.Background(), key).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(isExists)
}

//hash类型数据操作
//redis命令：hset hashTable key val
func hset(hashTable, key, val string) {

	isSetSuccessful, err := client.HSet(context.Background(), hashTable, key, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	//如果键存在这返回false，如果键不存在则返回true
	fmt.Println(isSetSuccessful)
}

//redis命令：hget hashTable key
func hget(hashTable, key string) {

	val, err := client.HGet(context.Background(), hashTable, key).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：hmset hashTable key1 val1 key2 val2 key3 val3 ...
//该函数本身有问题，只插入一个键值对的话相当于hset，可以成功
//如果插入一个以上的键值对则会报错：ERR wrong number of arguments for 'hset' command
//且go-redis官方本身也不推荐是用该函数
//func hmset(hashTable,key1,val1,key2,val2,key3,val3 string){

//	_,err := client.HMSet(hashTable,key1,val1,key2,val2,key3,val3).Result()
//
//	if err != nil {

//		log.Fatal(err)
//	}
//}
//redis命令：hmget hashTable key1 key2 key3 ...
func hmget(hashTable, key1, key2, key3 string) {

	vals, err := client.HMGet(context.Background(), hashTable, key1, key2, key3).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：hdel hashTable key1 key2 key3 ...
func hdel(hashTable, key1, key2, key3 string) {

	//返回1表示删除成功，返回0表示删除失败
	//只要至少有一个被删除则返回1（不存在的键不管），一个都没删除则返回0（不存在的则也算没删除）
	n, err := client.Del(context.Background(), hashTable, key1, key2, key3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：hgetall hashTable
func hgetall(hashTable string) {

	vals, err := client.HGetAll(context.Background(), hashTable).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：hexists hashTable key
func hexists(hashTable, key string) {

	isExists, err := client.HExists(context.Background(), hashTable, key).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(isExists)
}

//redis命令：hlen hashTable
func hlen(hashTable string) {

	len, err := client.HLen(context.Background(), hashTable).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//redis命令：hkeys hashTable
func hkeys(hashTable string) {

	keys, err := client.HKeys(context.Background(), hashTable).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range keys {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：hvals hashTable
func hvals(hashTable string) {

	vals, err := client.HVals(context.Background(), hashTable).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

func main() {
	goRedisConnect()
	set("key1", "value1")
	get("key1")
}
