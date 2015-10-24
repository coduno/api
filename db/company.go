package db

import (
	"github.com/coduno/api/model"
	"golang.org/x/net/context"
)

type CompanyService interface {
	Save(c model.Company) (*model.KeyedCompany, error)
	GetByAddress(address string) (*model.KeyedCompany, error)
}

type companyService struct {
	ctx context.Context
}

func NewCompanyService(ctx context.Context) companyService {
	return companyService{ctx}
}

func (cs companyService) Save(c model.Company) (*model.KeyedCompany, error) {
	key, err := c.Put(cs.ctx, nil)
	if err != nil {
		return nil, err
	}
	return c.Key(key), nil
}

func (cs companyService) GetByAddress(address string) (*model.KeyedCompany, error) {
	var companies []model.Company
	keys, err := model.NewQueryForCompany().
		Filter("Address =", address).
		Limit(1).
		GetAll(cs.ctx, &companies)
	if err != nil {
		return nil, err
	}
	if len(companies) != 1 {
		return nil, nil
	}
	return companies[0].Key(keys[0]), nil
}
