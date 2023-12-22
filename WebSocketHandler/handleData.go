package WebSocketHandler

import (
	"encoding/json"
	"fmt"
)

type TheHandleData struct {
	Source string     `form:"source" json:"source" uri:"source" xml:"source" binding:"required"`
	Cid    string     `form:"cid" json:"cid" uri:"cid" xml:"cid" binding:"required"`
	Msg    TheMsgData `form:"msg" json:"msg" uri:"msg" xml:"msg" binding:"required"`
}

type TheMsgData struct {
	Content  string       `form:"content" json:"content" uri:"content" xml:"content" binding:"required"`
	Msg_type int          `form:"msg_type" json:"msg_type" uri:"msg_type" xml:"msg_type" binding:"required"`
	Param    TheParamData `form:"param" json:"param" uri:"param" xml:"param" binding:"required"`
}

type TheParamData struct {
	Id       string                 `form:"id" json:"id" uri:"id" xml:"id" binding:"required"`
	Msg_id   string                 `form:"msg_id" json:"msg_id" uri:"msg_id" xml:"msg_id" binding:"required"`
	Msg_info map[string]interface{} `form:"msg_info" json:"msg_info" uri:"msg_info" xml:"msg_info" binding:"required"`
}

func HandleData(msgType int, data []byte) (result []byte, err bool) {
	// 判断是否为文本数据，如果二进制数据则直接返回数据格式错误
	r := []byte("no data")
	if msgType != 1 {
		return r, false
	}

	var the_handle_data TheHandleData
	json.Unmarshal(data, &the_handle_data)

	source := the_handle_data.Source
	cid := the_handle_data.Cid
	//fmt.Println(the_handle_data, source, cid, msgType, data)

	// 判断来源是否web
	if source == "web" {
		// 获取系统消息数据
		result, err := getXtData(cid)
		if err == true {
			return result, err
		} else {
			// 获取用户专有消息数据
			result, err := getData(cid)
			return result, err
		}
	} else if source == "server" {
		// 设置数据
		msg := the_handle_data.Msg
		msg_type := msg.Msg_type
		//系统消息
		if msg_type == 2 {
			fmt.Println(msg_type, cid, "xt server->msg_type")
			//设置系统消息
			result, err := SetXtData(cid, data, msg_type, msg.Param.Msg_id)
			return result, err
		} else {
			fmt.Println(msg_type, cid, "not-xt server->msg_type")
			//fmt.Println(msg, the_handle_data, "php")
			result, err := setData(cid, data)
			return result, err
		}
	}
	re := []byte("null data")
	return re, false
}

func getData(cid string) (result []byte, err bool) {
	redis_list_key := REDIS_LIST_KEY + cid
	len, _ := RedisListLen(redis_list_key)
	if len == 0 {
		r := []byte("no data")
		return r, false
	}
	re, e := RedisListRpop(redis_list_key)
	if e == false {
		r := []byte("no data")
		return r, false
	}
	fmt.Println(redis_list_key, len, re, e, "getData")

	return re, true
}

func setData(cid string, msg interface{}) (result []byte, err bool) {
	redis_list_key := REDIS_LIST_KEY + cid
	r, _ := RedisListLpush(redis_list_key, msg)
	if r == true {
		re := []byte("success")
		return re, true
	} else {
		re := []byte("fail1")
		return re, false
	}
}

func getXtData(cid string) (result []byte, err bool) {
	// 生产环境测试，限定cid，测试用，未来需注释掉
	//if cid != "110_128" {
	//	re := []byte("fail8")
	//	return re, false
	//}

	res, err := getData("_")
	if err == false {
		re := []byte("fail2")
		return re, false
	}

	//将系统消息放回系统消息list
	setData("_", res)

	//判断集合数据是否失效
	var xt_handle_data TheHandleData
	json.Unmarshal(res, &xt_handle_data)
	msg := xt_handle_data.Msg
	msg_param := msg.Param
	msg_id := msg_param.Msg_id
	redis_set_key := REDIS_SET_KEY + msg_id
	fmt.Println("getXtData", redis_set_key, xt_handle_data, cid)

	_, er := RdbGetSetScard(redis_set_key)
	if er == false {
		fmt.Println("RdbGetSetScard is err", redis_set_key, cid)
		re := []byte("fail3")
		return re, false
	}
	//判断用户是否已接收过数据
	eb := RdbSISMembers(redis_set_key, cid)
	if eb == true {
		fmt.Println("RdbSISMembers is err", redis_set_key, cid)
		re := []byte("fail4")
		return re, false
	}

	//添加用户接收记录
	RdbSAdd(redis_set_key, cid)

	return res, true
}

func SetXtData(cid string, msg interface{}, msg_type int, msg_id string) (result []byte, err bool) {
	if msg_type != 2 {
		re := []byte("fail7")
		return re, false
	}
	//添加过期时间为7天的集合，"_"用于占位
	redis_set_key := REDIS_SET_KEY + msg_id
	rs, _ := RdbSAdd(redis_set_key, "_")
	if rs == false {
		re := []byte("fail5")
		return re, false
	}
	//设置过期时间为7天
	RdbSetKeyExp(redis_set_key, 7*24*3600)

	redis_list_key := REDIS_LIST_KEY + cid
	//清空之前数据
	RedisListRpop(redis_list_key)
	//插入新数据
	r, _ := RedisListLpush(redis_list_key, msg)
	if r == true {
		re := []byte("success")
		return re, true
	} else {
		re := []byte("fail6")
		return re, false
	}
}
