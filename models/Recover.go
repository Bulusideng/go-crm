package models

import (
	"errors"
	"fmt"
	"path"
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

var excelFileName = "./customer.xlsx"

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
			Uname: uname,
			Cname: "管理员",
			Title: Titles[EAdmin],
			//Email:  "364718765@qq.com",
			Pwd:    pwd,
			Status: "Active",
			ErrCnt: 0,
		}
		AddAccount(usr)

	}

}

var usrs = []DefUsr{}

type DefUsr struct {
	uname   string
	cname   string
	title   string
	manager string
}

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
			acct.Pwd = "1"
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
		sname := c.Country + c.Project_type
		if len(sname) > 30 {
			beego.Warn("Cut sheet name from: ", sname, " to: ", sname[:30])
			sname = sname[:30]
		}
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

		row.WriteStruct(c, -1)
		cid, _ := strconv.ParseInt(c.Contract_id, 10, 64)
		row.Cells[1].SetInt64(cid)
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
	legacy := false
	legacy, _ = beego.AppConfig.Bool("old_form")
	beego.BeeLogger.Warn("Old format:%t", legacy)
	contracts, err := GetAllContracts()
	if err == nil && len(contracts) > 0 {
		return
	}

	xlFile, err := xlsx.OpenFile(excelFileName)

	if err != nil {
		beego.Error("invalid file name")
		return
	}
	for sheetId, sheet := range xlFile.Sheets {
		rowId := 0
		for rowId = 1; rowId < len(sheet.Rows); rowId++ {
			row := sheet.Rows[rowId]

			if len(row.Cells) < 12 {
				fmt.Printf("IMPORT_CONTRACT_FAILED, Missing field in sheet:%s row: %d, cell cnt:%d: .", sheet.Name, rowId, len(row.Cells))
				for _, cell := range row.Cells {
					fmt.Printf("%s; ", cell.String())
				}
				fmt.Println()
				break
			}
			c := &Contract{}
			if !legacy {
				row.ReadStruct(c)
			} else {
				c = readOldForm(row)
			}

			if c.Contract_id != "" {
				if c.Create_by == "" {
					c.Create_by = "liupan"
				}
				if c.Create_date == "" {
					c.Create_date = time.Now().Format("2006-01-02")
				}
				c.Create_date = getDate(c.Create_date)
				if err := AddContract(c); err != nil {
					fmt.Printf("IMPORT_CONTRACT_FAILED:%s, sheet:%s row:%d. %+v\n", err.Error(), sheet.Name, rowId, *c)
				} else {
					//fmt.Printf("IMPORT_CONTRACT_SUCCESS: sheet:%s row:%d. %+v\n", sheet.Name, rowId, *c)
				}
			} else {
				fmt.Printf("IMPORT_CONTRACT_FAILED, Empty contract id, sheet:%s row: %d, parsed:%+v\n", sheet.Name, rowId, *c)
			}
		}
		fmt.Printf("Sheet:%d %s contracts: %d\n", sheetId, sheet.Name, rowId)
	}
}

func readOldForm(row *xlsx.Row) (c *Contract) {

	idx := 0
	c.Seq = strings.TrimSpace(row.Cells[idx].String())
	idx += 1
	c.Contract_id = strings.TrimSpace(row.Cells[1].String())
	idx += 1
	c.Client_name = strings.TrimSpace(row.Cells[idx].String())
	idx += 1
	c.Country = strings.TrimSpace(row.Cells[idx].String())
	idx += 1
	c.Project_type = strings.TrimSpace(row.Cells[idx].String())
	idx += 1

	idx += 1 //备注

	c.Contract_date = getDateCell(row.Cells[idx])
	idx += 1
	c.Consulters = strings.TrimSpace(row.Cells[idx].String())
	idx += 1

	idx += 1 //合同明细

	c.Secretaries = strings.TrimSpace(row.Cells[idx].String())
	idx += 1
	c.Zhuan_an_date = getDateCell(row.Cells[idx])
	idx += 1
	c.Current_state = getDateCell(row.Cells[idx])
	idx += 1

	//Optional fields
	if len(row.Cells) > idx {
		c.Didang_date = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Danganhao_date = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Buliao_date = getDateCell(row.Cells[idx])
	}
	idx += 1

	idx += 1 //面试排队

	if len(row.Cells) > idx {
		c.Interview_date1 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Interview_date2 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Pay_date1 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Pay_date2 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Nominate_date = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Federal_date1 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Federal_date2 = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Physical_date = getDateCell(row.Cells[idx])
	}
	idx += 1

	if len(row.Cells) > idx {
		c.Visa_date = getDateCell(row.Cells[idx])
	}
	idx += 1
	if len(row.Cells) > idx {
		c.Fail_date = getDateCell(row.Cells[idx])
	}
	return
}

func getDateCell(cel *xlsx.Cell) string {
	return getDate(cel.String())
}

func getDate(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, ".", "-", -1)
	strs := strings.Split(str, "-")
	if len(strs) != 3 { //YYYY-MM-DD
		beego.BeeLogger.Warn("Invalid date format: %s", str)
		return ""
	}
	if len(strs[1]) == 1 {
		strs[1] = "0" + strs[1]
	}
	if len(strs[2]) == 1 {
		strs[2] = "0" + strs[2]
	}
	return strs[0] + "-" + strs[1] + "-" + strs[2]
}
