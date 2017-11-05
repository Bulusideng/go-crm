package models

import (
	"errors"
)

const (
	EAdmin = iota //0
	EManager
	EConsulter
	ESecretary
	EGuest
)

var Titles = []string{
	"Admin",
	"Manager",
	"Consulter",
	"Secretary",
	"Guest",
}

func Guest() *Account {
	return &Account{
		Uname:  "Guest",
		Title:  "Guest",
		Cname:  "访客",
		Status: "Disabled",
	}
}

type Account struct {
	Uname      string `orm:"pk"`
	Title      string //Administrator, Manager, Consulter or Secretary
	Cname      string //Chinese name
	Manager    string
	Email      string
	Mobile     string
	Status     string //"Locked":锁定 	"Disabled": 禁用 	"Active"：正常
	ErrCnt     int    //3 times wll lock acct
	CreateDate string

	Pwd string //Encrypted
}

func (this *Account) IsAdmin() bool {
	return this.Title == "Admin"
}

func (this *Account) IsManager() bool {
	return this.Title == "Manager"
}

func (this *Account) IsWorker() bool {
	return this.Title == "Manager" || this.Title == "Consulter" || this.Title == "Secretary"
}

func (this *Account) IsManagerOf(u *Account) bool { //Admin is manager of all
	return this.IsAdmin() || (this.IsManager() && (u.Manager == "" || this.Uname == u.Manager))
}

func (this *Account) IsGuest() bool {
	return this.Title == "Guest"
}

//Set Status
func (this *Account) Lock() { //Lock by pwd error 3 times
	if !this.Disabled() { //Don't set disabled to lock
		this.Status = "Locked"
	}
}
func (this *Account) UnLock() {
	if !this.Disabled() {
		this.ErrCnt = 0
		this.Status = "Active"
	}
}
func (this *Account) Disable() {
	this.Status = "Disabled"
}
func (this *Account) Enable() {
	this.ErrCnt = 0
	this.Status = "Active"
}

func (this *Account) GetPwd() string {
	pwd, _ := ENC.Decrypt(this.Pwd)
	return pwd
}

func (this *Account) SetPwd(pwd string) {
	epwd, _ := ENC.Encrypt(pwd)
	this.Pwd = string(epwd)
}

func (this *Account) ValidPwd(pwd string) bool {
	return this.GetPwd() == pwd
}

func (this *Account) ValidEPwd(epwd string) bool {
	return (this.Pwd == epwd)
}

//Get status
func (this *Account) Locked() bool {
	return this.Status == "Locked"
}
func (this *Account) Disabled() bool {
	return this.Status == "Disabled"
}
func (this *Account) IsActive() bool {
	return this.Status == "Active"
}

func (this *Account) Register() error {
	if this.Uname == "" {
		return errors.New("Invalid Username")
	}
	if this.Pwd == "" {
		return errors.New("Invalid Password")
	}
	if this.Manager == "" {
		return errors.New("Invalid Manager")
	} else if _, err := GetAccount(this.Manager); err != nil {
		return errors.New("Invalid Manager")
	}
	this.ErrCnt = 0
	this.Disable() //New account need manager enable
	return AddAccount(this)

}

func (this *Account) GetReporters() ([]*Account, error) {
	if this.IsActive() {
		if this.IsAdmin() || this.IsManager() {
			filter := map[string]string{
				"manager": this.Uname,
			}
			return GetAcctounts(filter)
		}
	}
	return nil, errors.New("No reporter")
}
