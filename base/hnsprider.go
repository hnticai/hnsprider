package base

import (
	"hnsprider/paser"

	"github.com/loticket/gospider"
	"github.com/zhshch2002/goreq"
)

type Hnsprider struct {
	sprider *gospider.Spider
	stop    chan bool
}

func (hs *Hnsprider) Initialization() {
	hs.sprider = gospider.NewSpider(func(s *gospider.Spider) {
		s.Logging = true
	})

	hs.sprider.OnItem(func(ctx *gospider.Context, i interface{}) interface{} { // 收集并存储结果
		var lotname string = ctx.Meta["name"].(string)
		paser.PaserLottery(lotname, ctx.Resp, hs.sprider)
		return i
	})
	hs.priderInitUrl()
	for {
		select {
		case <-hs.stop:
			goto END
		}
	}
END:
	/*var timers *time.Timer = time.NewTimer(6 * time.Second)

		for {
			select {
			case <-timers.C:

				timers.Reset(6 * time.Second)
			case <-hs.stop:
				goto END
			}
		}
	END:*/
}

func (hs *Hnsprider) priderInitUrl() {
	for name, url := range GspriderInitUrl {
		hs.AddSprider(url, name)
	}
}

func (hs *Hnsprider) AddSprider(url string, names string) {
	hs.sprider.SeedTask(goreq.Get(url), map[string]interface{}{"name": names}, func(ctx *gospider.Context) {
		ctx.AddItem(ctx.Resp.Text)
	})
}

func NewHnsprider() *Hnsprider {
	return &Hnsprider{
		stop: make(chan bool),
	}
}
