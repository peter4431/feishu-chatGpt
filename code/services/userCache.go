package services

import (
	"github.com/patrickmn/go-cache"
	"start-feishubot/types"
	"time"
)

type UserService struct {
	cache *cache.Cache
}

var userServices *UserService

// GetList 获得用于请求 chat/completion 的数据
func (u UserService) GetList(userId string) []*types.ChatMsg {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(userId)
	if !ok {
		return nil
	}
	//list to string
	list := sessionContext.([]*types.ChatMsgPair)
	var pairs = types.NewChatMsgPairs(list)
	return pairs.GetList()
}

func (u UserService) Get(userId string) string {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(userId)
	if !ok {
		return ""
	}
	//list to string
	list := sessionContext.([]*types.ChatMsgPair)
	var pairs = types.NewChatMsgPairs(list)
	return pairs.GetReqStr()
}

func (u UserService) Set(userId string, question, reply string) {
	// 列表，最多保存8个
	//如果满了，删除最早的一个
	//如果没有满，直接添加
	maxCache := 8
	maxLength := 2048
	maxCacheTime := time.Minute * 30
	listOut := make([]*types.ChatMsgPair, 0)

	value := &types.ChatMsgPair{
		Q: question,
		A: reply,
	}

	raw, ok := u.cache.Get(userId)
	if ok {
		listOut = raw.([]*types.ChatMsgPair)
		if len(listOut) == maxCache {
			listOut = listOut[1:]
		}
		listOut = append(listOut, value)
	} else {
		listOut = append(listOut, value)
	}

	//限制对话上下文长度
	if getStrPoolTotalLength(listOut) > maxLength {
		listOut = listOut[1:]
	}
	u.cache.Set(userId, listOut, maxCacheTime)
}

func (u UserService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

type UserCacheInterface interface {
	GetList(userId string) []*types.ChatMsg
	Get(userId string) string
	Set(userId string, question, reply string)
	Clear(userId string) bool
}

func GetUserCache() UserCacheInterface {
	if userServices == nil {
		userServices = &UserService{cache: cache.New(30*time.Minute, 30*time.Minute)}
	}
	return userServices
}

func getStrPoolTotalLength(strPool []*types.ChatMsgPair) int {
	var total int
	for _, v := range strPool {
		total += len(v.GetReqStr())
	}
	return total
}
