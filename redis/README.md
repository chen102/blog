//做缓存的前提是数据不保持实时一致，如果需要实时一致的数据，那就失去了为MySQL做缓存的意义，还不如直接在redis里实现业务。
//MySQL和Redis的双写一致性
//策略
//1.更redis，更db
//若更db失败，造成缓存脏数据
//2.更db,更redis
//线程a更新了db，但由于某种原因，更redis延迟了，这时线程b更新了db，更redis，然后线程a才更新redis即线程b对缓存的更新丢失了
//3.删redis，更db
//线程a在写操作，删除redis，然后准备更新数据库，这时线程b进行读操作，发现缓存未命中，读db，更新了redis，线程a完成了数据库的更新，即数据库与缓存又不一致
//4.更db，删redis
//线程a读取数据，准备将数据写入缓存时，线程b更新了数据库，然后执行了删除缓存的操作，线程a将原来的旧值写入缓存不过这种情况出现的概率较低，因为写操作的时间一般要比读操作时间长

//缓存的三大问题
//1.缓存穿透
//业务系统要查询的数据根本就存在！当业务系统发起查询时，按照上述流程，首先会前往缓存中查询，由于缓存中不存在，然后再前往数据库中查询。由于该数据压根就不存在，因此数据库也返回空。这就是缓存穿透
//缓存穿透的解决方案
//缓存空数据
//之所以发生缓存穿透，是因为缓存中没有存储这些空数据的key，导致这些请求全都打到数据库上。那么，我们可以稍微修改一下业务系统的代码，将数据库查询结果为空的key也存储在缓存中。当后续又出现该key的查询请求时，缓存直接返回null，而无需查询数据库。
//缓存空对象会有两个问题：
//第一，空值做了缓存，意味着缓存层中存了更多的键，需要更多的内存空间 ( 如果是攻击，问题更严重 )，比较有效的方法是针对这类数据设置一个较短的过期时间，让其自动剔除。
//第二，缓存层和存储层的数据会有一段时间窗口的不一致，可能会对业务有一定影响。例如过期时间设置为 5 分钟，如果此时存储层添加了这个数据，那此段时间就会出现缓存层和存储层数据的不一致，此时可以利用消息系统或者其他方式清除掉缓存层中的空对象。
// BloomFilter布隆过滤器
//它需要在缓存之前再加一道屏障，里面存储目前数据库中存在的所有key
//当业务系统有查询请求的时候，首先去BloomFilter中查询该key是否存在。若不存在，则说明数据库中也不存在该数据，因此缓存都不要查了，直接返回null。若存在，则继续执行后续的流程，先前往缓存中查询，缓存中没有的话再前往数据库中的查询。
//这种方法适用于数据命中不高，数据相对固定实时性低（通常是数据集较大）的应用场景，代码维护较为复杂，但是缓存空间占用少。
//布隆过滤器
//布隆过滤器（Bloom Filter） 是由 Howard Bloom在1970年提出的二进制向量数据结构，它具有很好的空间和时间效率，被用来检测一个元素是不是集合中的一个成员，即判定 “可能已存在和绝对不存在” 两种情况。如果检测结果为是，该元素不一定在集合中；但如果检测结果为否，该元素一定不在集合中,因此Bloom filter具有100%的召回率。
//布隆过滤器的核心是一个超大的位数组和几个哈希函数。假设位数组的长度为m,哈希函数的个数为k。下图表示有三个hash函数，比如一个集合中有x，y，z三个元素，分别用三个hash函数映射到二进制序列的某些位上，假设我们判断w是否在集合中，同样用三个hash函数来映射，结果发现取得的结果不全为1，则表示w不在集合里面。
//可以使用redis的bitmap+原生语言实现一个简单的布隆过滤器
//2.缓存雪崩
//缓存其实扮演了一个保护数据库的角色。它帮数据库抵挡大量的查询请求，从而避免脆弱的数据库受到伤害。如果缓存因某种原因发生了宕机，那么原本被缓存抵挡的海量查询请求就会像疯狗一样涌向数据库。此时数据库如果抵挡不了这巨大的压力，它就会崩溃。这就是缓存雪崩。
//和飞机都有多个引擎一样，如果缓存层设计成高可用的，即使个别节点、个别机器、甚至是机房宕掉，依然可以提供服务，例如前面介绍过的 Redis Sentinel 和 Redis Cluster 都实现了高可用。
//Hystrix是一款开源的“防雪崩工具”，它通过 熔断、降级、限流三个手段来降低雪崩发生后的损失。
//Hystrix就是一个Java类库，它采用命令模式，每一项服务处理请求都有各自的处理器。所有的请求都要经过各自的处理器。处理器会记录当前服务的请求失败率。一旦发现当前服务的请求失败率达到预设的值，Hystrix将会拒绝随后该服务的所有请求，直接返回一个预设的结果。这就是所谓的“熔断”。当经过一段时间后，Hystrix会放行该服务的一部分请求，再次统计它的请求失败率。如果此时请求失败率符合预设值，则完全打开限流开关；如果请求失败率仍然很高，那么继续拒绝该服务的所有请求。这就是所谓的“限流”。而Hystrix向那些被拒绝的请求直接返回一个预设结果，被称为“降级”。

//3.缓存击穿
//我们一般都会给缓存设定一个失效时间，过了失效时间后，该数据库会被缓存直接删除，从而一定程度上保证数据的实时性。但是，对于一些请求量极高的热点数据而言，一旦过了有效时间，此刻将会有大量请求落在数据库上，从而可能会导致数据库崩溃
//如果某一个热点数据失效，那么当再次有该数据的查询请求[req-1]时就会前往数据库查询。但是，从请求发往数据库，到该数据更新到缓存中的这段时间中，由于缓存中仍然没有该数据，因此这段时间内到达的查询请求都会落到数据库上，这将会对数据库造成巨大的压力。此外，当这些请求查询完成后，都会重复更新缓存。
//互斥锁解决
//此方法只允许一个线程重建缓存，其他线程等待重建缓存的线程执行完，重新从缓存获取数据即可
//当第一个数据库查询请求发起后，就将缓存中该数据上锁；此时到达缓存的其他查询请求将无法查询该字段，从而被阻塞等待；当第一个请求完成数据库查询，并将数据更新值缓存后，释放锁；此时其他被阻塞的查询请求将可以直接从缓存中查到该数据。当某一个热点数据失效后，只有第一个数据库查询请求发往数据库，其余所有的查询请求均被阻塞，从而保护了数据库。但是，由于采用了互斥锁，其他请求将会阻塞等待，此时系统的吞吐量将会下降。这需要结合实际的业务考虑是否允许这么做。
//互斥锁可以避免某一个热点数据失效导致数据库崩溃的问题而在实际业务中，往往会存在一批热点数据同时失效的场景。那么，对于这种场景该如何防止数据库过载呢？
//设置不同的失效时间
//当我们向缓存中存储这些数据的时候，可以将他们的缓存失效时间错开。这样能够避免同时失效。如：在一个基础时间上加/减一个随机数，从而将这些缓存的失效时间错开
//永远不过期解决
//“永远不过期”包含两层意思：
//从缓存层面来看，确实没有设置过期时间，所以不会出现热点 key 过期后产生的问题，也就是“物理”不过期。
//从功能层面来看，为每个 value 设置一个逻辑过期时间，当发现超过逻辑过期时间后，会使用单独的线程去构建缓存。
//从实战看，此方法有效杜绝了热点 key 产生的问题，但唯一不足的就是重构缓存期间，会出现数据不一致的情况，这取决于应用方是否容忍这种不一致。