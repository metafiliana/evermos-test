package order

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type cacheRepository struct {
	redisPool *redis.Pool
}

func NewCacheRepository(redisPool *redis.Pool) CacheRepository {
	// set initial stock values

	return &cacheRepository{redisPool}
}

const (
	CacheKeyTemplateStock        = "S_%d"  // itemID
	CacheKeyTemplateReservedItem = "RI_%d" // userID
)

type CacheRepository interface {
	SetReservedItems(userID int, reservedItems []*ReservedItem) error
	GetReservedItemsByUserID(userID int) ([]*ReservedItem, error)
	RemoveReservedItemsByUserID(userID int) error

	GetItemStockMapByIDs(itemIDs []int) (map[int]int, error)
	IncreaseStockByItemID(itemID, qty int) error
	DecreaseStockByItemID(itemID, qty int) error
}

func (cr *cacheRepository) SetReservedItems(userID int, reservedItems []*ReservedItem) error {
	// persist reserved item data uniquely by user
	key := fmt.Sprintf(CacheKeyTemplateReservedItem, userID)

	dataBytes, err := json.Marshal(reservedItems)
	if err != nil {
		return err
	}

	conn := cr.redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, dataBytes)
	if err != nil {

		return err
	}

	return nil
}

func (cr *cacheRepository) GetReservedItemsByUserID(userID int) ([]*ReservedItem, error) {
	// persist reserved item data uniquely by user
	key := fmt.Sprintf(CacheKeyTemplateReservedItem, userID)

	conn := cr.redisPool.Get()
	defer conn.Close()

	dataInBytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}

		return nil, err
	}

	if len(dataInBytes) == 0 {
		return nil, nil
	}

	var reservedItems []*ReservedItem
	if err := json.Unmarshal(dataInBytes, &reservedItems); err != nil {
		return nil, err
	}

	return reservedItems, nil
}

func (cr *cacheRepository) RemoveReservedItemsByUserID(userID int) error {
	// persist reserved item data uniquely by user
	key := fmt.Sprintf(CacheKeyTemplateReservedItem, userID)

	conn := cr.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)

	if err != nil {
		return err
	}

	return nil
}

func (cr *cacheRepository) GetItemStockMapByIDs(itemIDs []int) (map[int]int, error) {
	conn := cr.redisPool.Get()
	defer conn.Close()

	mapItemIDsStock := make(map[int]int)
	for _, itemID := range itemIDs {
		key := fmt.Sprintf(CacheKeyTemplateStock, itemID)
		val, err := redis.Int(conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				return nil, nil
			}

			return nil, err
		}

		mapItemIDsStock[itemID] = val
	}

	return mapItemIDsStock, nil
}

func (cr *cacheRepository) DecreaseStockByItemID(itemID, qty int) error {
	key := fmt.Sprintf(CacheKeyTemplateStock, itemID)
	conn := cr.redisPool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("DECRBY", key, qty))
	if err != nil {
		return err
	}

	return nil
}

func (cr *cacheRepository) IncreaseStockByItemID(itemID, qty int) error {
	key := fmt.Sprintf(CacheKeyTemplateStock, itemID)
	conn := cr.redisPool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("INCRBY", key, qty))
	if err != nil {
		return err
	}

	return nil
}
