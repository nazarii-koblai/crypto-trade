package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/crypto-trade/model"
	"github.com/gofrs/uuid"
)

const (
	addUser    = `INSERT INTO users (email, user_name, surname, password) VALUES($1, $2, $3, $4)`
	addWallet  = `INSERT INTO wallets (id, currency, balance) VALUES($1, $2, $3)`
	addWallets = `INSERT INTO users_wallets (email, wallet_id) VALUES($1, $2)`

	passwordByEmail = `SELECT password FROM users WHERE email = $1`
)

const (
	defaultBalance = 100
)

// HashedPassword returns user hashed password by email.
func (s *Storage) HashedPassword(ctx context.Context, email model.Email) (model.Password, error) {
	var password model.Password
	row := s.db.QueryRowContext(ctx, passwordByEmail, email)
	if err := row.Scan(&password); err != nil {
		return "", errorFromCode(err)
	}
	return password, nil
}

// AddNewUser adds new user to storage.
func (s *Storage) AddNewUser(ctx context.Context, user model.User) (model.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, fmt.Errorf("can't create new user transaction: %w", errorFromCode(err))
	}

	if _, err = tx.Exec(addUser, user.Email, user.Name, user.Surname, user.Password); err != nil {
		tx.Rollback()
		return model.User{}, fmt.Errorf("can't create user transaction: %w", errorFromCode(err))
	}

	defaultWallets := []model.CurrencyType{model.FBTC, model.FETH}

	valueStrings := make([]string, 0, len(defaultWallets))
	valueArgs := make([]interface{}, 0, len(defaultWallets)*3)

	for i, currency := range defaultWallets {
		uuid, err := uuid.NewV4()
		if err != nil {
			return model.User{}, fmt.Errorf("can't generate uuid: %w", errorFromCode(err))
		}

		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d)",
			i*3+1,
			i*3+2,
			i*3+3,
		))

		valueArgs = append(valueArgs, uuid, currency.String(), defaultBalance)

		user.Wallets = append(user.Wallets, model.Wallet{
			ID:       uuid,
			Currency: currency.String(),
			Balance:  defaultBalance,
		})
	}

	createDefaultWallets := fmt.Sprintf("INSERT INTO wallets (id, currency, balance) VALUES %s", strings.Join(valueStrings, ","))
	if _, err = tx.Exec(createDefaultWallets, valueArgs...); err != nil {
		tx.Rollback()
		return model.User{}, fmt.Errorf("can't create default wallets transaction: %w", errorFromCode(err))
	}

	valueStrings = make([]string, 0, len(defaultWallets))
	valueArgs = make([]interface{}, 0, len(defaultWallets)*2)

	for i := range defaultWallets {
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d)",
			i*2+1,
			i*2+2,
		))
		valueArgs = append(valueArgs, user.Email, user.Wallets[i].ID)
	}

	attachWalletsToEmail := fmt.Sprintf("INSERT INTO users_wallets (email, wallet_id) VALUES %s", strings.Join(valueStrings, ","))
	if _, err = tx.Exec(attachWalletsToEmail, valueArgs...); err != nil {
		tx.Rollback()
		return model.User{}, fmt.Errorf("can't attach wallets transaction: %w", errorFromCode(err))
	}

	return user, tx.Commit()
}
