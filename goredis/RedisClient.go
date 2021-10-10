package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
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

//list类型数据操作

//list类型数据操作
//redis命令：lpush mylist val1 val2 val3 ...
func lpush(ctx context.Context, mylist, val1, val2, val3 string) {

	//返回列表的总长度（即有多少个元素在列表中）
	n, err := client.LPush(ctx, mylist, val1, val2, val3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：rpush mylist val1 val2 val3 ...
func rpush(ctx context.Context, mylist, val1, val2, val3 string) {

	//返回列表的总长度（即有多少个元素在列表中）
	n, err := client.RPush(ctx, mylist, val1, val2, val3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：lpop mylist
func lpop(ctx context.Context, mylist string) {

	//返回被删除的值
	val, err := client.LPop(ctx, mylist).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：rpop mylist
func rpop(ctx context.Context, mylist string) {

	//返回被删除的值
	val, err := client.RPop(ctx, mylist).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：lrem mylist count val
func lrem(ctx context.Context, mylist, val string, count int64) {

	//返回成功删除的val的数量
	n, err := client.LRem(ctx, mylist, count, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：ltrim mylist start end
func ltrim(ctx context.Context, mylist string, start, end int64) {

	//返回状态（OK）
	status, err := client.LTrim(ctx, mylist, start, end).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(status)
}

//redis命令：lset mylist index val
func lset(ctx context.Context, mylist, val string, index int64) {

	status, err := client.LSet(ctx, mylist, index, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(status)
}

//redis命令：lindex mylist index
func lindex(ctx context.Context, mylist string, index int64) {

	//通过索引查找字符串
	val, err := client.LIndex(ctx, mylist, index).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：lrange mylist start end
func lrange(ctx context.Context, mylist string, start, end int64) {

	vals, err := client.LRange(ctx, mylist, start, end).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：llen mylist
func llen(ctx context.Context, mylist string) {

	len, err := client.LLen(ctx, mylist).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//无序集合set类型数据操作
//redis命令：sadd myset val1 val2 val3 ...
func sadd(ctx context.Context, myset, val1, val2, val3 string) {

	n, err := client.SAdd(ctx, myset, val1, val2, val3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：srem myset val
func srem(ctx context.Context, myset, val string) {

	//删除集合中的值并返回其索引
	index, err := client.SRem(ctx, myset, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(index)
}

//redis命令：spop myset
func spop(ctx context.Context, myset string) {

	//随机删除一个值并返回
	val, err := client.SPop(ctx, myset).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：smembers myset
func smembers(ctx context.Context, myset string) {

	vals, err := client.SMembers(ctx, myset).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：scard myset
func scard(ctx context.Context, myset string) {

	len, err := client.SCard(ctx, myset).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//redis命令：sismember myset val
func sismember(ctx context.Context, myset, val string) {

	//判断值是否为集合中的成员
	isMember, err := client.SIsMember(ctx, myset, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(isMember)
}

//redis命令：srandmember myset count
func srandmembers(ctx context.Context, myset string, count int64) {

	vals, err := client.SRandMemberN(ctx, myset, count).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//该函数是上一个函数在只随机取一个元素的情况
func srandmember(ctx context.Context, myset string) {

	val, err := client.SRandMember(ctx, myset).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(val)
}

//redis命令：smove myset myset2 val
func smove(ctx context.Context, myset, myset2, val string) {

	isSuccessful, err := client.SMove(ctx, myset, myset2, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(isSuccessful)
}

//redis命令：sunion myset myset2 ...
func sunion(ctx context.Context, myset, myset2 string) {

	vals, err := client.SUnion(ctx, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：sunionstore desset myset myset2 ...
func sunionstore(ctx context.Context, desset, myset, myset2 string) {

	//返回新集合的长度
	n, err := client.SUnionStore(ctx, desset, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：sinter myset myset2 ...
func sinter(ctx context.Context, myset, myset2 string) {

	vals, err := client.SInter(ctx, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：sinterstore desset myset myset2 ...
func sinterstore(ctx context.Context, desset, myset, myset2 string) {

	n, err := client.SInterStore(ctx, desset, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：sdiff myset myset2 ...
func sdiff(ctx context.Context, myset, myset2 string) {

	vals, err := client.SDiff(ctx, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}

	for k, v := range vals {

		fmt.Printf("k = %v v = %s\n", k, v)
	}
}

//redis命令：sdiffstore desset myset myset2 ...
func sdiffstore(ctx context.Context, desset, myset, myset2 string) {

	n, err := client.SDiffStore(ctx, desset, myset, myset2).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//有序集合zset类型数据操作
//redis命令：zadd myzset score1 val1 score2 val2 score3 val3 ...
func zadd(ctx context.Context, myzset, val1, val2, val3 string, score1, score2, score3 float64) {

	member1 := &redis.Z{
		score1, val1}
	member2 := &redis.Z{
		score2, val2}
	member3 := &redis.Z{
		score3, val3}

	n, err := client.ZAdd(ctx, myzset, member1, member2, member3).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：zrem myzset val1 val2 ...
func zrem(ctx context.Context, myzset, val1, val2 string) {

	n, err := client.ZRem(ctx, myzset, val1, val2).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：srange myzset start end [withscores]
func zrange(ctx context.Context, myzset string, start, end, flag int64) {

	if flag == 0 {

		//不加withscores
		vals, err := client.ZRange(ctx, myzset, start, end).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range vals {

			fmt.Printf("k = %v v = %s\n", k, v)
		}
	} else if flag == 1 {

		//加withscores
		svals, err := client.ZRangeWithScores(ctx, myzset, start, end).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range svals {

			fmt.Printf("k = %v v = %s s = %.2f\n", k, v.Member, v.Score)
		}
	}
}

//redis命令：srevrange myzset start end [withscores]
func zrevrange(ctx context.Context, myzset string, start, end, flag int64) {

	if flag == 0 {

		//不加withscores
		vals, err := client.ZRevRange(ctx, myzset, start, end).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range vals {

			fmt.Printf("k = %v v = %s\n", k, v)
		}
	} else if flag == 1 {

		//加withscores
		svals, err := client.ZRevRangeWithScores(ctx, myzset, start, end).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range svals {

			fmt.Printf("k = %v v = %s s = %.2f\n", k, v.Member, v.Score)
		}
	}
}

//redis命令：zrangebyscore myzset start end [withscores]
func zrangebyscore(ctx context.Context, myzset, start, end string, flag int) {

	if flag == 0 {

		//不加withscores
		vals, err := client.ZRangeByScore(ctx, myzset, &redis.ZRangeBy{
			Min: start, Max: end, Count: 0}).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range vals {

			fmt.Printf("k = %v v = %s\n", k, v)
		}
	} else if flag == 1 {

		//加withscores
		svals, err := client.ZRangeByScoreWithScores(ctx, myzset, &redis.ZRangeBy{
			Min: start, Max: end, Count: 0}).Result()

		if err != nil {

			log.Fatal(err)
		}

		for k, v := range svals {

			fmt.Printf("k = %v v = %s s = %.2f\n", k, v.Member, v.Score)
		}
	}
}

//redis命令：zcard myzset
func zcard(ctx context.Context, myzset string) {

	len, err := client.ZCard(ctx, myzset).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(len)
}

//redis命令：zcount myzset minscore maxscore
func zcount(ctx context.Context, myzset, minscore, maxscore string) {

	n, err := client.ZCount(ctx, myzset, minscore, maxscore).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(n)
}

//redis命令：zrank myzset val
func zrank(ctx context.Context, myzset, val string) {

	index, err := client.ZRank(ctx, myzset, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(index)
}

//redis命令：zscore myzset val
func zscore(ctx context.Context, myzset, val string) {

	score, err := client.ZScore(ctx, myzset, val).Result()

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(score)
}

func main() {
	goRedisConnect()
	//set("key1", "value1")
	//get("key1")

	//hset("hs","k1","v1")
	//hset("hs","k2","v2")
	//hset("hs","k3","v3")
	//
	//
	//hget("hs","k2")
	//hkeys("hs")
	//hvals("hs")
	//
	//hlen("hs")
	var ctx context.Context
	ctx = context.Background()

	//lset(ctx,"list","value-1",1)
	//lindex(ctx,"list",0)
	lrange(ctx, "list", 0, 1)
}
