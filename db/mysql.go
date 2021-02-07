package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/core"
	"xorm.io/xorm"
	logs "xorm.io/xorm/log"
)

var Engine *xorm.Engine

//实例化数据库引擎
func newDb() {
	var err error
	var DbConf map[string]string = map[string]string{
		"user":      "hntcai",
		"pwd":       "hntcai0831",
		"host":      "192.168.1.70",
		"databases": "hntc",
		"charset":   "utf8",
		"port":      "1234",
	}

	var dbConn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", DbConf["user"], DbConf["pwd"], DbConf["host"], DbConf["port"], DbConf["databases"], DbConf["charset"])
	Engine, err = xorm.NewEngine("mysql", dbConn)
	if err != nil {
		panic("链接数据库失败")
	}

	Engine.ShowSQL(true)
	var path string = "mysqls.log"
	var newLoggerHander *FileHandler = NewFileHandler(path)
	//logger := xorm.NewSimpleLogger(logWriter)
	logger := logs.NewSimpleLogger(newLoggerHander)
	Engine.SetLogger(logger)
	//Engine.Logger().SetLevel(core.LOG_INFO)
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "wp_")
	Engine.SetMapper(tbMapper)
	//Engine.SetMaxIdleConns(10)
	//Engine.SetMaxOpenConns(1000)
	err = Engine.Ping()
	if err != nil {
		panic("链接数据库失败")
	}
	//设置数据库链接心跳
	//go dbHeartbeat(Engine)
}

//设置一个心跳 来保证数据库长时间的链接有新鲜

func dbHeartbeat(engine *xorm.Engine) {
	if engine == nil {
		return
	}

	var timer *time.Timer = time.NewTimer(3 * time.Minute)
	for {
		select {
		case <-timer.C:
			var err error = engine.Ping()
			if err != nil {
				fmt.Println("数据库心跳，连接数据库失败")
			}
			timer.Reset(3 * time.Minute)
		}

	}

}

//实例化操作
func init() {
	newDb()
}

/*"go.formatTool": "goimports",
"editor.quickSuggestions": null,
"go.useLanguageServer": false,
"go.useCodeSnippetsOnFunctionSuggest": true*/
