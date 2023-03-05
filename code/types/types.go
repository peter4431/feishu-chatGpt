package types

import (
	"fmt"
	"strings"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type ChatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatMsgPair 一次问答
type ChatMsgPair struct {
	Q string `json:"q"`
	A string `json:"a"`
}

func (c *ChatMsgPair) GetReqStr() string {
	return fmt.Sprintf("Q:%s\nA:%s\n\n", c.Q, c.A)
}

func (c *ChatMsgPair) GetChatMsg() []*ChatMsg {
	return []*ChatMsg{
		&ChatMsg{
			Role:    RoleUser,
			Content: c.Q,
		},
		&ChatMsg{
			Role:    RoleAssistant,
			Content: c.A,
		},
	}
}

// ChatMsgParis 缓存数据
type ChatMsgParis struct {
	pairs []*ChatMsgPair
}

func NewChatMsgPairs(pairs []*ChatMsgPair) *ChatMsgParis {
	return &ChatMsgParis{
		pairs: pairs,
	}
}

func (p *ChatMsgParis) GetList() []*ChatMsg {
	var ret []*ChatMsg
	for _, item := range p.pairs {
		ret = append(ret, item.GetChatMsg()...)
	}
	return ret
}

func (p *ChatMsgParis) GetReqStr() string {
	var strSlice []string
	for _, item := range p.pairs {
		strSlice = append(strSlice, item.GetReqStr())
	}
	return strings.Join(strSlice, "")
}
