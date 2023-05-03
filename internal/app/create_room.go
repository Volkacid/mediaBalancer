package app

import (
	"context"
	"errors"
	"github.com/Volkacid/mediaBalancer/internal/service"
	"github.com/Volkacid/mediaBalancer/internal/storage"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"strings"
	"time"
)

func NewRoom(ctx context.Context, duration int, servPool *service.AmazonServices, rdb *storage.RedisStorage) (string, error) {
	if duration <= 0 {
		return "", errors.New("duration must be positive")
	}

	curTime := time.Now().Unix()
	newRoomID := uuid.New().String()
	newRoomScore := float64(curTime) + float64(duration) // Далее использование ZRANGEBYSCORE(current_time) даст нам список активных комнат
	bestScore := newRoomScore
	var serverAddr string

	iter := rdb.Client.Scan(ctx, 0, "serv:*", 0).Iterator()
	for iter.Next(ctx) {
		res, _ := rdb.Client.ZRangeByScoreWithScores(ctx, iter.Val(), &redis.ZRangeBy{Min: strconv.FormatInt(curTime, 10), Max: "+inf"}).Result()
		if len(res) >= 8 { //Исходим из того, что на сервере может поместиться 8 комнат. На реальном сервере берём метрики и решаем, поместится ли новая комната
			continue
		}

		var curScoreSum float64
		for _, s := range res {

			if s.Score > newRoomScore {
				curScoreSum += s.Score - newRoomScore
			} else {
				curScoreSum += newRoomScore - s.Score
			}
		}

		curScoreAvg := curScoreSum / float64(len(res))
		if curScoreAvg < bestScore {
			bestScore = curScoreAvg
			serverAddr = iter.Val()
		}
	}

	if err := iter.Err(); err != nil {
		return "", err
	}

	_, serverAddr, _ = strings.Cut(serverAddr, "serv:")

	if bestScore == newRoomScore { //Если ни один из серверов не подошёл
		serverAddr = servPool.CreateNewServer()
		log.Printf("New server was created(%s)", serverAddr)
	}

	rdb.Client.ZAdd(ctx, "serv:"+serverAddr, redis.Z{Score: newRoomScore, Member: newRoomID})
	log.Printf("Room %s assigned to %s for %d seconds", newRoomID, serverAddr, duration)

	return serverAddr, nil
}
