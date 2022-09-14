package controller

import (
	"IpTellServer/mydbs"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func reqAddrApi(ip string) string {
	clien := &http.Client{
		Timeout: time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://sp1.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=%s&resource_id=5809", ip), nil)
	checkErr(err)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Version/3.1 Safari/525.13")
	resq, err := clien.Do(req)
	checkErr(err)
	defer resq.Body.Close()
	body, err := ioutil.ReadAll(resq.Body)
	return string(body)
}

func GetAddr(c *gin.Context) {
	data, err := mydbs.DataSelect(c.RemoteIP())
	nowtime := int(time.Now().Unix())
	if err != nil {
		//未查询到数据，插入数据,并返回
		data := mydbs.Mydata{
			Myip:   c.RemoteIP(),
			Mytime: nowtime,
		}
		data.Myvalue = reqAddrApi(c.RemoteIP())
		mydbs.DataInsert(&data)
		c.String(200, data.Myvalue)
	} else if int(time.Now().Unix())-data.Mytime > 3600*24 {
		//时间超出一天,重新请求数据,并更新数据
		data.Myvalue = reqAddrApi(data.Myip)
		data.Mytime = nowtime
		mydbs.DataUpdate(&data)
		c.String(200, data.Myvalue)
	} else {
		//查询到数据,返回数据
		c.String(200, data.Myvalue)
	}
}

func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"RemoteIP": c.RemoteIP(),
		"ClientIP": c.ClientIP(),
	})
}

func GetIp(c *gin.Context) {
	c.JSON(200, gin.H{
		"origin": c.RemoteIP(),
	})
}

func MyGet(c *gin.Context) {

	c.JSON(200, gin.H{
		"headers": c.Request.Header,
	})
}
