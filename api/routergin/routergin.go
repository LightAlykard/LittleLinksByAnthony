package routergin

import (
	"fmt"
	"net/http"

	"github.com/LightAlykard/LittleLinksByAnthony/api/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RouterGin struct {
	*gin.Engine
	hs *handler.Handlers
}

func GinAuthMW(c *gin.Context) {
	if u, p, ok := c.Request.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unautorized"))
		return
	}
	c.Next()
}

func NewRouterGin(hs *handler.Handlers) *RouterGin {
	r := gin.Default()
	rtg := &RouterGin{
		hs: hs,
	}

	r.Use(GinAuthMW)

	r.POST("/create", rtg.CreateLink)
	r.GET("/read/:id", rtg.ReadLink)
	r.DELETE("/delete/:id", rtg.DeleteLink)

	rtg.Engine = r
	return rtg
}

type User handler.Link

func (rt *RouterGin) CreateLink(c *gin.Context) {
	sid := c.Param("id")

	lid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := rt.hs.ReadLink(c.Request.Context(), lid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (rt *RouterGin) ReadLink(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := rt.hs.ReadLink(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (rt *RouterGin) DeleteLink(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := rt.hs.DeleteLink(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}
