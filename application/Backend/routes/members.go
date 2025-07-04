package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/gin-gonic/gin"
)

func RegisterMemberRoutes(router *gin.Engine) {
	memberRoutes := []struct {
		route string
		file  string
	}{
		{"/elias", "src/html/member-pages/aboutme-elias.html"},
		{"/zachary", "src/html/member-pages/aboutme-zachary.html"},
		{"/jiarong", "src/html/member-pages/aboutme-jiarong.html"},
		{"/hemasri", "src/html/member-pages/aboutme-hemasri.html"},
		{"/weiping", "src/html/member-pages/aboutme-weiping.html"},
	}

	for _, m := range memberRoutes {
		route := m.route
		file := m.file
		router.GET(route, func(c *gin.Context) {
			htmlData, err := core.LoadFrontendFile(file)
			if err != nil {
				renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
				return
			}
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, htmlData)
		})
	}
}
