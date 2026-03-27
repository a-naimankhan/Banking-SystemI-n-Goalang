package postgres

import (
	"BankingSystem/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) GetByID(id string) (*domain.Account, error) {
	var acc domain.Account

	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc, "id = ?", id).Error
	return &acc, err
}

func (r *AccountRepo) Create(acc interface{}) error {
	return r.db.Create(acc).Error
}

func (r *AccountRepo) Save(acc domain.Account) error {
	return r.db.Save(&acc).Error
}

func (r *AccountRepo) Update(acc interface{}) error {
	return r.db.Save(acc).Error
}

func (r *AccountRepo) Delete(id string) error {
	return r.db.Delete(&domain.Account{}, "id = ?", id).Error
}

func (r *AccountRepo) GetAll() ([]domain.Account, error) {
	var accounts []domain.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepo) BeginTx() (domain.AccountRepository, error) {
	tx := r.db.Begin()
	return &AccountRepo{db: tx}, tx.Error
}

func (r *AccountRepo) CommitTx() error {
	return r.db.Commit().Error
}

func (r *AccountRepo) RollbackTx() error {
	return r.db.Rollback().Error
}
