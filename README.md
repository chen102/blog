# blog cudr
博客发布平台
目前采用gin+gorm+go-redis
路由：
    使用gin作为路由处理与路由相应，后续可能会添加gRPC支持
业务：
    使用redis实现blog业务逻辑包括blog的cudr、文章的点击量、排行榜、发布订阅、评论、私信等等,使用mysql实现数据的持久化
项目进度：
    用户管理 ×
    私信管理 ×
    blog管理 ×
    评论管理 ×
