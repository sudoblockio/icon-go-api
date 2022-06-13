package crud

import (
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sudoblockio/icon-go-api/models"
)

// AddressCrud - type for address table model
type AddressCrud struct {
	db    *gorm.DB
	model *models.Address
}

var addressCrud *AddressCrud
var addressCrudOnce sync.Once

// GetAddressCrud - create and/or return the addresss table model
func GetAddressCrud() *AddressCrud {
	addressCrudOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		addressCrud = &AddressCrud{
			db:    dbConn,
			model: &models.Address{},
		}
	})

	return addressCrud
}

// SelectOne - select one from addresses table
func (m *AddressCrud) SelectOne(
	_address string,
) (*models.Address, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Address{})

	// Address
	db = db.Where("address = ?", _address)

	address := &models.Address{}
	db = db.First(address)

	return address, db.Error
}

// SelectMany - select many from addreses table
func (m *AddressCrud) SelectMany(
	limit int,
	skip int,
	address string,
) (*[]models.AddressList, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Address{})

	// Order balances
	db = db.Order("balance DESC")

	// Address
	if address != "" {
		db = db.Where("address = ?", address)
	}

	// Limit
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	addresses := &[]models.AddressList{}
	db = db.Find(addresses)

	return addresses, db.Error
}

// SelectManyContracts - select many from addreses table
func (m *AddressCrud) SelectManyContracts(
	limit int,
	skip int,
) (*[]models.ContractList, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Address{})

	// Order balances
	// db = db.Order("transaction_count DESC")

	// Is contract
	db = db.Where("is_contract = ?", true)

	// Order by name with nulls in the back
	db = db.Order("nullif(name, '') nulls last")

	// Limit
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	contracts := &[]models.ContractList{}
	db = db.Find(contracts)

	return contracts, db.Error
}
