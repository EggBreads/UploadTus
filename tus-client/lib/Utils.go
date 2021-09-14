package lib

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

func ConvertToMarshalJson(i interface{}) []byte {
	jsonBytes, err := json.Marshal(i)

	if err != nil {
		log.Print("[ERORR] Fail to converted json marshal")
		return []byte("")
	}

	return jsonBytes
}

func ConvertToUnMarshalJson(b []byte, i interface{}) error {
	log.Println("[DEBUG] UnMarshal String" + string(b))
	return json.Unmarshal(b, i)
}

func SetHeaders(r *http.Request, h map[string]string) {
	for k, v := range h {
		r.Header.Set(k, v)
	}

}

func SetRedisConn() *redis.Client {
	options := &redis.FailoverOptions{
		MasterName:    REDISSENTINELMASTER,
		SentinelAddrs: []string{REDISSENTINELHOST + ":" + REDISSENTINELPORT},
	}
	return redis.NewFailoverClient(options)
}
