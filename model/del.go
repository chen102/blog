package model

//目前是直接读环境参数，后续会添加配置文件
import "log"
import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var RedisReadDB *redis.Client
var RedisWriteDB *redis.Client
var RedisSysDB *redis.Client
var DB *gorm.DB

const RedisNil = redis.Nil

const RedisKeyNull = redis.Nil //结果为空
//检查redis主从是否开启
func check() bool {
	RedisWriteDB.Set("MasterSlave", 1, 3*time.Second)
	ok, err := RedisReadDB.Get("MasterSlave").Int()
	if err != nil && err != RedisNil {
		panic(err)
	}
	if ok == 1 {
		return true
	}
	return false
}

//读写分离
func DelRedis() {
	RedisReadDB = redis.NewClient(&redis.Options{
		Addr:     "192.168.122.1:6380", //这里不能两个都用回环，走的是网络层，需要模拟两个网段
		Password: "000000",
		DB:       0,
	})
	RedisWriteDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", //
		Password: "000000",         //若主不设置密码，redis主从起不来 具体看redis配置
		DB:       0,
	})
	//系统redis客户端
	RedisSysDB = redis.NewClient(&redis.Options{
		Addr:     "192.168.122.1:6379",
		Password: "000000",
		DB:       0,
	})
	pong, err := RedisWriteDB.Ping().Result()
	log.Println(pong, err)
	pong, err = RedisReadDB.Ping().Result()
	log.Println(pong, err)
	if err := RedisReadDB.SlaveOf("172.17.0.2", "6379").Err(); err != nil { //redis主从
		//这里的ip是主的IP，因为这里是本地，所以直接用的容器的IP (需固定容器IP)
		panic(err)
	}
	if !check() {
		panic("主从开启失败")
	}
	log.Println("redis主从启动成功")
}
func DelMysql() error {
	db, err := gorm.Open("mysql", "root:chenxi1234@tcp(10.177.3.141:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return err
	}
	db.LogMode(true) //开启gorm日志
	//gorm连接池
	//查看连接池是否有空闲连接，若有直接从连接池取出一个连接，若没有，判断连接数是否大于最大连接数，若小于创建新的连接
	db.DB().SetMaxIdleConns(20) //最大可空闲连接数
	//理论上maxIdleConns连接的上限越高，也即允许在连接池中的空闲连接最大值越大，可以有效减少连接创建和销毁的次数，提高程序的性能。但是连接对象也是占用内存资源的，而且如果空闲连接越多，存在于连接池内的时间可能越长。连接在经过一段时间后有可能会变得不可用，而这时连接还在连接池内没有回收的话，后续被征用的时候就会出问题。一般建议maxIdleConns的值为MaxOpenConns的1/2，仅供参考。
	db.DB().SetMaxOpenConns(100) //最大连接数
	//默认情况下，连接池的最大数量是没有限制的。一般来说，连接数越多，访问数据库的性能越高。但是系统资源不是无限的，数据库的并发能力也不是无限的。因此为了减少系统和数据库崩溃的风险，可以给并发连接数设置一个上限，这个数值一般不超过进程的最大文件句柄打开数，不超过数据库服务自身支持的并发连接数，比如1000。
	db.DB().SetConnMaxLifetime(time.Second * 30) //一个连接使用的最大时长
	//设置一个连接被使用的最长时间，即过了一段时间后会被强制回收，理论上这可以有效减少不可用连接出现的概率。当数据库方面也设置了连接的超时时间时，这个值应当不超过数据库的超时参数值。

	db.Set("gorm:table_options", "charset=utf8mb4") //tips:mysql容器中的默认编码是临时的,容器重启了就没了
	db.AutoMigrate(&User{}, &Article{}, &Stat{})
	DB = db
	log.Println("mysql启动成功")
	return nil
}
