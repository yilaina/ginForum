package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go_code/ginStudy/gin_b2/bluebell/models"
	"strconv"
	"time"
)

func getIDsFromKey(key string, page, size int64, c *gin.Context) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZRevRange 查询
	return rdb.ZRevRange(c, key, start, end).Result()

}
func GetPostIDsInorder(p *models.ParamPostList, c *gin.Context) ([]string, error) {
	// 从redis获取id
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2.确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size, c)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string, c *gin.Context) (data []int64, err error) {
	//data = make([]int64,0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := rdb.ZCount(c, key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用pipeline一次发送多条命令，减少rtt
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(c, key, "1", "1")
	}
	cmders, err := pipeline.Exec(c)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

// GetCommunityPostIDsInorder 按社区查询ids
func GetCommunityPostIDsInorder(p *models.ParamPostList, c *gin.Context) ([]string, error) {

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 缓存的key
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(c, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(
			c,
			cKey, // 目标键（缓存结果）
			&redis.ZStore{
				Keys:      []string{key, orderKey}, // 源键列表
				Aggregate: "MAX",                   // 聚合方式：取最大值
			},
		) // zinterstore 计算
		pipeline.Expire(c, orderKey, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(c)
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询id
	return getIDsFromKey(key, p.Page, p.Size, c)
}
