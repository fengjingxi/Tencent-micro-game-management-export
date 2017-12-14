package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"net/url"
	"strconv"
	"github.com/tealeg/xlsx"
)

type game struct {
	IAppId        int    `json:"iAppId"`
	SzAppName     string `json:"szAppName"`
	SzTypeName    string `json:"szTypeName"`
	SzLastVersion string `json:"szLastVersion"`
	DtUpdateTime  string `json:"dtUpdateTime"`
	SzIconUrl     string `json:"szIconUrl"`
	SzDownloadUrl string `json:"szDownloadUrl"`
	SzCdnUrl      string `json:"szCdnUrl"`
	DtStartTime   string `json:"dtStartTime"`
	DtEndTime     string `json:"dtEndTime"`
	DtSignedTime  string `json:"dtSignedTime"`
	BIsBrokenRule int    `json:"bIsBrokenRule"`
	BIsNewVer     int    `json:"bIsNewVer"`
}

type channels struct {
	IChannelID     int    `json:"iChannelID"`
	SzChannelName  string `json:"szChannelName"`
	Games          []game `json:"games"`
	InvalidGameIDs []int  `json:"invalidGameIDs"`
}

type result struct {
	Channels []channels `json:"channels"`
}

type resp struct {
	Ret    int    `json:"ret"`
	Msg    string `json:"msg"`
	Result result `json:"result"`
}

type data struct {
	DtChannelOnlineTime string `json:"dtChannelOnlineTime"`
	DtCreateTime        string `json:"dtCreateTime"`
	DtLastUpdateTime    string `json:"dtLastUpdateTime"`
	DtOnlineTime        string `json:"dtOnlineTime"`
	IAppID              int    `json:"iAppID"`
	IStatus             int    `json:"iStatus"`
	SzAppDetail         string `json:"szAppDetail"`
	SzAppIconUrl        string `json:"szAppIconUrl"`
	SzAppName           string `json:"szAppName"`
	SzAppType           string `json:"szAppType"`
	SzLastVersion       string `json:"szLastVersion"`
	SzVerDesc           string `json:"szVerDesc"`
}

type resultdetails struct {
	GameInfo data `json:"gameInfo"`
}

type gamedetails struct {
	Msg    string        `json:"msg"`
	Result resultdetails `json:"result"`
	Ret    int           `json:"ret"`
}

type config struct {
	cookie     string
	name       string
	iChannelID string
}

func main() {
	config := config{
		cookie:     "RK=zH9eob16GS; _qpsvr_localtk=0.21118394657969475; pgv_pvi=1420285952; pgv_si=s580742144; ptui_loginuin=2880997838; ptisp=ctc; ptcz=30176a2feaae43b130ff1d3266e0537af749354108426df819640dfa5a95448f; uin=o2880997838; skey=@UmeQnuT8W; pt2gguin=o2880997838; isuser=1; key=2880997838-88d7ed90-dfd5-11e7-964e-d39c5844c4b4; iCPID=488; isregist=1; isaudit=1; issign=1; ilevel=0; iHead=http%3A//thirdqq.qlogo.cn/g%3Fb%3Dsdk%26k%3DfgFdOq0Tzuwt6micsbQF7bA%26s%3D140%26t%3D1493197858",
		name:       "file.xlsx",
		iChannelID: "10024328",
	}
	v := &resp{} //游戏管理列表
	d := &gamedetails{} //游戏详情
	client := &http.Client{}
	//post请求
	reqest, _ := http.NewRequest("POST", "http://s.qq.com/service/managecenter/gamemanage", nil)
	//设置header头，设置Cookie
	reqest.Header.Set("Cookie", config.cookie)
	response, _ := client.Do(reqest)//发送
	//是否成功
	if response.StatusCode == 200 {
		//如果成功
		body, _ := ioutil.ReadAll(response.Body) //从一个connection中读出数据
		bodystr := string(body) //转行成字符串、
		//字符串转json
		bs := []byte(bodystr)
		json.Unmarshal(bs, v)
		//准备循环游戏列表
		gameid := v.Result.Channels[0].Games
		//准备导出Excel
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("Sheet1")
		//设置第一行的名字
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		cell := row.AddCell()
		cell.Value = "appid"
		cell = row.AddCell()
		cell.Value = "游戏名称"
		cell = row.AddCell()
		cell.Value = "游戏类型"
		cell = row.AddCell()
		cell.Value = "icon"
		cell = row.AddCell()
		cell.Value = "版本号"
		cell = row.AddCell()
		cell.Value = "更新时间"
		cell = row.AddCell()
		cell.Value = "说明"
		//开始循环游戏列表
		for _, v := range gameid {
			//准备post参数
			postValues := url.Values{}
			postValues.Set("iAppID", strconv.Itoa(v.IAppId))
			postValues.Set("iChannelID", config.iChannelID)
			postDataStr := postValues.Encode()
			postDataBytes := []byte(postDataStr)
			postBytesReader := bytes.NewReader(postDataBytes)
			//post请求
			reqest_details, _ := http.NewRequest("POST", "http://s.qq.com/service/managecenter/gamemanage/gamedetail", postBytesReader)
			//这个一定要加，不加form的值post不过去，被坑了两小时
			reqest_details.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			reqest_details.Header.Set("Cookie", config.cookie)
			response_details, _ := client.Do(reqest_details)//发送
			if response_details.StatusCode == 200 {
				bodyDetails, _ := ioutil.ReadAll(response_details.Body)
				bbs := []byte(bodyDetails)
				json.Unmarshal(bbs, d)
				gameinfo := d.Result.GameInfo
				row := sheet.AddRow()
				row.SetHeightCM(1) //设置每行的高度
				cell := row.AddCell()
				cell.Value = strconv.Itoa(gameinfo.IAppID)
				cell = row.AddCell()
				cell.Value = gameinfo.SzAppName
				cell = row.AddCell()
				cell.Value = gameinfo.SzAppType
				cell = row.AddCell()
				cell.Value = gameinfo.SzAppIconUrl
				cell = row.AddCell()
				cell.Value = gameinfo.SzLastVersion
				cell = row.AddCell()
				cell.Value = gameinfo.DtLastUpdateTime
				cell = row.AddCell()
				cell.Value = gameinfo.SzVerDesc

			}
		}
		//设置游戏名称
		err := file.Save(config.name)
		if err != nil {
			panic(err)
		}

	}

}
