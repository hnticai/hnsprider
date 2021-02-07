package paser

import (
	"errors"
	"hnsprider/db"

	"github.com/loticket/gospider"

	"github.com/zhshch2002/goreq"
)

func ZcJqcPaserLottery(ctx *goreq.Response, hs *gospider.Spider) error {
	ts, err := ctx.JSON()
	if err != nil {
		return nil
	}

	if !ts.Exists() {
		return nil
	}

	result := ts.Map()
	if _, ok := result["value"]; !ok {
		return errors.New("抓取信息不包含信息解析信息")
	}

	var jqc LotteryTraditionalBall = ZcPaserLotteryCommon(result["value"].Map(), "进球彩")
	var has bool
	if !has && err == nil && jqc.Number != "" && jqc.TotalSales != "" {
		db.Engine.Insert(&jqc)
	}

	return nil
}
