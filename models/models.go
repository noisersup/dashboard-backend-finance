package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Title 			string
	MaxExpenses		int64
	CurrExpenses	int64
	Expenses		[]Expense
}

type Expense struct {
	gorm.Model
	Title	string
	Cost	Currency
}


type Currency struct {
	First 	int64
	Second	int8
}