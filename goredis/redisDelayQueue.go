package main

/*
  redis延时队列

redis的zset是一种能自动排序的数据结构，我们可以用这个特性来实现简单的延时队列。

利用zadd将数据添加到zset中，每个数据的score值设置为数据的延时时间+当前时间戳，后台goroutine不断zrange轮询zset，取出score值小于当前时间戳的数据，然后再对数据进一步处理，这样就实现了简单延时队列的功能。

Zadd
zadd支持批量添加的功能，有需要可以自己去探索下，这里我们先提供添加一个元素的方法

import "github.com/garyburd/redigo/redis"

func AddZset(data, zsetName string, score int64) (err error) {
	con := pool.Get()
	defer con.Close()
	_, err = con.Do("zadd", zsetName, score, data)
	if err != nil {
		return
	}
	return
}
Zrange
zrange按照score值从小到大遍历zset，提供start和end两个下标参数，都以0为底。下标超出范围不会出错。start=0，end=1表示遍历前2个元素，end=-1表示最后一个元素。

可以通过使用 WITHSCORES 选项，来让成员和它的 score 值一并返回，返回列表以 value1,score1, ..., valueN,scoreN 的格式表示。

我们可以每次取几个数据，判断当前时间戳与数据对应score的大小关系，并决定是否处理这些数据。

//index sorted set from start to end, [start:end], eg: [0:1] will return[member1, score1, member2, score2]
func RangeZset(start, end int, zsetName string) (data []string, err error) {
	con := pool.Get()
	defer con.Close()
	data, err = redis.Strings(con.Do("zrange", zsetName, start, end, "withscores"))
	return
}
ZrangeWithScore
zrangewithscore提供了min和max两个参数，它会返回zset中介于min和max（包含两者）之间的所有元素。

我们可以以0为底，以当前时间戳为max，取出zset中所有到期的数据，然后进行处理。

func RangeZsetByScore(start, end int64, zsetName string) (data []string, err error) {
	con := pool.Get()
	defer con.Close()
	data, err = redis.Strings(con.Do("zrangebyscore", zsetName, start, end, "withscores"))
	return
}
Zrem
我们对取出的数据操作完成之后，需要将其删除，这里用到zrem命令，zrem命令支持批量删除的。

func RemZset(zsetName string, keys []string) (err error) {
	if len(keys) == 0 {
		return
	}
	con := pool.Get()
	defer con.Close()
	_, err = con.Do("zrem", redis.Args{}.Add(zsetName).AddFlat(keys)...)
	return
}

*/
