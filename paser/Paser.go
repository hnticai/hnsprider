package paser

import (
	"github.com/loticket/gospider"

	"github.com/zhshch2002/goreq"
)

var GspriderZucaiUrl map[string]string = map[string]string{
	"sfc": "https://webapi.sporttery.cn/gateway/lottery/getFootBallDrawInfoByDrawNumV1.qry?isVerify=1&lotteryGameNum=90&lotteryDrawNum=%s",
	"bqc": "https://webapi.sporttery.cn/gateway/lottery/getFootBallDrawInfoByDrawNumV1.qry?isVerify=1&lotteryGameNum=98&lotteryDrawNum=%s",
	"jqc": "https://webapi.sporttery.cn/gateway/lottery/getFootBallDrawInfoByDrawNumV1.qry?isVerify=1&lotteryGameNum=94&lotteryDrawNum=%s",
}

func PaserLottery(name string, ctx *goreq.Response, hs *gospider.Spider) {

	if name == "zucai" {
		ZcAllPaserLottery(ctx, hs)
	} else if name == "sfc" {
		ZcSfcPaserLottery(ctx, hs)
	} else if name == "bqc" {
		ZcBqcPaserLottery(ctx, hs)
	} else if name == "jqc" {
		ZcJqcPaserLottery(ctx, hs)
	}

}
