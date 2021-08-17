package api

import (
	"context"
	"log"

	"github.com/noisersup/dashboard-backend-finance/api/pb"
	"github.com/noisersup/dashboard-backend-finance/database"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Api struct {
	db        *database.Database
	restartDb chan<- bool
}

func InitAPI(database *database.Database, restartChan chan bool) *Api {
	return &Api{database, restartChan}
}

func (api *Api) GetGroups(ctx context.Context, empty *emptypb.Empty) (*pb.Groups, error) {
	groups, err := api.db.GetGroups()
	if err != nil {
		log.Print(err)
		api.restartDb <- true
		return nil, err
	}

	var respGroups pb.Groups
	for _, group := range groups {
		respGroup := pb.Group{
			Id:           int64(group.Id),
			Title:        group.Title,
			MaxExpenses:  float32(group.MaxExpenses),
			CurrExpenses: float32(group.CurrExpenses),
		}

		for _, expense := range group.Expenses {
			respGroup.Expenses = append(respGroup.Expenses, &pb.Expense{
				Id:    int64(expense.Id),
				Title: expense.Title,
				Cost:  float32(expense.Cost),
			})
		}

		respGroups.Groups = append(respGroups.Groups, &respGroup)
	}

	return &respGroups, nil
}

func (api *Api) CreateGroup(ctx context.Context, group *pb.PostGroup) (*emptypb.Empty, error) {
	_, err := api.db.CreateGroup(group.Title, float64(group.MaxExpenses))
	return &emptypb.Empty{}, err
}
