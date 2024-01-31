package internal

import (
	"go-micro.dev/v4/broker"
	"reflect"
	"testing"
)

func TestEngineService_Init(t *testing.T) {
	type fields struct {
		mq   *mq.Service
		buy  *BuyBook
		sell *SellBook
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EngineService{
				mq:   tt.fields.mq,
				buy:  tt.fields.buy,
				sell: tt.fields.sell,
			}
			e.Init()
		})
	}
}

func TestEngineService_createOrder(t *testing.T) {
	type fields struct {
		mq   *mq.Service
		buy  *BuyBook
		sell *SellBook
	}
	type args struct {
		order *TrustOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EngineService{
				mq:   tt.fields.mq,
				buy:  tt.fields.buy,
				sell: tt.fields.sell,
			}
			e.createOrder(tt.args.order)
		})
	}
}

func TestEngineService_processMsg(t *testing.T) {
	type fields struct {
		mq   *mq.Service
		buy  *BuyBook
		sell *SellBook
	}
	type args struct {
		event broker.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EngineService{
				mq:   tt.fields.mq,
				buy:  tt.fields.buy,
				sell: tt.fields.sell,
			}
			if err := e.processMsg(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("processMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngineService_processOrder(t *testing.T) {
	type fields struct {
		mq   *mq.Service
		buy  *BuyBook
		sell *SellBook
	}
	type args struct {
		takerOrder   *TrustOrder
		makerBooks   Books
		anotherBooks Books
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EngineService{
				mq:   tt.fields.mq,
				buy:  tt.fields.buy,
				sell: tt.fields.sell,
			}
			e.processOrder(tt.args.takerOrder, tt.args.makerBooks, tt.args.anotherBooks)
		})
	}
}

func Test_minFloat(t *testing.T) {
	type args struct {
		x *big.Float
		y *big.Float
	}
	tests := []struct {
		name string
		args args
		want *big.Float
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minFloat(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("minFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
