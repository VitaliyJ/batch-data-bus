package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"data-bus/internal/entity"
	"data-bus/internal/usecase/mocks"
	"data-bus/pkg/externalservice"
	"github.com/stretchr/testify/mock"
)

func TestBatchUseCase_Process(t *testing.T) {
	tests := []struct {
		name          string
		items         []externalservice.Item
		itemsLimit    uint64
		batchesNumber int
		wantErr       bool
		serviceCalled bool
	}{
		{
			name:          "nil items",
			items:         nil,
			itemsLimit:    10,
			serviceCalled: false,
		},
		{
			name:          "empty items",
			items:         make([]externalservice.Item, 0),
			itemsLimit:    10,
			batchesNumber: 0,
			serviceCalled: false,
		},
		{
			name:          "success 10 items",
			items:         make([]externalservice.Item, 10),
			itemsLimit:    4,
			batchesNumber: 3,
			serviceCalled: true,
		},
		{
			name:          "success 1000 items",
			items:         make([]externalservice.Item, 1000),
			itemsLimit:    10,
			batchesNumber: 100,
			serviceCalled: true,
		},
		{
			name:          "success 1000 items with limit 1",
			items:         make([]externalservice.Item, 1000),
			itemsLimit:    1,
			batchesNumber: 1000,
			serviceCalled: true,
		},
		{
			name:          "with error from service",
			items:         make([]externalservice.Item, 10),
			itemsLimit:    10,
			batchesNumber: 1,
			wantErr:       true,
			serviceCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// init mocks
			repo := mocks.NewBatchRepo(t)
			service := mocks.NewService(t)

			// set mocks calls assertions
			if len(tt.items) > 0 {
				if tt.wantErr {
					repo.On("AddBatch", ctx, mock.Anything).Return(errors.New("AddBatch error"))
				} else {
					repo.On("AddBatch", ctx, mock.Anything).Return(nil)
				}
			}
			if tt.serviceCalled {
				service.On("GetLimits").Return(tt.itemsLimit, time.Second)
			}

			// run Process()
			uc := NewBatchUseCase(repo, service)
			err := uc.Process(context.Background(), tt.items)
			// check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error: %v; wantErr: %v", err, tt.wantErr)
			}

			// check assertions
			if len(tt.items) == 0 {
				repo.AssertNotCalled(t, "AddBatch", ctx, make([]externalservice.Item, tt.itemsLimit))
			} else {
				repo.AssertNumberOfCalls(t, "AddBatch", tt.batchesNumber)
			}
		})
	}
}

func TestBatchUseCase_SendNextBatch(t *testing.T) {
	tests := []struct {
		name                string
		limit               time.Duration
		batch               entity.Batch
		lastRequestAfterNow bool
		serviceProcessCall  bool
		repoAddBatchCall    bool
		saveLastRecordCall  bool
		wantErr             bool
	}{
		{
			name:               "success process zero last request time",
			limit:              1 * time.Second,
			serviceProcessCall: true,
			repoAddBatchCall:   false,
			saveLastRecordCall: true,
			wantErr:            false,
		},
		{
			name:               "failed process zero last request time",
			limit:              1 * time.Second,
			serviceProcessCall: true,
			repoAddBatchCall:   true,
			saveLastRecordCall: false,
			wantErr:            true,
		},
		{
			name:                "last request time out of limits",
			limit:               100 * time.Second,
			lastRequestAfterNow: true,
			serviceProcessCall:  false,
			repoAddBatchCall:    false,
			saveLastRecordCall:  false,
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// init mocks
			repo := mocks.NewBatchRepo(t)
			service := mocks.NewService(t)

			service.On("GetLimits").Return(uint64(10), tt.limit)
			if tt.lastRequestAfterNow {
				repo.On("GetLastRequestTime", ctx).Return(time.Now().Add(time.Minute), nil)
			} else {
				repo.On("GetLastRequestTime", ctx).Return(time.Now().Add(-time.Minute), nil)
			}

			// set mocks calls assertions
			if tt.repoAddBatchCall {
				repo.On("AddBatch", ctx, mock.Anything).Return(nil)
			}
			if tt.saveLastRecordCall {
				repo.On("SaveLastRequestTime", ctx, mock.Anything).Return(nil)
			}
			if tt.serviceProcessCall {
				repo.On("GetFirstBatchItems", ctx).Return(make(entity.Batch, 1), nil)

				if tt.repoAddBatchCall {
					service.On("Process", ctx, mock.Anything).Return(errors.New("process error"))
				} else {
					service.On("Process", ctx, mock.Anything).Return(nil)
				}
			}

			// run Process()
			uc := NewBatchUseCase(repo, service)
			err := uc.SendNextBatch(context.Background())
			// check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error: %v; wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
