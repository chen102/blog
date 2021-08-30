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
//如何防止重复点赞？
//幂等性原本是数学上的概念，即使公式：f(x)=f(f(x)) 能够成立的数学性质。用在编程领域，则意为对同一个系统，使用同样的条件，一次请求和重复的多次请求对系统资源的影响是一致的。
//1.去重表：在表上建立唯一索引，保证某一类数据一旦执行完毕，后续同样的请求再也无法成功写入,比如这里可以将文章ID与用户id绑定唯一索引，这样重复点赞的数据就无法写入
//2.token
//是为每一次操作生成一个唯一性的凭证，也就是token。一个token在操作的每一个阶段只有一次执行权，一旦执行成功则保存执行结果。对重复的请求，返回同一个结果

func StatCronService() {
	//计算列表长度，按长度取相应的数据
	lenght, err := model.RedisSysDB.LLen(redis.UserStatQueueKey()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("点赞定时任务：长度为:", lenght, "正在写入mysql")
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
		stat.UserID = uint(userid)
		stat.ArticleID = uint(artid)
		if len(arr) == 3 {
			if err := model.DB.Model(&stat).Update("stat", 1); err != nil {
				panic(err)
			}
			continue
		}
		if err := model.DB.Create(&stat).Error; err != nil {
			panic(err)
		}

	}
}
