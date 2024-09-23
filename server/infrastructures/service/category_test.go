package service

import (
	"context"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"reflect"
	"testing"
)

func TestCategoryService_FindByCategoryID(t *testing.T) {
	type fields struct {
		baseService baseService
	}
	type args struct {
		ctx   context.Context
		input domains.FindByCategoryIDInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domains.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryService{
				baseService: tt.fields.baseService,
			}
			got, err := c.FindByCategoryID(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByCategoryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByCategoryID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_Save(t *testing.T) {
	type fields struct {
		baseService baseService
	}
	type args struct {
		ctx   context.Context
		input domains.SaveCategoryInput
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
			c := CategoryService{
				baseService: tt.fields.baseService,
			}
			if err := c.Save(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCategoryService(t *testing.T) {
	type args struct {
		db     bun.Tx
		logger echo.Logger
	}
	tests := []struct {
		name string
		args args
		want domains.CategoryRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCategoryService(tt.args.db, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCategoryService() = %v, want %v", got, tt.want)
			}
		})
	}
}
