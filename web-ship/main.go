package main

import (
	"github.com/xgfone/ship/v5"
	"github.com/xgfone/ship/v5/middleware"
)

func main() {
	router := ship.New()
	router.Use(middleware.Logger(), middleware.Recover()) // Use the middlewares.

	router.Route("/ping").GET(func(c *ship.Context) error {
		return c.JSON(200, map[string]interface{}{"message": "pong"})
	})

	router.Route("/test").GET(func(c *ship.Context) error {
		return c.JSON(200, map[string]interface{}{"message": c.Query("name")})
	})

	group := router.Group("/group")
	group.Route("/ping").GET(func(c *ship.Context) error {
		return c.Text(200, "group")
	})

	subgroup := group.Group("/subgroup")
	subgroup.Route("/ping").GET(func(c *ship.Context) error {
		return c.Text(200, "subgroup")
	})

	// Start the HTTP server.
	ship.StartServer(":8080", router)
	// or
	// http.ListenAndServe(":8080", router)
}
