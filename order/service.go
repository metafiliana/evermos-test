package order

import (
	"fmt"

	"github.com/metafiliana/evermos-test/util"
)

type service struct {
	repo      Repository
	cacheRepo CacheRepository
}

func NewService(repo Repository, cacheRepo CacheRepository) Service {
	return &service{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

type Service interface {
	CheckoutItems(req *CheckoutOrderRequest) error
	CreateOrder(req *CreateOrderRequest) error
}

func (s *service) CheckoutItems(req *CheckoutOrderRequest) error {
	// get items id with desired qty map
	mapItemIDQty := make(map[int]int)
	var itemIDs []int
	for _, item := range req.Items {
		mapItemIDQty[item.ItemID] = item.Qty
		itemIDs = append(itemIDs, item.ItemID)
	}

	// get items stock from stock cache data
	stockFromCache, err := s.cacheRepo.GetItemStockMapByIDs(itemIDs)
	if err != nil {
		return util.ErrWrap(err, `fail to get item stock cache by ids`, util.RepositoryError)
	}

	if len(stockFromCache) == 0 {
		return util.ErrWrap(fmt.Errorf("not found"), `desired items have no stock data`, util.NotFound)
	}

	// compare desired qty vs master stock
	var reservedItems []*ReservedItem
	for itemID, desiredQty := range mapItemIDQty {
		cacheStock := stockFromCache[itemID]

		if cacheStock >= desiredQty {
			reservedItems = append(reservedItems, &ReservedItem{
				UserID: req.UserID,
				ItemID: itemID,
				Qty:    desiredQty,
			})

			// minus cacheStock, reserve desired item qty
			go s.cacheRepo.DecreaseStockByItemID(itemID, desiredQty)
		}
	}

	if len(reservedItems) == 0 {
		return util.ErrWrap(fmt.Errorf("bad request"), `stock not available cant reserve any items`, util.BadRequest)
	}

	// set reserved items
	go s.cacheRepo.SetReservedItems(req.UserID, reservedItems)

	return nil
}

func (s *service) CreateOrder(req *CreateOrderRequest) error {
	reservedItems, err := s.cacheRepo.GetReservedItemsByUserID(req.UserID)
	if err != nil {
		return util.ErrWrap(err, `fail to get user reserved item`, util.RepositoryError)
	}

	if len(reservedItems) == 0 {
		return util.ErrWrap(fmt.Errorf("not found"), `no reserved item for current user found`, util.NotFound)
	}

	// check status order
	if req.Status == StatusOrderFailed {
		for _, item := range reservedItems {
			err := s.cacheRepo.IncreaseStockByItemID(item.ItemID, item.Qty)
			if err != nil {
				return util.ErrWrap(err, `fail to increase stock`, util.RepositoryError)
			}
		}
	} else {
		var itemIDs []int
		mapDesiredItemQty := make(map[int]int)
		for _, reservedItem := range reservedItems {
			itemIDs = append(itemIDs, reservedItem.ItemID)
			mapDesiredItemQty[reservedItem.ItemID] = reservedItem.Qty
		}

		// get master stock data from db
		masterStocks, err := s.repo.GetMasterStockDataByItemIDs(itemIDs)
		if err != nil {
			return util.ErrWrap(err, `fail to get master stock data from db`, util.RepositoryError)
		}

		if len(masterStocks) == 0 {
			return util.ErrWrap(err, `master stock for current items have no data`, util.NotFound)
		}

		// build new master stock data
		mapItemsStock := make(map[int]*MasterItemStocks)
		for _, masterStock := range masterStocks {
			desiredQty := mapDesiredItemQty[masterStock.ItemID]

			if masterStock.StockQty < desiredQty {
				fmt.Errorf("master stock for itemID:%d is lesser than desiredQty", masterStock.ItemID)
				continue
			}

			mapItemsStock[masterStock.ItemID] = &MasterItemStocks{
				ItemID:     masterStock.ItemID,
				StockQty:   masterStock.StockQty - desiredQty,
				DesiredQty: desiredQty,
			}
		}

		// escape when order item is 0
		if len(mapItemsStock) == 0 {
			return util.ErrWrap(err, `no order item created`, util.BadRequest)
		}

		// create order
		err = s.repo.CreateOrder(mapItemsStock, &Orders{
			UserID: req.UserID,
			Status: string(req.Status),
		})
		if err != nil {
			return util.ErrWrap(err, `fail to create order`, util.RepositoryError)
		}
	}

	go s.cacheRepo.RemoveReservedItemsByUserID(req.UserID)

	return nil
}
