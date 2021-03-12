package handler

import (
	"context"

	"github.com/baoshuai123/go-micro-order/domain/model"

	common "github.com/baoshuai123/go-micro-common"

	. "github.com/baoshuai123/go-micro-order/proto/order"

	"github.com/baoshuai123/go-micro-order/domain/service"
)

type Order struct {
	OrderDataService service.IOrderDataService
}

//根据订单id查询订单
func (o *Order) GetOrderByID(ctx context.Context, req *OrderID, rsp *OrderInfo) error {
	order, err := o.OrderDataService.FindOrderByID(req.OrderId)
	if err != nil {
		return err
	}
	if err := common.SwapTo(order, rsp); err != nil {
		return err
	}
	return nil
}

//查找所有订单
func (o *Order) GetAllOrder(ctx context.Context, req *AllOrderRes, rsp *AllOrder) error {
	orderAll, err := o.OrderDataService.FindAllOrder()
	if err != nil {
		return err
	}
	for _, v := range orderAll {
		order := &OrderInfo{}
		if err := common.SwapTo(v, order); err != nil {
			return err
		}
		rsp.OrderInfo = append(rsp.OrderInfo, order)
	}
	return nil
}

//生成订单
func (o *Order) CreateOrder(ctx context.Context, req *OrderInfo, rsp *OrderID) error {
	orderADD := &model.Order{}
	if err := common.SwapTo(req, orderADD); err != nil {
		return err
	}
	orderID, err := o.OrderDataService.AddOrder(orderADD)
	if err != nil {
		return err
	}
	rsp.OrderId = orderID
	return nil
}

//删除订单
func (o *Order) DeleteOrderByID(ctx context.Context, req *OrderID, rsp *Response) error {
	err := o.OrderDataService.DeleteOrder(req.OrderId)
	if err != nil {
		return err
	}
	rsp.Msg = "success"
	return nil
}

//更新支付状态
func (o *Order) UpdateOrderPayStatus(ctx context.Context, req *PayStatus, rsp *Response) error {
	if err := o.OrderDataService.UpdatePayStatus(req.OrderId, req.PayStatus); err != nil {
		return err
	}
	rsp.Msg = "订单支付成功"
	return nil
}

//更新发货状态
func (o *Order) UpdateOrderShipStatus(ctx context.Context, req *ShipStatus, rsp *Response) error {
	if err := o.OrderDataService.UpdateShipStatus(req.OrderId, req.ShipStatus); err != nil {
		return err
	}
	rsp.Msg = "订单发货成功"
	return nil
}

//更新订单状态
func (o *Order) UpdateOrder(ctx context.Context, req *OrderInfo, rsp *Response) error {
	order := &model.Order{}
	if err := common.SwapTo(req, order); err != nil {
		return err
	}
	if err := o.OrderDataService.UpdateOrder(order); err != nil {
		return err
	}
	rsp.Msg = "订单更新成功"
	return nil
}
