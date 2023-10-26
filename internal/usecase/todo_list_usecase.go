package usecase

import (
	"context"
	"fmt"
	"time"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
)

type TodoList struct {
	validatorUsecase ValidatorUsecase
	todoListRepo     mysql.TodoListRepository
}

func NewTodoListUsecase(
	validatorUsecase ValidatorUsecase,
	todoListRepo mysql.TodoListRepository,
) *TodoList {
	return &TodoList{validatorUsecase, todoListRepo}
}

type TodoListUsecase interface {
	GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error)
	GetByID(ctx context.Context, walletID int64) (*entity.TodoListResponse, error)
	Create(ctx context.Context, todoListReq *entity.TodoListReq) (*entity.TodoListResponse, error)
	UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error
	DeleteByID(ctx context.Context, todoListID int64) error
}

func (t *TodoList) GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error) {
	funcName := "TodoListUsecase.GetByUserID"
	captureFieldError := entity.CaptureFields{
		"user_id": userID,
	}

	result, err := t.todoListRepo.GetByUserID(ctx, userID)
	if err != nil {
		helper.LogError("todoListRepo.GetByUserID", funcName, err, captureFieldError, "")

		return nil, err
	}

	for _, v := range result {
		res = append(res, &entity.TodoListResponse{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			DoingAt:     v.DoingAt,
			CreatedAt:   v.CreatedAt,
		})
	}

	return res, nil
}

func (t *TodoList) GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error) {
	funcName := "TodoListUsecase.GetByID"
	captureFieldError := entity.CaptureFields{
		"user_id": todoListID,
	}

	data, err := t.todoListRepo.GetByID(ctx, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.TodoListResponse{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
	}, nil
}

func (t *TodoList) Create(ctx context.Context, todoListReq *entity.TodoListReq) (*entity.TodoListResponse, error) {
	funcName := "TodoListUsecase.Create"
	captureFieldError := entity.CaptureFields{
		"user_id": todoListReq.UserID,
		"payload": todoListReq,
	}

	if errMsg := ValidateStruct(*todoListReq); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(apperr.INVALID_PAYLOAD_CODE), errMsg)
	}

	todoListPayload := &mentity.TodoList{
		UserID:      todoListReq.UserID,
		Title:       todoListReq.Title,
		Description: todoListReq.Description,
		DoingAt:     todoListReq.DoingAt,
		CreatedAt:   time.Now(),
	}

	err := t.todoListRepo.Create(ctx, nil, todoListPayload, false)
	if err != nil {
		helper.LogError("todoListRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.TodoListResponse{
		ID:          todoListPayload.ID,
		Title:       todoListPayload.Title,
		Description: todoListPayload.Description,
		DoingAt:     todoListPayload.DoingAt,
		CreatedAt:   todoListPayload.CreatedAt,
	}, nil
}

func (t *TodoList) UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error {
	funcName := "TodoListUsecase.UpdateByID"
	todoListID := todoListReq.ID

	captureFieldError := entity.CaptureFields{
		"user_id": todoListReq.UserID,
		"payload": todoListReq,
	}

	if err := mysql.DBTransaction(t.todoListRepo, func(trx mysql.TrxObj) error {
		lockedWallet, err := t.todoListRepo.LockByID(ctx, trx, todoListID)
		if err != nil {
			helper.LogError("todoListRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}

		if err := t.todoListRepo.Update(ctx, trx, lockedWallet, &mentity.TodoList{
			Title:       todoListReq.Title,
			Description: todoListReq.Description,
			DoingAt:     todoListReq.DoingAt,
			UpdatedAt:   time.Now(),
		}); err != nil {
			helper.LogError("todoListRepo.Update", funcName, err, captureFieldError, "")

			return err
		}

		return nil
	}); err != nil {
		helper.LogError("todoListRepo.DBTransaction", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}

func (t *TodoList) DeleteByID(ctx context.Context, todoListID int64) error {
	funcName := "TodoListUsecase.DeleteByID"
	captureFieldError := entity.CaptureFields{
		"todo_list_id": todoListID,
	}

	err := t.todoListRepo.DeleteByID(ctx, nil, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
