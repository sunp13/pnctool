package pnctool

import (
	"fmt"
	"time"

	"github.com/raff/godet"
)

var (
	remoteAddr string
)

// Turl struct
type Turl struct {
	Type string
	URL  string
}

// Init func
func Init(addr string) error {
	remote, err := godet.Connect(addr, false)
	if err != nil {
		return err
	}
	defer remote.Close()

	tabs, _ := remote.TabList("")
	fmt.Println("tabs length", len(tabs))

	remoteAddr = addr
	return nil
}

// CollectURL 获取网页内容
func CollectURL(url string, timeout ...time.Duration) ([]Turl, error) {
	if "" == remoteAddr {
		return nil, fmt.Errorf("Remote addr invalid")
	}
	t := 3 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	remote, err := godet.Connect(remoteAddr, false)
	if err != nil {
		return nil, err
	}
	defer remote.Close()

	// 创建tab
	tab, err := remote.NewTab("about:blank")
	if err != nil {
		return nil, err
	}
	defer remote.CloseTab(tab)

	urls := make([]Turl, 0)
	// 开启network event
	remote.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
		// fmt.Println(params["type"], params["request"].(map[string]interface{})["url"].(string))
		urls = append(urls, Turl{
			Type: params["type"].(string),
			URL:  params["request"].(map[string]interface{})["url"].(string),
		})
	})
	remote.NetworkEvents(true)
	remote.Navigate(url)
	remote.EvaluateWrap(`
		var i = 0;
		setInterval(function(){
			i+=500;
			window.scrollTo(0,i)
		},500);
	`)
	time.Sleep(t)
	return urls, nil
}
