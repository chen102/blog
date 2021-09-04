package model

import (
	"blog/tool"
	"strconv"
	"time"
)

type Article struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"omitempty"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"omitempty"`
	Title     string
	UserID    uint   //用户的外键
	UserName  string `gorm:"-" map:"omitempty" `
	Stat      uint
	Tags      string `mapstructure:",omitempty"`
	Content   string
}

//手动处理redis sort的数据
func ArticleList(data []string) ([]Article, error) {
	////Sort返回的结果为[]string，将string转为多个文章模型进行响应
	article := make([]Article, len(data)/6) //5个string为一个article,分别是id,title,time,stat,tags
	id := 0
	for i := 0; i < len(data); i++ {
		if i != 0 && i%6 == 0 {
			id++
		}
		//这样写真的很蠢,但是又没有想到其他的方法，因为返回的是[]string
		switch i % 6 {
		case 0:
			articleid, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}
			article[id].ID = uint(articleid)
		case 1:
			article[id].UserName = data[i]
		case 2:
			article[id].Title = data[i]
		case 3:
			article[id].UpdatedAt = tool.StringToTime(data[i])
		case 4:
			stat, err := strconv.Atoi(data[i])
			if err != nil {
				return nil, err
			}
			article[id].Stat = uint(stat)
		case 5:
			if data[i] != "" {

				article[id].Tags = data[i]
			}
		}

	}
	return article, nil

}
