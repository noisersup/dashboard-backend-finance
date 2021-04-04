package models

type Group struct {
	Id 				int
	Title 			string
	MaxExpenses		float64
	CurrExpenses	float64
	Expenses		[]Expense
}

type Expense struct {
	Id 		int
	Title	string
	Cost	float64
}