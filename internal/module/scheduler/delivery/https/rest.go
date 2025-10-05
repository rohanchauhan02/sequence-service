package https

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/sequence-service/internal/module/scheduler"
)

type schedulerHandler struct {
	usecase scheduler.Usecase
}

func NewSchedulerHandler(e *echo.Echo, usecase scheduler.Usecase) {
	h := schedulerHandler{
		usecase: usecase,
	}
	_ = h
}
