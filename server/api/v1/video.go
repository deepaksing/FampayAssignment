package v1

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func (s *ApiV1Service) RegisterVideoRoute(g *echo.Group) {
	g.GET("/videos", s.GetVideos)
	g.GET("/search", s.SearchVideos)
}

func (s *ApiV1Service) GetVideos(c echo.Context) error {
	// Default to page 1 if invalid or not provided
	pageNumStr := c.QueryParam("pagenum")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	// Default page size if invalid or not provided
	pageSizeStr := c.QueryParam("pagesize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	videos, err := s.store.GetVideosFromDB(c.Request().Context(), pageNum, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, videos)
}

func (s *ApiV1Service) SearchVideos(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("query"))
	query = strings.ToLower(query)

	matchedVideos, err := s.store.SearchInVideos(c.Request().Context(), query)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, matchedVideos)
}
