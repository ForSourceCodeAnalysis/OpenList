package extensions

import (
	"github.com/OpenListTeam/OpenList/v4/extensions/handles"
	"github.com/OpenListTeam/OpenList/v4/extensions/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterRoute register extension routes
func RegisterRoute(g *gin.RouterGroup) {
	ext := g.Group("/ext")
	fsRoute(ext)
}

func fsRoute(g *gin.RouterGroup) {
	fs := g.Group("/fs")
	fs.POST("/batch_rename", middlewares.BatchRename, handles.BatchRename)
	fs.POST("/batch_remove", middlewares.BatchRemove, handles.BatchRemove)
}
