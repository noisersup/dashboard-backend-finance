package models

type Group struct {
	Id 				int
	Title 			string
	MaxExpenses		int64
	CurrExpenses	int64
	Expenses		[]Expense
}

type Expense struct {
	Id 		int
	Title	string
	Cost	float64
}