package http

import (
	"github.com/LittleLe6owski/LinkShortener/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LinkHandler struct {
	linkUseCase domain.LinkUseCase
}

func NewLinkHandler(e *echo.Echo, linkUseCase domain.LinkUseCase) {
	handler := &LinkHandler{
		linkUseCase: linkUseCase,
	}
	l := e.Group("/link")

	l.POST("/create/", handler.addLink)
	l.GET("/", handler.getLink)
}

func (l *LinkHandler) addLink(e echo.Context) error {
	qp := e.QueryParams()
	if !qp.Has("url") {
		return e.NoContent(http.StatusBadRequest)
	}
	link, err := l.linkUseCase.GenShortLink(qp.Get("url"))
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}
	return e.JSON(http.StatusCreated, link)
}

func (l *LinkHandler) getLink(e echo.Context) error {
	qp := e.QueryParams()
	if !qp.Has("url") {
		return e.NoContent(http.StatusBadRequest)
	}
	link, err := l.linkUseCase.GetOriginalLink(qp.Get("url"))
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}
	return e.Redirect(http.StatusFound, link)
}
