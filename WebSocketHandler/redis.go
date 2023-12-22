package WebSocketHandler

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var Rdb redis.Conn

func RedisSet(key string, value string, time int) (bool, error) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
		//Rdb.Do("SELECT", "index")
	} else {
		fmt.Println("connect redis error :", err)
		return false, err
	}
	defer Rdb.Close()
	_, errs := Rdb.Do("SET", key, value)
	if errs == nil {
		if time > 0 {
			Rdb.Do("expire", key, time)
		}
		return true, errs
	} else {
		return false, errs
	}
}

func RedisGet(key string) (string, bool) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
		//Rdb.Do("SELECT", "index")
	} else {
		fmt.Println("connect redis error :", err)
		return "", false
	}
	defer Rdb.Close()
	v, err := redis.String(Rdb.Do("GET", key))
	if err != nil {
		return "", false
	} else {
		return v, true
	}
}

func RedisListLen(key string) (int, bool) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return 0, false
	}
	defer Rdb.Close()
	len, errs := redis.Int(Rdb.Do("LLEN", key))
	if errs == nil {
		return len, true
	} else {
		return 0, false
	}
}

func RedisListLpush(key string, value interface{}) (bool, error) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return false, err
	}
	defer Rdb.Close()

	_, errs := Rdb.Do("LPUSH", key, value)
	if errs == nil {
		return true, errs
	} else {
		return false, errs
	}
}

func RedisListRpop(key string) ([]byte, bool) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		r := []byte("fail8")
		return r, false
	}
	defer Rdb.Close()
	data, errs := redis.Bytes(Rdb.Do("RPOP", key))
	if errs == nil {
		return data, true
	} else {
		r := []byte("fail9")
		return r, false
	}
}

func RedisHset(key, id string, v int) (bool, error) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return false, err
	}
	defer Rdb.Close()
	_, errs := Rdb.Do("hSET", key, id, v)
	if err != nil {
		return false, errs
	} else {
		return true, errs
	}
}

func RedisHget(key, id string) (string, bool) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return "", false
	}
	defer Rdb.Close()
	v, err := redis.String(Rdb.Do("hGET", key, id))
	if err != nil {
		return "", false
	} else {
		return v, true
	}
}

/**
 *redis SADD 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
 */
func RdbSAdd(key, v string) (bool, error) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return false, err
	}
	defer Rdb.Close()
	_, errs := Rdb.Do("SADD", key, v)
	if errs == nil {
		return true, errs
	} else {
		fmt.Println("RdbSAdd", err)
		return false, errs
	}
}

// 获取集合中元素的数量
func RdbGetSetScard(key string) (int, bool) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return 0, false
	}
	defer Rdb.Close()
	len, errs := redis.Int(Rdb.Do("SCARD", key))
	if errs == nil {
		return len, true
	} else {
		//fmt.Println("RdbGetSetScard", errs)
		return 0, false
	}
}

/**
 *redis SMEMBERS 返回集合 key 中的所有成员。return map
 */
func RdbSMembers(key string) (interface{}, error) {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return nil, err
	}
	defer Rdb.Close()
	data, err := redis.Strings(Rdb.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

/**
 *redis SISMEMBER 判断 member 元素是否集合 key 的成员。return bool
 */
func RdbSISMembers(key, v string) bool {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return false
	}
	defer Rdb.Close()
	b, err := redis.Bool(Rdb.Do("SISMEMBER", key, v))
	if err != nil {
		fmt.Println("RdbSISMembers false", err)
		return false
	} else {
		fmt.Println("RdbSISMembers true", b)
		return b
	}
}

/**
 *设置可以过期，redis EXPIRE
 */
func RdbSetKeyExp(key string, ex int) bool {
	Rdb, err := redis.Dial("tcp", REDIS_CONN)
	if err == nil {
		Rdb.Do("AUTH", REDIS_PASS)
		Rdb.Do("SELECT", REDIS_DB)
	} else {
		fmt.Println("connect redis error :", err)
		return false
	}
	defer Rdb.Close()
	_, errs := Rdb.Do("EXPIRE", key, ex)
	if errs != nil {
		fmt.Println("connect redis error :", errs)
		return false
	}
	return true
}
