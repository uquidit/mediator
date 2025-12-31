package mediatorscript

import (
	"github.com/labstack/echo/v4"
)

func AddMediatorscriptAPI(g *echo.Group) {
	g.GET("", GetAll)
	g.GET("/:slug", GetAllByType)

	g.POST("/register", RegisterScript)

	g.DELETE("/unregister-all", UnregisterAll)
	g.DELETE("/unregister/:slug/:script", UnregisterScript)
	g.DELETE("/unregister/:slug", UnregisterScript)

	g.POST("/refresh-all", RefreshAllScript)
	g.POST("/refresh/:slug/:script", RefreshScript)
	g.POST("/refresh/:slug", RefreshScript)

	g.POST("/execute/:script", ExecuteScript)
	g.POST("/execute-scripted-condition/:id", ExecuteScriptedCondition)
	g.POST("/execute-scripted-task/:id", ExecuteScriptedTask)
	g.POST("/execute-pre-assignment", ExecutePreAssignment)
	g.POST("/execute-risk-analysis", ExecuteRiskAnalysis)

	g.POST("/test-all", TestAllScripts)
	g.POST("/test/:slug/:script", TestScript)
	g.POST("/test/:slug", TestScript)

}
