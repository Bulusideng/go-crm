package models

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

func Init() {
	RegisterAccounts()
	ImportContracts()

}

var (
	CommentSheet = "History"
)

//var usrs = []DefUsr{}

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

func ExportContracts(cat string, withComments bool) (string, error) {
	beego.Warn("cat:", cat)
	usr := &Account{
		Title: "Admin",
	}

	contracts, err := GetContracts(usr, nil)
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
		sname := "不分类"
		switch cat {
		case "国家":
			sname = c.Country
		case "项目":
			sname = c.Country + c.Project_type
		case "顾问":
			sname = c.Consulters
		case "文案":
			sname = c.Secretaries
		case "年份":
			sname = c.Contract_date
			if len(sname) >= 4 { //Get the year
				sname = sname[:4]
			}
		}
		if len(sname) == 0 {
			sname = "未知" + cat
		} else if len(sname) > 30 {
			beego.Warn("Cut sheet name from: ", sname, " to: ", sname[:30])
			sname = sname[:30]
		}
		sname = strings.Replace(sname, "/", ",", -1)
		if sheet, ok = sm[sname]; !ok {
			sm[sname], err = file.AddSheet(sname)
			beego.Error("Add sheet: " + sname)
			if err != nil {
				beego.Error("Add sheet: " + sname + "panic, error:" + err.Error())
				return "", err
			}
			row := sm[sname].AddRow()
			row.WriteSlice(&[]string{
				"序号", "合同号", "客户姓名", "客户电话", "申请国家", "项目", "签约日期",
				"顾问", "文案", "转案日期", "状态", "递档日期", "档案号", "补料", "通知面试",
				"面试", "打款通知", "打款确认", "省提名", "递交联邦", "联邦档案号", "通知体检",
				"获签", "拒签", "录入时间", "录入人", "状态"}, -1)
		}
		sheet = sm[sname]
		row := sheet.AddRow()
		row.WriteStruct(c, -1)
	}
	if withComments {
		history, _ := file.AddSheet(CommentSheet)

		row := history.AddRow()
		row.WriteSlice(&[]string{"Id", "合同号", "修改者职位", "用户名", "姓名", "日期", "改动", "内容"}, -1)
		comments, _ := GetComments("")
		for _, comment := range comments {
			history.AddRow().WriteStruct(comment, -1)
		}
	}

	ts := strings.Replace(time.Now().Format(time.RFC1123), ":", "-", -1)

	fn := "客户备份 " + cat + ts + ".xlsx"
	err = file.Save(path.Join(`files`, fn))
	if err != nil {
		return "", errors.New("保存备份文件失败: " + err.Error())
	} else {
		return fn, nil
	}
}

var (
	sheetName = ""
	rowId     = 0
)

func ImportContracts() {
	orig_format := false
	orig_format, _ = beego.AppConfig.Bool("orig_format")
	beego.BeeLogger.Warn("Old format:%t", orig_format)
	usr := &Account{
		Title: "Admin",
	}
	contracts, err := GetContracts(usr, nil)
	if err == nil && len(contracts) > 0 {
		return
	}

	var excelFileName = "./backup.xlsx"
	if orig_format {
		excelFileName = "./orig.xlsx"
	}

	xlFile, err := xlsx.OpenFile(excelFileName)
	total_cnt := 0
	if err != nil {
		beego.Error("invalid file name")
		return
	}
	for sheetId, sheet := range xlFile.Sheets {
		cnt := 0
		rowId = 0
		sheetName = sheet.Name
		for rowId = 1; rowId < len(sheet.Rows); rowId++ {
			row := sheet.Rows[rowId]
			if sheetName == CommentSheet {
				comment := &Comment{}
				row.ReadStruct(comment)
				AddComment(comment)
				continue
			}

			if len(row.Cells) < 12 {
				fmt.Printf("IMPORT_CONTRACT_FAILED, Missing field in sheet:%s row: %d, cell cnt:%d: .", sheet.Name, rowId, len(row.Cells))
				for _, cell := range row.Cells {
					fmt.Printf("%s; ", cell.String())
				}
				fmt.Println()
				break
			}
			c := &Contract{}
			if !orig_format {
				row.ReadStruct(c)
			} else {
				c, err = readOldForm(row)
				if err != nil {
					fmt.Printf("%s: %+v\n", err.Error(), row.Cells)
					break
				}
			}

			if c.Contract_id != "" {
				if c.Create_by == "" {
					c.Create_by = "liupan"
				}
				if c.Create_date == "" {
					c.Create_date = time.Now().Format("2006-01-02")
				}
				c.Create_date = getDate(c.Create_date)
				err := AddContract(c)
				if err != nil {
					fmt.Printf("IMPORT_CONTRACT_FAILED:%s, sheet:%s row:%d. %+v\n\n", err.Error(), sheet.Name, rowId, *c)
					c.Contract_id = c.Contract_id + "_dup"
					err = AddContract(c)

				} else {
					cnt++
					//fmt.Printf("IMPORT_CONTRACT_SUCCESS: sheet:%s row:%d. %+v\n", sheet.Name, rowId, *c)
				}
			} else {
				fmt.Printf("IMPORT_CONTRACT_FAILED, Empty contract id, sheet:%s row: %d, parsed:%+v\n", sheet.Name, rowId, *c)
			}
		}
		fmt.Printf("Sheet:%d %s rows: %d, contracts: %d\n", sheetId, sheet.Name, rowId, cnt)
		total_cnt += cnt
	}

	fmt.Printf("Total contracts: %d\n", total_cnt)
}

func readOldForm(row *xlsx.Row) (*Contract, error) {
	c := &Contract{}

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

	if c.Seq == "" && c.Contract_id == "" && c.Client_name == "" {
		return nil, errors.New("Invalid contract sheet[" + sheetName + "], row:" + string(rowId))
	}

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
	c.Current_state = strings.TrimSpace(row.Cells[idx].String())
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
	return c, nil
}

func getDateCell(cel *xlsx.Cell) string {
	return getDate(cel.String())
}

func getDate(str string) string {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return ""
	}
	str = strings.Replace(str, ".", "-", -1)
	strs := strings.Split(str, "-")
	if len(strs) < 3 { //YYYY-MM-DD
		beego.BeeLogger.Warn("Invalid date format: %s in sheet[%s], row:%d", str, sheetName, rowId)
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
