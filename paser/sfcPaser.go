package paser

import (
	"errors"
	"hnsprider/db"

	"github.com/loticket/gospider"

	"github.com/zhshch2002/goreq"
)

func ZcSfcPaserLottery(ctx *goreq.Response, hs *gospider.Spider) error {
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

	var sfc LotteryTraditionalBall = ZcPaserLotteryCommon(result["value"].Map(), "胜负彩")

	var has bool
	has, err = db.Engine.Where("term = ? and type = ? ", sfc.Term, sfc.Types).Cols("id").Exist(&LotteryTraditionalBall{})
	if !has && err == nil && sfc.Number != "" && sfc.TotalSales != "" {
		db.Engine.Insert(&sfc)
	}
	return nil
}
