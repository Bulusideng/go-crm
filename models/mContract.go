package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

//合同明细
//序号	合同号	客户姓名	客户电话 申请国家	项目	签约日期	顾问	 文案姓名 转案日期	 状态（案件进度）
//递档日期	档案号	补料	 通知面试	面试	打款通知	打款确认	省提名	递交联邦	联邦档案号	通知体检	获签	拒签

type Contract struct {
	Seq             string //序号
	Contract_id     string `orm:"pk"` //合同号
	Client_name     string //客户姓名
	Client_tel      string //客户电话
	Country         string //申请国家
	Project_type    string //项目
	Contract_date   string //签约日期
	Consulters      string //顾问
	Secretaries     string //文案
	Zhuan_an_date   string //转案日期
	Current_state   string //状态（案件进度）
	Didang_date     string //递档日期
	Danganhao_date  string //档案号
	Buliao_date     string //补料
	Interview_date1 string //通知面试
	Interview_date2 string //面试
	Pay_date1       string //打款通知
	Pay_date2       string //打款确认
	Nominate_date   string //省提名
	Federal_date1   string //递交联邦
	Federal_date2   string //联邦档案号
	Physical_date   string //通知体检
	Visa_date       string //获签
	Fail_date       string //拒签

	Create_date string
	Create_by   string
}

type ContractSelector struct {
	Contract_id   FieldSelector
	Seq           FieldSelector
	Client_name   FieldSelector
	Client_tel    FieldSelector
	Country       FieldSelector
	Project_type  FieldSelector
	Consulter     FieldSelector
	Secretary     FieldSelector
	Create_date   FieldSelector
	Create_by     FieldSelector
	Zhuan_an_date FieldSelector
	Current_state FieldSelector
}
type FieldSelector struct {
	CurSelected string          //Filter value of field
	List        map[string]bool //All values of the field
}

func NewContractSelectors() *ContractSelector {
	return &ContractSelector{
		Contract_id:   FieldSelector{"ALL", map[string]bool{}},
		Seq:           FieldSelector{"ALL", map[string]bool{}},
		Client_name:   FieldSelector{"ALL", map[string]bool{}},
		Client_tel:    FieldSelector{"ALL", map[string]bool{}},
		Country:       FieldSelector{"ALL", map[string]bool{}},
		Project_type:  FieldSelector{"ALL", map[string]bool{}},
		Consulter:     FieldSelector{"ALL", map[string]bool{}},
		Secretary:     FieldSelector{"ALL", map[string]bool{}},
		Create_date:   FieldSelector{"ALL", map[string]bool{}},
		Create_by:     FieldSelector{"ALL", map[string]bool{}},
		Zhuan_an_date: FieldSelector{"ALL", map[string]bool{}},
		Current_state: FieldSelector{"ALL", map[string]bool{}},
	}
}

func AddContract(c *Contract) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		fmt.Printf("Add Contract failed:%s %+v\n", err.Error(), *c)
		return err
	} else {
		//fmt.Printf("Add Contract success %+v\n", *c)
	}
	return nil
}

func DelContract(cid string) error {
	c := &Contract{
		Contract_id: cid,
	}
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
	if err != nil {
		fmt.Printf("GetContract[%s] failed %+v\n", contractId, err.Error())
	}
	return c, err
}

func UpdateContract(c *Contract) (*ChangeSlice, error) {

	o := orm.NewOrm()
	old := *c
	err := o.Read(&old)
	changes := ChangeSlice{}
	if err == nil {
		if !reflect.DeepEqual(*c, old) { //compae and update reflect.DeepEqual()
			if c.Client_name != old.Client_name {
				changes = append(changes, Change{Item: "客户姓名", Last: old.Client_name, Current: c.Client_name})
			}
			if c.Client_tel != old.Client_tel {
				changes = append(changes, Change{Item: "客户电话", Last: old.Client_tel, Current: c.Client_tel})
			}
			if c.Country != old.Country {
				changes = append(changes, Change{Item: "国家", Last: old.Country, Current: c.Country})
			}
			if c.Consulters != old.Consulters {
				changes = append(changes, Change{Item: "咨询", Last: old.Consulters, Current: c.Consulters})
			}
			if c.Secretaries != old.Secretaries {
				changes = append(changes, Change{Item: "文案", Last: old.Secretaries, Current: c.Secretaries})
			}
			if c.Zhuan_an_date != old.Zhuan_an_date {
				changes = append(changes, Change{Item: "转案日期", Last: old.Zhuan_an_date, Current: c.Zhuan_an_date})
			}
			if c.Current_state != old.Current_state {
				changes = append(changes, Change{Item: "状态", Last: old.Current_state, Current: c.Current_state})
			}
			_, err = o.Update(c)
		}
	}
	if err != nil {
		return nil, err
	}
	return &changes, nil
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

func NewContract() *Contract {
	return &Contract{}
	/*
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
	*/

}
