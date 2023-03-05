# oapi-sdk-gin
an  oapi-sdk-go extension package that integrates the Gin Web framework


# 使用示例

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/larksuite/oapi-sdk-gin"
)


func main() {
	// 创建注册消息处理器
	handler := dispatcher.NewEventDispatcher("v", "e").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	})

	// 创建卡片行为处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))

		// 返回卡片消息
		//return getCard(), nil

		//custom resp
		//return getCustomResp(),nil

		// 无返回值
		return nil, nil
	})

	// 注册处理器
	g := gin.Default()
	g.POST("/webhook/event", sdkginext.NewEventHandlerFunc(handler))
	g.POST("/webhook/card", sdkginext.NewCardActionHandlerFunc(cardHandler))

	// 启动服务
	err := g.Run(":9999")
	if err != nil {
		panic(err)
	}
}


```
