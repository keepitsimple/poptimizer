package port

import (
	"context"
	"fmt"

	"github.com/WLM1ke/poptimizer/opt/internal/domain"
	"github.com/WLM1ke/poptimizer/opt/internal/domain/portfolio"
)

// AccountsService сервис для редактирования счетов.
type AccountsService struct {
	repo domain.FullRepo[Portfolio]
}

// NewAccountsService создает сервис редактирования брокерских счетов.
func NewAccountsService(repo domain.FullRepo[Portfolio]) *AccountsService {
	return &AccountsService{repo: repo}
}

// AccountsDTO содержит перечень доступных счетов.
type AccountsDTO []string

// GetAccountNames возвращает перечень существующих счетов.
func (s AccountsService) GetAccountNames(ctx context.Context) (AccountsDTO, domain.ServiceError) {
	qids, err := s.repo.List(ctx, portfolio.Subdomain, _AccountsGroup)
	if err != nil {
		return nil, domain.NewServiceInternalErr(err)
	}

	if len(qids) <= 1 {
		return nil, nil
	}

	acc := make(AccountsDTO, 0, len(qids)-1)

	for _, qid := range qids {
		if qid != _NewAccount {
			acc = append(acc, qid)
		}
	}

	return acc, nil
}

// CreateAccount создает новый счет с выбранными бумагами.
func (s AccountsService) CreateAccount(ctx context.Context, name string) domain.ServiceError {
	if name == _NewAccount {
		return domain.NewBadServiceRequestErr("reserved name %s", name)
	}

	tmplAgg, err := s.repo.Get(ctx, AccountID(_NewAccount))
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	agg, err := s.repo.Get(ctx, AccountID(name))
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	if !agg.Timestamp.IsZero() {
		return domain.NewBadServiceRequestErr("can't create existing account %s", name)
	}

	agg.Timestamp = tmplAgg.Timestamp
	agg.Entity = tmplAgg.Entity

	if err := s.repo.Save(ctx, agg); err != nil {
		return domain.NewServiceInternalErr(err)
	}

	return nil
}

// DeleteAccount удаляет счет.
func (s AccountsService) DeleteAccount(ctx context.Context, name string) domain.ServiceError {
	err := s.repo.Delete(ctx, AccountID(name))
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	return nil
}

// PositionDTO информация об отдельной позиции.
type PositionDTO struct {
	Ticker   string  `json:"ticker"`
	Shares   int     `json:"shares"`
	Lot      int     `json:"lot"`
	Price    float64 `json:"price"`
	Turnover float64 `json:"turnover"`
}

// AccountDTO информация об отдельном счете.
type AccountDTO struct {
	Positions []PositionDTO `json:"positions"`
	Cash      int           `json:"cash"`
}

// GetAccount выдает информацию о счет по указанному имени.
func (s AccountsService) GetAccount(ctx context.Context, name string) (AccountDTO, domain.ServiceError) {
	var dto AccountDTO

	agg, err := s.repo.Get(ctx, AccountID(name))
	if err != nil {
		return dto, domain.NewServiceInternalErr(err)
	}

	if len(agg.Entity.Positions) == 0 {
		return dto, domain.NewServiceInternalErr(fmt.Errorf("wrong account name %s", name))
	}

	dto.Cash = agg.Entity.Cash

	for _, pos := range agg.Entity.Positions {
		dto.Positions = append(dto.Positions, PositionDTO{
			Ticker:   pos.Ticker,
			Shares:   pos.Shares,
			Lot:      pos.Lot,
			Price:    pos.Price,
			Turnover: pos.Turnover,
		})
	}

	return dto, nil
}

// UpdateDTO содержит информацию об обновленных позициях.
type UpdateDTO []struct {
	Ticker string `json:"ticker"`
	Shares int    `json:"shares"`
}

// UpdateAccount меняет значение количества акций для заданного счета и тикера.
//
// Для изменения количества денег необходимо указать тикер CASH.
func (s AccountsService) UpdateAccount(ctx context.Context, name string, dto UpdateDTO) domain.ServiceError {
	aggs, err := s.repo.GetGroup(ctx, portfolio.Subdomain, _AccountsGroup)
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	portQID := PortfolioDateID(aggs[0].Timestamp)

	port, err := s.repo.Get(ctx, portQID)
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	port.Timestamp = aggs[0].Timestamp
	port.Entity = aggs[0].Entity

	for count, agg := range aggs {
		if agg.QID().ID == name {
			for _, pos := range dto {
				if err := agg.Entity.SetAmount(pos.Ticker, pos.Shares); err != nil {
					return domain.NewBadServiceRequestErr("%s", err)
				}
			}

			if err := s.repo.Save(ctx, agg); err != nil {
				return domain.NewServiceInternalErr(err)
			}
		}

		if count == 0 {
			continue
		}

		port.Entity = port.Entity.Sum(agg.Entity)
	}

	err = s.repo.Save(ctx, port)
	if err != nil {
		return domain.NewServiceInternalErr(err)
	}

	return nil
}

// PortfolioService сервис для редактирования счетов.
type PortfolioService struct {
	repo domain.GetListRepo[Portfolio]
}

// NewPortfolioService создает сервис для просмотра информации о портфеле.
func NewPortfolioService(repo domain.GetListRepo[Portfolio]) *PortfolioService {
	return &PortfolioService{repo: repo}
}

// GetPortfolioDates выдает перечень дат, для которых есть информация о портфеле.
func (s PortfolioService) GetPortfolioDates(ctx context.Context) (AccountsDTO, domain.ServiceError) {
	qids, err := s.repo.List(ctx, portfolio.Subdomain, _PortfolioGroup)
	if err != nil {
		return nil, domain.NewServiceInternalErr(err)
	}

	return qids, nil
}

// GetPortfolio выдает информацию о портфеле для заданной даты.
func (s PortfolioService) GetPortfolio(ctx context.Context, date string) (AccountDTO, domain.ServiceError) {
	var dto AccountDTO

	agg, err := s.repo.Get(ctx, PortfolioID(date))
	if err != nil {
		return dto, domain.NewServiceInternalErr(err)
	}

	if len(agg.Entity.Positions) == 0 {
		return dto, domain.NewBadServiceRequestErr("wrong portfolio date %s", date)
	}

	dto.Cash = agg.Entity.Cash

	for _, pos := range agg.Entity.Positions {
		dto.Positions = append(dto.Positions, PositionDTO{
			Ticker:   pos.Ticker,
			Shares:   pos.Shares,
			Lot:      pos.Lot,
			Price:    pos.Price,
			Turnover: pos.Turnover,
		})
	}

	return dto, nil
}
