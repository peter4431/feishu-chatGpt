package handlers

import (
	"context"
	"fmt"
	"start-feishubot/services"
	"start-feishubot/types"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/viper"
)

type GroupMessageHandler struct {
	userCache services.UserCacheInterface
	msgCache  services.MsgCacheInterface
}

func (p GroupMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	ifMention := judgeIfMentionMe(event)
	if !ifMention {
		return nil
	}
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	sender := event.Event.Sender
	openId := sender.SenderId.OpenId
	chatId := event.Event.Message.ChatId

	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	qParsed := parseContent(*content)
	if len(qParsed) == 0 {
		sendMsg(ctx, "🤖️：你想知道什么呢~", chatId)
		fmt.Println("msgId", *msgId, "message.text is empty")
		return nil
	}

	if qParsed == "/clear" || qParsed == "清除" {
		p.userCache.Clear(*openId)
		sendMsg(ctx, "🤖️：AI机器人已清除记忆", chatId)
		return nil
	}

	var (
		completions string
		err         error
	)
	completions, err = p.sendPrompt(ctx, *openId, qParsed)

	ok := true
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
	}
	if len(completions) == 0 {
		ok = false
	}
	if ok {
		p.userCache.Set(*openId, qParsed, completions)
		err := replyMsg(ctx, completions, msgId)
		if err != nil {
			replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
			return nil
		}
	}
	return nil

}

func (p GroupMessageHandler) sendPrompt(ctx context.Context, userId string, qStr string) (res string, err error) {
	var useChatCompletion = true // 使用新接口
	if useChatCompletion {
		msgList := p.userCache.GetList(userId)
		msgList = append(msgList, &types.ChatMsg{
			Role:    types.RoleUser,
			Content: qStr,
		})
		res, err = services.ChatCompletion(msgList)
	} else {
		prompt := p.userCache.Get(userId)
		prompt = fmt.Sprintf("%s\nQ:%s\nA:", prompt, qStr)
		res, err = services.Completions(prompt)
	}
	return
}

var _ MessageHandlerInterface = (*PersonalMessageHandler)(nil)

func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{
		userCache: services.GetUserCache(),
		msgCache:  services.GetMsgCache(),
	}
}

func judgeIfMentionMe(event *larkim.P2MessageReceiveV1) bool {
	mention := event.Event.Message.Mentions
	if len(mention) != 1 {
		return false
	}
	return *mention[0].Name == viper.GetString("BOT_NAME")
}
