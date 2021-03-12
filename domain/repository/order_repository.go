package repository

import (
	"errors"

	"github.com/baoshuai123/go-micro-order/domain/model"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	InitTable() error
	FindOrderByID(int64) (*model.Order, error)
	CreateOrder(*model.Order) (int64, error)
	DeleteOrderByID(int64) error
	UpdateOrder(*model.Order) error
	FindAll() ([]model.Order, error)
	UpdatePayStatus(orderID int64, payStatus int32) error
	UpdateShipStatus(orderID int64, shipStatus int32) error
}

//创建orderRepository
func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{mysqlDb: db}
}

type OrderRepository struct {
	mysqlDb *gorm.DB
}

//初始化表
func (u *OrderRepository) InitTable() error {
	return u.mysqlDb.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Order{}, &model.OrderDetail{})
}

//根据ID查找Order信息
func (u *OrderRepository) FindOrderByID(orderID int64) (order *model.Order, err error) {
	order = &model.Order{}
	return order, u.mysqlDb.Preload("OrderDetail").First(order, orderID).Error
}

//创建Order信息
func (u *OrderRepository) CreateOrder(order *model.Order) (int64, error) {
	return order.ID, u.mysqlDb.Create(order).Error
}

//根据ID删除Order信息
func (u *OrderRepository) DeleteOrderByID(orderID int64) error {
	tx := u.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	//彻底删除order
	if err := tx.Unscoped().Where("order_id", orderID).Delete(&model.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//彻底删除OrderDetail
	if err := tx.Unscoped().Where("odder_id", orderID).Delete(&model.OrderDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//更新Order信息
func (u *OrderRepository) UpdateOrder(order *model.Order) error {
	return u.mysqlDb.Model(order).Updates(order).Error
}

//获取结果集
func (u *OrderRepository) FindAll() (orderAll []model.Order, err error) {
	return orderAll, u.mysqlDb.Preload("OrderDetail").Find(&orderAll).Error
}

// 更新订单发货状态
func (u *OrderRepository) UpdateShipStatus(orderID int64, shipStatus int32) error {
	db := u.mysqlDb.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("ship_status", shipStatus)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

//更新订单支付状态
func (u *OrderRepository) UpdatePayStatus(orderID int64, payStatus int32) error {
	db := u.mysqlDb.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("ship_status", payStatus)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil

}
