package crud

import (
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sudoblockio/icon-go-api/models"
)

// TransactionCrud - type for transaction table model
type TransactionCrud struct {
	db    *gorm.DB
	model *models.Transaction
}

var transactionCrud *TransactionCrud
var transactionCrudOnce sync.Once

// GetTransactionCrud - create and/or return the transactions table model
func GetTransactionCrud() *TransactionCrud {
	transactionCrudOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		transactionCrud = &TransactionCrud{
			db:    dbConn,
			model: &models.Transaction{},
		}
	})

	return transactionCrud
}

// SelectMany - select from transactions table
// Returns: models, error (if present)
func (m *TransactionCrud) SelectMany(
	limit int,
	skip int,
	from string,
	to string,
	_type string,
	blockNumber int,
	startBlockNumber int,
	endBlockNumber int,
	method string,
	sort string,
) (*[]models.TransactionList, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Latest transactions first
	if sort != "" {
		db = db.Order("block_number " + sort + ", transaction_index")
	} else {
		db = db.Order("transaction_index")
	}

	// from
	if from != "" {
		db = db.Where("from_address = ?", from)
	}

	// to
	if to != "" {
		db = db.Where("to_address = ?", to)
	}

	// type
	if _type != "" {
		db = db.Where("type = ?", _type)
	}

	// block number
	if blockNumber != 0 {
		db = db.Where("block_number = ?", blockNumber)
	}

	// start block number
	if startBlockNumber != 0 {
		db = db.Where("block_number >= ?", startBlockNumber)
	}

	// end block number
	if endBlockNumber != 0 {
		db = db.Where("block_number <= ?", endBlockNumber)
	}

	// method
	if method != "" {
		db = db.Where("method = ?", method)
	}

	// Limit is required and defaulted to 1
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	transactions := &[]models.TransactionList{}
	db = db.Find(transactions)

	return transactions, db.Error
}

// SelectManyByAddress - select from transactions table
// Returns: models, error (if present)
func (m *TransactionCrud) SelectManyByAddress(
	limit int,
	skip int,
	address string,
) (*[]models.TransactionList, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Latest transactions first
	db = db.Order("block_number DESC")

	// Address
	// This replaces a common query with select * from __ where from_address = ... or to_address = ... sort by block_number
	//  which was really slow so we do this subquery to speed up requests from the single page view.
	db = db.Where(`hash
	IN (
		SELECT
			transaction_hash
		FROM
			transaction_by_addresses
		WHERE
			address = ?
		ORDER BY block_number desc
		LIMIT ?
		OFFSET ?
	)`, address, limit, skip)

	// Type
	db = db.Where("type = ?", "transaction")

	transactions := &[]models.TransactionList{}
	db = db.Find(transactions)

	return transactions, db.Error
}

// CountManyIcxByAddress - select from transactions table
// Returns: int64, error (if present)
func (m *TransactionCrud) CountManyIcxByAddress(address string) (int64, error) {
	db := m.db

	db = db.Model(&models.Transaction{}).Where("type='transaction'")
	db = db.Model(&models.Transaction{}).Where("to_address = ? or from_address = ?", address, address)
	db = db.Model(&models.Transaction{}).Where("value_decimal != 0")

	var count int64
	db = db.Count(&count)
	return count, db.Error
}

// SelectManyIcxByAddress - select from transactions table
// Returns: models, error (if present)
func (m *TransactionCrud) SelectManyIcxByAddress(
	limit int,
	skip int,
	address string,
) (*[]models.TransactionList, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Latest transactions first
	db = db.Order("block_number DESC")

	// Address
	db = db.Where("from_address = ? OR to_address = ?", address, address)

	// Type
	db = db.Where("type = ?", "transaction")

	// Non-zero ICX amount
	db = db.Where("value_decimal != 0")

	// Limit is required and defaulted to 1
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	transactions := &[]models.TransactionList{}
	db = db.Find(transactions)

	return transactions, db.Error
}

// SelectManyInternal- select many internal transaction table
// Returns: models, error (if present)
func (m *TransactionCrud) SelectManyInternal(
	limit int,
	skip int,
	hash string,
	blockNumber int,
) (*[]models.TransactionInternalList, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Latest transactions first
	db = db.Order("block_number desc")

	// Hash
	if hash != "" {
		db = db.Where("hash = ?", hash)
	}

	// Block Number
	if blockNumber != 0 {
		db = db.Where("block_number = ?", blockNumber)
	}

	// Internal transactions only
	db = db.Where("type = ?", "log")

	// Limit is required and defaulted to 1
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	transactions := &[]models.TransactionInternalList{}
	db = db.Find(transactions)

	return transactions, db.Error
}

// SelectManyInternalByAddress - select from internal transactions table
// Returns: models, error (if present)
func (m *TransactionCrud) SelectManyInternalByAddress(
	limit int,
	skip int,
	address string,
) (*[]models.TransactionInternalList, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Latest transactions first
	db = db.Order("transactions.block_number DESC")

	// Address
	db = db.Where(`(hash, log_index)
	IN (
		SELECT
			transaction_hash, log_index
		FROM
			transaction_internal_by_addresses
		WHERE
			address = ?
		ORDER BY block_number desc
		LIMIT ?
		OFFSET ?
	)`, address, limit, skip)

	// Type
	db = db.Where("type = ?", "log")

	transactions := &[]models.TransactionInternalList{}
	db = db.Find(transactions)

	return transactions, db.Error
}

// SelectOne - select from transactions table
func (m *TransactionCrud) SelectOne(
	hash string,
	logIndex int32, // Used for internal transactions
) (*models.Transaction, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.Transaction{})

	// Hash
	db = db.Where("hash = ?", hash)

	// Log Index
	db = db.Where("log_index = ?", logIndex)

	transaction := &models.Transaction{}
	db = db.First(transaction)

	return transaction, db.Error
}
