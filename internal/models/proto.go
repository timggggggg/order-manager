package models

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/pkg/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderModelToProto(order *Order) *api.TOrder {
	return &api.TOrder{
		ID:             order.ID,
		UserID:         order.UserID,
		Status:         StatusModelToProto(order.Status),
		AcceptTime:     TimeModelToProto(order.AcceptTime),
		ExpireTime:     TimeModelToProto(order.ExpireTime),
		IssueTime:      TimeModelToProto(order.IssueTime),
		Weight:         order.Weight,
		Cost:           order.Cost.String(),
		Packaging:      PackagingModelToProto(order.Package),
		ExtraPackaging: PackagingModelToProto(order.ExtraPackage),
	}
}

func OrderProtoToModel(order *api.TOrder) *Order {
	cost, err := NewMoney(order.Cost)
	if err != nil {
		fmt.Printf("error pasrsing money string: %v", err)
	}

	return &Order{
		ID:           order.ID,
		UserID:       order.UserID,
		Status:       StatusProtoToModel(order.Status),
		AcceptTime:   TimeProtoToModel(order.AcceptTime),
		ExpireTime:   TimeProtoToModel(order.ExpireTime),
		IssueTime:    TimeProtoToModel(order.IssueTime),
		Weight:       order.Weight,
		Cost:         cost,
		Package:      PackagingProtoToModel(order.Packaging),
		ExtraPackage: PackagingProtoToModel(order.ExtraPackaging),
	}
}

func OrdersSliceModelToProto(orders OrdersSliceStorage) []*api.TOrder {
	res := make([]*api.TOrder, 0, len(orders))
	for _, order := range orders {
		res = append(res, OrderModelToProto(order))
	}
	return res
}

func OrdersSliceProtoToModel(orders []*api.TOrder) OrdersSliceStorage {
	res := make(OrdersSliceStorage, 0, len(orders))
	for _, order := range orders {
		res = append(res, OrderProtoToModel(order))
	}
	return res
}

func TimeModelToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func TimeProtoToModel(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

func PackagingModelToProto(p PackagingType) api.TPackagingType {
	return api.TPackagingType(api.TPackagingType_value[string(p)])
}

func PackagingProtoToModel(p api.TPackagingType) PackagingType {
	return PackagingType(api.TPackagingType_name[int32(p)])
}

func StatusModelToProto(p OrderStatus) api.TOrderStatus {
	return api.TOrderStatus(api.TOrderStatus_value[string(p)])
}

func StatusProtoToModel(p api.TOrderStatus) OrderStatus {
	return OrderStatus(api.TOrderStatus_name[int32(p)])
}
