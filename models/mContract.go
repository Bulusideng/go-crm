package models

import (
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//合同明细
//序号	合同号	客户姓名	客户电话 申请国家	项目	签约日期	顾问	 文案姓名 转案日期	 状态（案件进度）
//递档日期	档案号	补料	 通知面试	面试	打款通知	打款确认	省提名	递交联邦	联邦档案号	通知体检	获签	拒签

type Contract struct {
	Seq             string `xlsx:"0"`          //序号
	Contract_id     string `xlsx:"1" orm:"pk"` //合同号
	Client_name     string `xlsx:"2"`          //客户姓名
	Client_tel      string `xlsx:"3"`          //客户电话
	Country         string `xlsx:"4"`          //申请国家
	Project_type    string `xlsx:"5"`          //项目
	Contract_date   string `xlsx:"6"`          //签约日期
	Consulters      string `xlsx:"7"`          //顾问
	Secretaries     string `xlsx:"8"`          //文案
	Zhuan_an_date   string `xlsx:"9"`          //转案日期
	Current_state   string `xlsx:"10"`         //状态（案件进度）
	Didang_date     string `xlsx:"11"`         //递档日期
	Danganhao_date  string `xlsx:"12"`         //档案号
	Buliao_date     string `xlsx:"13"`         //补料
	Interview_date1 string `xlsx:"14"`         //通知面试
	Interview_date2 string `xlsx:"15"`         //面试
	Pay_date1       string `xlsx:"16"`         //打款通知
	Pay_date2       string `xlsx:"17"`         //打款确认
	Nominate_date   string `xlsx:"18"`         //省提名
	Federal_date1   string `xlsx:"19"`         //递交联邦
	Federal_date2   string `xlsx:"20"`         //联邦档案号
	Physical_date   string `xlsx:"21"`         //通知体检
	Visa_date       string `xlsx:"22"`         //获签
	Fail_date       string `xlsx:"23"`         //拒签
	Create_date     string `xlsx:"24"`         //录入时间
	Create_by       string `xlsx:"25"`         //录入人
}

func NewAllFilter() *Contract {
	return &Contract{
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
		"ALL",
	}
}

type ContractSelector struct {
	Contract_id   FieldSelector
	Seq           FieldSelector
	Client_name   FieldSelector
	Client_tel    FieldSelector
	Country       FieldSelector
	Project_type  FieldSelector
	Contract_date FieldSelector
	Consulters    FieldSelector
	Secretaries   FieldSelector
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
		Contract_date: FieldSelector{"ALL", map[string]bool{}},
		Consulters:    FieldSelector{"ALL", map[string]bool{}},
		Secretaries:   FieldSelector{"ALL", map[string]bool{}},
		Create_date:   FieldSelector{"ALL", map[string]bool{}},
		Create_by:     FieldSelector{"ALL", map[string]bool{}},
		Zhuan_an_date: FieldSelector{"ALL", map[string]bool{}},
		Current_state: FieldSelector{"ALL", map[string]bool{}},
	}
}

func AddContract(c *Contract) error {
	if c.Consulters == "" {
		c.Consulters = "空缺"
	}
	if c.Secretaries == "" {
		c.Secretaries = "空缺"
	}
	if c.Zhuan_an_date == "" {
		c.Zhuan_an_date = "空缺"
	}
	o := orm.NewOrm()
	_, err := o.Insert(c)

	if err != nil {
		p, err1 := GetContract(c.Contract_id)
		if err1 == nil {
			beego.Error("Previous contract: ", *p)
		}
		//fmt.Printf("Add Contract failed:%s %+v\n", err.Error(), *c)
		return err
	} else {
		//fmt.Printf("Add Contract success %+v\n", *c)
	}
	return nil
}

func DelContract(cid string) error {
	c := &Contract{}
	c.Contract_id = cid
	o := orm.NewOrm()
	_, err := o.Delete(c)
	return err
}

func GetContract(cid string) (c *Contract, err error) {
	o := orm.NewOrm()
	c = &Contract{}
	c.Contract_id = cid
	err = o.Read(c)
	if err != nil {
		beego.Error("GetContract[%s]", cid, " failed: ", err.Error())
	}
	return c, err
}

