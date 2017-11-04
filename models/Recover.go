package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

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
		usr = &Account{
			Uname:  "dengx",
			Cname:  "管理员",
			Title:  Titles[EAdmin],
			Mobile: "15010275683",
			Email:  "364718765@qq.com",
			Pwd:    "dengx",
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

func ImportContracts() {
	contracts, err := GetAllContracts()
	if err == nil && len(contracts) > 0 {
		return
	}
	excelFileName := "./customer.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic("invalid file name")
	}
	for i, sheet := range xlFile.Sheets {
		j := 0
		for j = 1; j < len(sheet.Rows); j++ {
			row := sheet.Rows[j]
			if len(row.Cells) >= 26 {
				c := &Contract{
					Seq:          strings.TrimSpace(row.Cells[0].String()),
					Contract_id:  strings.TrimSpace(row.Cells[1].String()),
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
				if len(c.Contract_id) > 0 {
					c.Create_by = "liupan"
					c.Create_date = time.Now().Format(time.RFC1123)
					if err := AddContract(c); err != nil {
						fmt.Printf("Create contract failed, sheet:%d %s row: %d\n", i, sheet.Name, j)
					}
				} else {
					fmt.Printf("Invalid contract id:%+v\n", *c)
				}
			} else {
				fmt.Printf("Invalid contract in sheet:%d %s row: %d, cells:%d:\n", i, sheet.Name, j, len(row.Cells))
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
