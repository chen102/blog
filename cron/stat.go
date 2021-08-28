package cron

import (
	"blog/model"
	"blog/redis"
	//"github.com/go-redis/redis"
	//"fmt"
	"log"
	"strconv"
	"strings"
)

//点赞时间服务
//固定时间将redis的点赞情况写入mysql
const StatCronServiceSpec = "*/30 * * * * ?" //每10秒一次
//这个无论有没有点赞都会执行redis查操作(轮询)
func StatCronService() {
	//计算列表长度，按长度取相应的数据
	lenght, err := model.RedisSysDB.LLen(redis.UserStatQueueKey()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("点赞：长度为:", lenght, "正在写入mysql")
	for ; lenght > 0; lenght-- {
		userstats, err := model.RedisSysDB.RPop(redis.UserStatQueueKey()).Result()
		if err != nil {
			panic(err)
		}
		var stat model.Stat
		log.Println(userstats)
		arr := strings.Split(userstats, ":")
		userid, err := strconv.Atoi(arr[0])
		if err != nil {
			panic(err)
		}
		artid, err := strconv.Atoi(arr[1])
		if err != nil {
			panic(err)
		}
		if len(arr) == 3 {
			stat.Stat = true
		}
		stat.UserID = uint(userid)
		stat.ArticleID = uint(artid)
		if err := model.DB.Create(&stat).Error; err != nil {
			panic(err)
		}

	}
}
