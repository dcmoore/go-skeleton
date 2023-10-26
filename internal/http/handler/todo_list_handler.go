package handler

import (
	"fmt"
	"net/http"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/http/middleware"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"

	fiber "github.com/gofiber/fiber/v2"
)

type TodoListHandler struct {
	parser          parser.Parser
	presenter       json.JsonPresenter
	todoListUsecase usecase.TodoListUsecase
}

func NewTodoListHandler(
	parser parser.Parser,
	presenter json.JsonPresenter,
	todoListUsecase usecase.TodoListUsecase,
) *TodoListHandler {
	return &TodoListHandler{parser, presenter, todoListUsecase}
}

func (w *TodoListHandler) Register(app fiber.Router) {
	app.Get("/todo-lists/:id", middleware.VerifyJWTToken, w.GetByID)
	app.Get("/todo-lists", middleware.VerifyJWTToken, w.GetByUserID)
	app.Post("/todo-lists", middleware.VerifyJWTToken, w.Create)
	app.Put("/todo-lists/:id", middleware.VerifyJWTToken, w.Update)
	app.Delete("/todo-lists/:id", middleware.VerifyJWTToken, w.Delete)
}

// @Summary         Get Todo List by ID
// @Description     Get a Todo List by its ID
// @Tags            Todo List
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the wallet"
// @Success			201 {object} entity.GeneralResponse{data=entity.TodoListResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/todo-lists/{id} [get]
func (w *TodoListHandler) GetByID(c *fiber.Ctx) error {
	walletID, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.todoListUsecase.GetByID(c.Context(), walletID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Retrieve Todo Lists by User ID
// @Description     Retrieve a list of Todo Lists belonging to a user by their User ID
// @Tags            Todo List
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Success			201 {object} entity.GeneralResponse{data=[]entity.TodoListResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/todo-list [get]
func (w *TodoListHandler) GetByUserID(c *fiber.Ctx) error {
	userID, err := w.parser.ParserUserID(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.todoListUsecase.GetByUserID(c.Context(), userID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	fmt.Println(data)
	fmt.Println(userID)

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary			Create a new Todo List
// @Description		Create a new Todo List
// @Tags			Todo List
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param			req body entity.WalletReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.WalletReq} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/todo-list [post]
func (w *TodoListHandler) Create(c *fiber.Ctx) error {
	var req entity.TodoListReq

	err := w.parser.ParserBodyRequestWithUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.todoListUsecase.Create(c.Context(), &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Update an existing Todo List by ID
// @Description     Update an existing Todo List
// @Tags            Todo List
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the todo list"
// @Param			req body entity.WalletReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/todo-list [put]
func (w *TodoListHandler) Update(c *fiber.Ctx) error {
	var walletReq entity.TodoListReq
	err := w.parser.ParserBodyWithIntIDPathParamsAndUserID(c, &walletReq)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.todoListUsecase.UpdateByID(c.Context(), walletReq)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}

// @Summary         Delete Todo List by ID
// @Description     Delete an existing Todo List by its ID
// @Tags			Todo List
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param           id path int true "ID of the todo list"
// @Success			201 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/todo-lists/{id} [delete]
func (w *TodoListHandler) Delete(c *fiber.Ctx) error {
	walletID, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.todoListUsecase.DeleteByID(c.Context(), walletID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}
