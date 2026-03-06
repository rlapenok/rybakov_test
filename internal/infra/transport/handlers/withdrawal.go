package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rlapenok/rybakov_test/internal/domain"
	"github.com/rlapenok/rybakov_test/internal/usecase"
)

type WithdrawalHandler struct {
	withdrawalUseCase *usecase.WithdrawalUseCase
}

func NewWithdrawalHandler(withdrawalUseCase *usecase.WithdrawalUseCase) *WithdrawalHandler {
	return &WithdrawalHandler{withdrawalUseCase: withdrawalUseCase}
}

func (h *WithdrawalHandler) CreateWithdrawal(c *gin.Context) {
	var req usecase.CreateWithdrawalInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(domain.ErrInvalidRequestPayload.AddMeta("original_error", err.Error()))
		return
	}

	out, err := h.withdrawalUseCase.CreateWithdrawal(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (h *WithdrawalHandler) GetWithdrawalByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(domain.ErrWithdrawalNotFound)
		return
	}

	out, err := h.withdrawalUseCase.GetWithdrawalByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}
