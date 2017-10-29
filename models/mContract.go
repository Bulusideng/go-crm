package models

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"

	"github.com/astaxie/beego/orm"
)

//合同号 序号	客户姓名 国家 项目 咨询 文案 转案日期 目前状态

type Contract struct {
	Contract_id   string `orm:"pk"`
	Seq           string
	Client_name   string
	Country       string
	Project_type  string
	Consulter     string
	Secretary     string
	Zhuan_an_date string
	Current_state string
}

func NewContract() *Contract {
	return &Contract{
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
	}
}

type Selector struct {
	Filter string
	Values map[string]bool
}

type ContractSelector struct {
	Contract_id   Selector
	Seq           Selector
	Client_name   Selector
	Country       Selector
	Project_type  Selector
	Consulter     Selector
	Secretary     Selector
	Zhuan_an_date Selector
	Current_state Selector
}

func NewContractSelector() *ContractSelector {
	return &ContractSelector{
		Contract_id:   Selector{"ALL", map[string]bool{}},
		Seq:           Selector{"ALL", map[string]bool{}},
		Client_name:   Selector{"ALL", map[string]bool{}},
		Country:       Selector{"ALL", map[string]bool{}},
		Project_type:  Selector{"ALL", map[string]bool{}},
		Consulter:     Selector{"ALL", map[string]bool{}},
		Secretary:     Selector{"ALL", map[string]bool{}},
		Zhuan_an_date: Selector{"ALL", map[string]bool{}},
		Current_state: Selector{"ALL", map[string]bool{}},
	}
}

func AddContract(c *Contract) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		fmt.Printf("Add Contract failed:%s %+v\n", err.Error(), *c)
		return err
	} else {
		fmt.Printf("Add Contract success %+v\n", *c)
	}
	return nil
}

func DelContract(c *Contract) error {
	o := orm.NewOrm()
	_, err := o.Delete(c)
	return err
}

func GetContract(contractId string) (c *Contract, err error) {
	o := orm.NewOrm()
	c = &Contract{
		Contract_id: contractId,
	}
	err = o.Read(c)
	return c, err
}

func UpdateContract(c *Contract) error {
	o := orm.NewOrm()
	old := *c
	err := o.Read(&old)
	if err == nil {
		if *c != old { //compae and update reflect.DeepEqual()
			_, err = o.Update(c)
		}
	}
	return err
}

func GetContracts(filters map[string]string) ([]*Contract, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("Contract")
	qs.Limit(-1)
	for k, v := range filters {
		qs = qs.Filter(strings.ToLower(k), v)
	}
	contracts := []*Contract{}
	_, err := qs.OrderBy("contract_id").All(&contracts)
	return contracts, err
}

func GetAllContracts() ([]*Contract, error) {
	o := orm.NewOrm()

	contracts := make([]*Contract, 0)

	qs := o.QueryTable("Contract")
	_, err := qs.Limit(-1).OrderBy("contract_id").All(&contracts)
	return contracts, err
}

func RegCase() {
	excelFileName := "./jjl.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {

	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) == 9 {
				c := &Contract{
					Seq:           row.Cells[0].String(),
					Contract_id:   row.Cells[1].String(),
					Client_name:   row.Cells[2].String(),
					Country:       row.Cells[3].String(),
					Project_type:  row.Cells[4].String(),
					Consulter:     row.Cells[5].String(),
					Secretary:     row.Cells[6].String(),
					Zhuan_an_date: row.Cells[7].String(),
					Current_state: row.Cells[8].String(),
				}
				AddContract(c)

			} else {
				for _, cell := range row.Cells {
					fmt.Printf("%s\n", cell.String())
				}
			}
		}
	}
}