func UpdateContract(oldContractId string, c *Contract) (*ChangeSlice, error) {
	o := orm.NewOrm()
	old := &Contract{
		Contract_id: oldContractId,
	}
	err := o.Read(old)
	if c.Zhuan_an_date == "" {
		c.Zhuan_an_date = "空缺"
	}
	changes := ChangeSlice{}
	if err == nil {
		if !reflect.DeepEqual(*c, old) { //compae and update reflect.DeepEqual()
			if c.Contract_id != old.Contract_id {
				changes = append(changes, Change{Item: "合同号", Last: old.Contract_id, Current: c.Contract_id})
			}
			if c.Client_name != old.Client_name {
				changes = append(changes, Change{Item: "客户姓名", Last: old.Client_name, Current: c.Client_name})
			}
			if c.Client_tel != old.Client_tel {
				changes = append(changes, Change{Item: "客户电话", Last: old.Client_tel, Current: c.Client_tel})
			}
			if c.Country != old.Country {
				changes = append(changes, Change{Item: "国家", Last: old.Country, Current: c.Country})
			}
			if c.Project_type != old.Project_type {
				changes = append(changes, Change{Item: "项目", Last: old.Project_type, Current: c.Project_type})
			}
			if c.Contract_date != old.Contract_date {
				changes = append(changes, Change{Item: "签约日期", Last: old.Contract_date, Current: c.Contract_date})
			}
			if c.Consulters != old.Consulters {
				changes = append(changes, Change{Item: "顾问", Last: old.Consulters, Current: c.Consulters})
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
			if c.Didang_date != old.Didang_date {
				changes = append(changes, Change{Item: "递档", Last: old.Didang_date, Current: c.Didang_date})
			}
			if c.Danganhao_date != old.Danganhao_date {
				changes = append(changes, Change{Item: "档案号", Last: old.Danganhao_date, Current: c.Danganhao_date})
			}

			if c.Buliao_date != old.Buliao_date {
				changes = append(changes, Change{Item: "补料", Last: old.Buliao_date, Current: c.Buliao_date})
			}
			if c.Interview_date1 != old.Interview_date1 {
				changes = append(changes, Change{Item: "通知面试", Last: old.Interview_date1, Current: c.Interview_date1})
			}
			if c.Interview_date2 != old.Interview_date2 {
				changes = append(changes, Change{Item: "面试", Last: old.Interview_date2, Current: c.Interview_date2})
			}
			if c.Pay_date1 != old.Pay_date1 {
				changes = append(changes, Change{Item: "打款通知", Last: old.Pay_date1, Current: c.Pay_date1})
			}
			if c.Pay_date2 != old.Pay_date2 {
				changes = append(changes, Change{Item: "打款确认", Last: old.Pay_date2, Current: c.Pay_date2})
			}
			if c.Nominate_date != old.Nominate_date {
				changes = append(changes, Change{Item: "省提名", Last: old.Nominate_date, Current: c.Nominate_date})
			}

			if c.Federal_date1 != old.Federal_date1 {
				changes = append(changes, Change{Item: "递交联邦", Last: old.Federal_date1, Current: c.Federal_date1})
			}
			if c.Federal_date2 != old.Federal_date2 {
				changes = append(changes, Change{Item: "联邦档案号", Last: old.Federal_date2, Current: c.Federal_date2})
			}
			if c.Physical_date != old.Physical_date {
				changes = append(changes, Change{Item: "通知体检", Last: old.Physical_date, Current: c.Physical_date})
			}
			if c.Visa_date != old.Visa_date {
				changes = append(changes, Change{Item: "获签", Last: old.Visa_date, Current: c.Visa_date})
			}
			if c.Fail_date != old.Fail_date {
				changes = append(changes, Change{Item: "拒签", Last: old.Fail_date, Current: c.Fail_date})
			}
			if c.Create_date != old.Create_date {
				changes = append(changes, Change{Item: "录入日期", Last: old.Create_date, Current: c.Create_date})
			}
			if c.Create_by != old.Create_by {
				changes = append(changes, Change{Item: "录入人", Last: old.Create_by, Current: c.Create_by})
			}
			if old.Contract_id == c.Contract_id {
				_, err = o.Update(c)
			} else {
				_, err = o.Insert(c)
				if err == nil {
					_, err = o.Delete(old)
					comments, _ := GetComments(old.Contract_id) //Update contract id for comments
					for _, comment := range comments {
						comment.Contract_id = c.Contract_id
						UpdateComment(comment)
					}
				}
			}
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
