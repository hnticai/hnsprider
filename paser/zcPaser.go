package paser

import (
	"errors"
	"fmt"
	"time"

	"strings"

	"github.com/loticket/gospider"
	"github.com/tidwall/gjson"
	"github.com/zhshch2002/goreq"
)

type LotteryTraditionalBall struct {
	Id           int    `xorm:"not null pk autoincr 'id'" json:"id"`
	Types        string `xorm:"type" json:"type"`
	Term         string `xorm:"term" json:"term"`
	Pool         string `xorm:"pool" json:"pool"`
	TotalSales   string `xorm:"totalSales" json:"totalSales"`
	Number       string `xorm:"number" json:"number"`
	OpenTime     string `xorm:"openTime_fmt" json:"openTime_fmt"`
	Level        string `xorm:"level" json:"level"`
	Allmoney     string `xorm:"allmoney" json:"allmoney"`
	Piece        string `xorm:"piece" json:"piece"`
	HomeTeamView string `xorm:"homeTeamView" json:"homeTeamView"`
	AwayTeamView string `xorm:"awayTeamView" json:"awayTeamView"`
	Money        string `xorm:"moneys" json:"moneys"`
}

//老足彩
func ZcAllPaserLottery(ctx *goreq.Response, hs *gospider.Spider) error {
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

	var values map[string]gjson.Result = result["value"].Map()

	/*var bqclist []gjson.Result = values["bqclist"].Array()
	for i := len(bqclist) - 1; i > 0; i-- {
		var issue string = bqclist[i].String()
		var urls string = fmt.Sprintf(GspriderZucaiUrl["bqc"], issue) //
		hs.SeedTask(goreq.Get(urls), map[string]interface{}{"name": "bqc"}, func(ctx *gospider.Context) {
			ctx.AddItem(ctx.Resp.Text)
		})

		time.Sleep(3 * time.Second)

	}*/

	var jqclist []gjson.Result = values["jqclist"].Array()
	for i := len(jqclist) - 1; i > 0; i-- {
		var issue string = jqclist[i].String()
		var urls string = fmt.Sprintf(GspriderZucaiUrl["jqc"], issue) //
		hs.SeedTask(goreq.Get(urls), map[string]interface{}{"name": "jqc"}, func(ctx *gospider.Context) {
			ctx.AddItem(ctx.Resp.Text)
		})

		time.Sleep(3 * time.Second)
	}

	/*var sfclist []gjson.Result = values["sfclist"].Array()
	for i := len(sfclist) - 1; i > 0; i-- {
		var issue string = sfclist[i].String()
		var urls string = fmt.Sprintf(GspriderZucaiUrl["sfc"], issue) //
		hs.SeedTask(goreq.Get(urls), map[string]interface{}{"name": "sfc"}, func(ctx *gospider.Context) {
			ctx.AddItem(ctx.Resp.Text)
		})

		time.Sleep(3 * time.Second)
	}*/

	return err
}

//解析公共数据
func ZcPaserLotteryCommon(sfcDetail map[string]gjson.Result, lotname string) LotteryTraditionalBall {
	var LotteryZuCai LotteryTraditionalBall = LotteryTraditionalBall{}
	LotteryZuCai.Types = lotname
	LotteryZuCai.Term = sfcDetail["lotteryDrawNum"].String()
	LotteryZuCai.Pool = strings.Replace(sfcDetail["poolBalanceAfterdraw"].String(), ",", "", -1)
	LotteryZuCai.TotalSales = strings.Replace(sfcDetail["drawFlowFund"].String(), ",", "", -1)
	LotteryZuCai.Number = sfcDetail["lotteryDrawResult"].String()
	LotteryZuCai.OpenTime = sfcDetail["estimateDrawTime"].String() + " 00:00:00"
	var sfHome []string = make([]string, 0)
	var sfGuest []string = make([]string, 0)
	var level []string = make([]string, 0)
	var allmoney []string = make([]string, 0)
	var money []string = make([]string, 0)
	var piece []string = make([]string, 0)
	var sfcDetails []gjson.Result = sfcDetail["prizeLevelList"].Array() //赛程信息
	var sfcMatchResults []gjson.Result = sfcDetail["matchList"].Array() //赛程信息
	for _, matchresult := range sfcMatchResults {
		detailMapRes := matchresult.Map()
		sfHome = append(sfHome, strings.Replace(detailMapRes["masterTeamName"].String(), " ", "", -1))
		if lotname == "进球彩" {
			sfHome = append(sfHome, strings.Replace(detailMapRes["guestTeamName"].String(), " ", "", -1))
		} else if lotname == "半全彩" {
			sfHome = append(sfHome, strings.Replace(detailMapRes["masterTeamName"].String(), " ", "", -1))
		}
		sfGuest = append(sfGuest, strings.Replace(detailMapRes["guestTeamName"].String(), " ", "", -1))
	}

	LotteryZuCai.HomeTeamView = strings.Join(sfHome, ",")
	LotteryZuCai.AwayTeamView = strings.Join(sfGuest, ",")

	for _, detail := range sfcDetails {
		detailMap := detail.Map()
		levelTemp := detailMap["prizeLevel"].String()
		allmoneyTemp := strings.Replace(detailMap["stakeAmount"].String(), ",", "", -1)
		moneyTemp := strings.Replace(detailMap["stakeAmount"].String(), ",", "", -1)
		pieceTemp := strings.Replace(detailMap["stakeCount"].String(), ",", "", -1)
		level = append(level, levelTemp)
		allmoney = append(allmoney, allmoneyTemp)
		money = append(money, moneyTemp)
		piece = append(piece, pieceTemp)
	}

	LotteryZuCai.Level = strings.Join(level, ",")
	LotteryZuCai.Allmoney = strings.Join(allmoney, ",")
	LotteryZuCai.Money = strings.Join(money, ",")
	LotteryZuCai.Piece = strings.Join(piece, "_")

	return LotteryZuCai
}
