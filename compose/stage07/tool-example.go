package stage07

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func ToolExample() {
	ctx := context.Background()
	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}

	url := "https://www.bing.com"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	time.Sleep(5 * time.Second)
	but.Cleanup()
}
