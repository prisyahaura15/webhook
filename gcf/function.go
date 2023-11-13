package gcf

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/whatsauth/webhook"
)

//biar si livelocationnya jalan juga

func init() {
	functions.HTTP("WebHook", webhook.Post)
}
