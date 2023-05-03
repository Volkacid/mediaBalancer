package service

import (
	"context"
	"github.com/Volkacid/mediaBalancer/internal/storage"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"strings"
	"time"
)

func ServerStatsService(ctx context.Context, servPool *AmazonServices, rdb *storage.RedisStorage) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			iter := rdb.Client.Scan(ctx, 0, "serv:*", 0).Iterator()
			for iter.Next(ctx) {

				res, _ := rdb.Client.ZRangeByScore(ctx, iter.Val(), &redis.ZRangeBy{Min: strconv.FormatInt(time.Now().Unix(), 10), Max: "+inf"}).Result()
				if len(res) == 0 {

					servAddr := iter.Val()
					rdb.Client.Del(ctx, servAddr)

					_, servAddr, _ = strings.Cut(servAddr, "serv:")
					servPool.DestroyServer(servAddr)

					log.Printf("Server %v is destroyed", servAddr)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
