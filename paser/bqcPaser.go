package paser

import (
	"errors"
	"hnsprider/db"

	"github.com/loticket/gospider"

	"github.com/zhshch2002/goreq"
)

func ZcBqcPaserLottery(ctx *goreq.Response, hs *gospider.Spider) error {
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

	var bqc LotteryTraditionalBall = ZcPaserLotteryCommon(result["value"].Map(), "半全彩")
	var has bool
	has, err = db.Engine.Where("term = ? and type = ? ", bqc.Term, bqc.Types).Cols("id").Exist(&LotteryTraditionalBall{})
	if !has && err == nil && bqc.Number != "" && bqc.TotalSales != "" {
		db.Engine.Insert(&bqc)
	}

	return nil
}
