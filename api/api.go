package api

import (
	"context"

	"github.com/noisersup/dashboard-backend-finance/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Api struct{

}

func (api *Api) GetGroups(ctx context.Context, empty *emptypb.Empty) (*pb.Groups,error) {
	expense := pb.Expense{
		Id: 123,
		Title: "Expense",
		Cost: 1011.1,
	}
	group := pb.Group{
		Id: 1,
		Title: "Group",
		MaxExpenses: 101.01,
		CurrExpenses: 101.01,
		Expenses: []*pb.Expense{
			&expense,
		},
	}
	groups := pb.Groups{
		Groups: []*pb.Group{
			&group,
		},
	}

	return &groups,nil
}