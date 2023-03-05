package handlers

import (
	"context"
	"start-feishubot/initialization"
	"testing"
)

func TestSendPrompt(t *testing.T) {
	t.Run("测试发送数据", func(t *testing.T) {
		initialization.LoadConfig()
		h := NewPersonalMessageHandler()
		ctx := context.Background()
		res, err := h.sendPrompt(ctx, "test", "你好")
		if err != nil {
			t.Fail()
		}
		t.Log(res)
	})
}
