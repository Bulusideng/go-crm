package models

import (
	"errors"
	"path"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

func Init() {
	RegisterAccounts()
	ImportContracts()

}

func RegisterAdmin() {
	users, err := GetAcctounts(map[string]string{"title": "Admin"})
	for err != nil {
		time.Sleep(time.Second)

		users, err = GetAllAccts()
	}
	if err == nil && len(users) == 0 { ///No user yet, add the user in config to DB
		uname := beego.AppConfig.String("uname")
		pwd := beego.AppConfig.String("pwd")
		usr := &Account{
			Uname:  uname,
			Cname:  "管理员",
			Title:  Titles[EAdmin],
			Email:  "364718765@qq.com",
			Pwd:    pwd,
			Status: "Active",
			ErrCnt: 0,
		}
		AddAccount(usr)

	}

}

type DefUsr struct {
	uname   string
	cname   string
	title   string
	manager string
}

var usrs = []DefUsr{}

func RegisterAccounts() {
	RegisterAdmin()
	acct := &Account{
		Pwd:    "1",
		Status: "Active",
		ErrCnt: 0,
	}
	users, _ := GetAllAccts()
	if len(users) < 3 {
		for _, u := range usrs {
			acct.Uname = u.uname
			acct.Cname = u.cname
			acct.Title = u.title
			acct.Manager = u.manager
			AddAccount(acct)
		}
	}

}

func ExportContracts() (string, error) {
	contracts, err := GetAllContracts()
	if err != nil {
		return "", err
	}
	if len(contracts) == 0 {
		return "", errors.New("没有客户信息")
	}

	file := xlsx.NewFile()

	sm := map[string]*xlsx.Sheet{}
	var sheet *xlsx.Sheet
	var ok bool
	for _, c := range contracts {
		sname := c.Country //+ c.Project_type
		if sheet, ok = sm[sname]; !ok {
			sm[sname], err = file.AddSheet(sname)
			if err != nil {
				beego.Error("Add sheet: " + sname + "panic, error:" + err.Error())
				return "", err
			}
			row := sm[sname].AddRow()
			row.SetHeightCM(1) //设置每行的高度
			cell := row.AddCell()
			cell.Value = "序号"
			cell = row.AddCell()
			cell.Value = "合同号"
			cell = row.AddCell()
			cell.Value = "客户姓名"
			cell = row.AddCell()
			cell.Value = "客户电话"
			cell = row.AddCell()
			cell.Value = "申请国家"
			cell = row.AddCell()
			cell.Value = "项目"
			cell = row.AddCell()
			cell.Value = "签约日期"
			cell = row.AddCell()
			cell.Value = "顾问"
			cell = row.AddCell()
			cell.Value = "文案"
			cell = row.AddCell()
			cell.Value = "转案日期"
			cell = row.AddCell()
			cell.Value = "状态"
			cell = row.AddCell()
			cell.Value = "递档日期"
			cell = row.AddCell()
			cell.Value = "档案号"
			cell = row.AddCell()
			cell.Value = "补料"
			cell = row.AddCell()
			cell.Value = "通知面试"
			cell = row.AddCell()
			cell.Value = "面试"
			cell = row.AddCell()
			cell.Value = "打款通知"
			cell = row.AddCell()
			cell.Value = "打款确认"
			cell = row.AddCell()
			cell.Value = "省提名"
			cell = row.AddCell()
			cell.Value = "递交联邦"
			cell = row.AddCell()
			cell.Value = "联邦档案号"
			cell = row.AddCell()
			cell.Value = "通知体检"
			cell = row.AddCell()
			cell.Value = "获签"
			cell = row.AddCell()
			cell.Value = "拒签"
			cell = row.AddCell()
			cell.Value = "录入时间"
			cell = row.AddCell()
			cell.Value = "录入人"

		}
		sheet = sm[sname]
		row := sheet.AddRow()

		//row.SetHeightCM(1)
		row.WriteStruct(c, 26)
	}

	ts := strings.Replace(time.Now().Format(time.RFC1123), ":", "-", -1)

	fn := "客户备份" + ts + ".xlsx"
	err = file.Save(path.Join(`files`, fn))
	if err != nil {
		return "", errors.New("保存备份文件失败: " + err.Error())
	} else {
		return fn, nil
	}
}

func ImportContracts() {
	contracts, err := GetAllContracts()
	if err == nil && len(contracts) > 0 {
		return
	}
	excelFileName := "./customer.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)

	if err != nil {
		beego.Error("invalid file name")
		return
	}
	for i, sheet := range xlFile.Sheets {
		j := 0
		for j = 1; j < len(sheet.Rows); j++ {
			row := sheet.Rows[j]
			if len(row.Cells) >= 26 {
				c := &Contract{
					Seq: strings.TrimSpace(row.Cells[0].String()),

					Client_name:  strings.TrimSpace(row.Cells[2].String()),
					Country:      strings.TrimSpace(row.Cells[3].String()),
					Project_type: strings.TrimSpace(row.Cells[4].String()),
					//备注
					Contract_date: strings.TrimSpace(row.Cells[6].String()),
					Consulters:    strings.TrimSpace(row.Cells[7].String()),
					//合同明细
					Secretaries:    strings.TrimSpace(row.Cells[9].String()),
					Zhuan_an_date:  strings.TrimSpace(row.Cells[10].String()),
					Current_state:  strings.TrimSpace(row.Cells[11].String()),
					Didang_date:    strings.TrimSpace(row.Cells[12].String()),
					Danganhao_date: strings.TrimSpace(row.Cells[13].String()),
					Buliao_date:    strings.TrimSpace(row.Cells[14].String()),
					//面试排队
					Interview_date1: strings.TrimSpace(row.Cells[16].String()),
					Interview_date2: strings.TrimSpace(row.Cells[17].String()),
					Pay_date1:       strings.TrimSpace(row.Cells[18].String()),
					Pay_date2:       strings.TrimSpace(row.Cells[19].String()),
					Nominate_date:   strings.TrimSpace(row.Cells[20].String()),
					Federal_date1:   strings.TrimSpace(row.Cells[21].String()),
					Federal_date2:   strings.TrimSpace(row.Cells[22].String()),
					Physical_date:   strings.TrimSpace(row.Cells[23].String()),
					Visa_date:       strings.TrimSpace(row.Cells[24].String()),
					Fail_date:       strings.TrimSpace(row.Cells[25].String()),
				}
				if len(row.Cells) >= 28 {
					c.Create_date = strings.TrimSpace(row.Cells[26].String())
					c.Create_by = strings.TrimSpace(row.Cells[27].String())
				} else {
					c.Create_date = time.Now().Format(time.RFC1123)
					c.Create_by = "liupan"
				}
				c.Contract_id, _ = strconv.ParseInt(strings.TrimSpace(row.Cells[1].String()), 10, 64)
				if c.Contract_id > 0 {
					c.Create_by = "liupan"
					c.Create_date = time.Now().Format(time.RFC1123)
					if err := AddContract(c); err != nil {
						fmt.Printf("Create contract failed:%s, sheet:%s row:%d. %+v\n", err.Error(), sheet.Name, j, *c)
					}
				} else {
					fmt.Printf("Create contract failed, sheet: %s row: %d, parsed:%+v\n", sheet.Name, j, *c)
				}
			} else {
				fmt.Printf("Invalid contract in sheet:%d %s row: %d, cells:%d: .", i, sheet.Name, j, len(row.Cells))
				for _, cell := range row.Cells {
					fmt.Printf("%s; ", cell.String())
				}
				fmt.Println()
				break
			}
		}
		fmt.Printf("Sheet:%d %s contracts: %d\n", i, sheet.Name, j)
	}
}
