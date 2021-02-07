package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func LotteryTimes(lotname string) bool {
	loc, _ := time.LoadLocation("Local")
	var nowTime time.Time = time.Now().In(loc)
	sendtime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-03-16 15:04:05", loc)

	if nowTime.After(sendtime) {
		return false
	}

	if lotname == "jclq" || lotname == "jczq" || lotname == "zcall" || lotname == "numpdf" {
		return true
	}

	var weeks string = nowTime.Weekday().String()
	var weekInt int = Weeks(weeks)
	if (weekInt == 2 || weekInt == 4 || weekInt == 5 || weekInt == 7) && lotname == "dlt" {
		return false
	}

	if (weekInt == 1 || weekInt == 3 || weekInt == 4 || weekInt == 6) && lotname == "qxc" {
		return false
	}

	var nowDate string = nowTime.Format("2006-01-02")
	var start string = nowDate + " 20:35:00"
	var end string = nowDate + " 23:55:00"
	starttime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, loc)
	endtime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, loc)
	var strartDiff bool = nowTime.After(starttime)
	var endDiff bool = nowTime.Before(endtime)

	if strartDiff && endDiff {
		return true
	}

	return false
}

func Weeks(week string) int {
	var nowWeet int = 0
	switch week {
	case "Monday":
		nowWeet = 1
	case "Tuesday":
		nowWeet = 2
	case "Wednesday":
		nowWeet = 3
	case "Thursday":
		nowWeet = 4
	case "Friday":
		nowWeet = 5
	case "Saturday":
		nowWeet = 6
	case "Sunday":
		nowWeet = 7
	default:
		nowWeet = 0
	}
	return nowWeet
}

func IsFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		return false
	}
	return false
}

func DownloadFile(url string, localPath string) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return err
	}

	if IsFileExist(localPath, fsize) {
		return err
	}
	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback

	}
	if err == nil {
		file.Close()
		os.Rename(tmpFilePath, localPath)
	}
	return err
}
