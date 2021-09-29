package order

import (
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	GetMasterStockDataByItemIDs(itemIDs []int) ([]*MasterItemStocks, error)
	CreateOrder(mapItemsStock map[int]*MasterItemStocks, order *Orders) error
}

func (r *repository) GetMasterStockDataByItemIDs(itemIDs []int) ([]*MasterItemStocks, error) {
	var data []*MasterItemStocks
	err := r.db.Model(MasterItemStocks{}).Where(`item_id IN (?)`, itemIDs).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *repository) CreateOrder(mapItemsStock map[int]*MasterItemStocks, order *Orders) error {
	tx := r.db.Begin()

	for itemID, masterStockData := range mapItemsStock {
		err := tx.Model(MasterItemStocks{}).Where(`item_id = ?`, itemID).Update(`stock_qty`, masterStockData.StockQty).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err := tx.Model(Orders{}).Create(order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// build order items data
	for itemID, masterStockData := range mapItemsStock {
		newOrderItem := OrderItems{
			UserID:  order.UserID,
			ItemID:  itemID,
			OrderID: order.ID,
			Qty:     masterStockData.DesiredQty,
		}

		err := tx.Model(OrderItems{}).Create(&newOrderItem).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
