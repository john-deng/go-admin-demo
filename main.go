package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/GoAdminGroup/go-admin/adapter/iris"             // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	_ "github.com/GoAdminGroup/themes/adminlte"                   // ui theme

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/kataras/iris/v12"

	"go-admin-demo/models"
	"go-admin-demo/pages"
	"go-admin-demo/tables"
)

func main() {
	startServer()
}

func startServer() {
	app := iris.Default()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	if err := eng.AddConfigFromYAML("./config.yml").
		AddGenerators(tables.Generators).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "//", pages.GetDashBoard)
	eng.HTMLFile("GET", "///hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	models.Init(eng.MysqlConnection())

	app.HandleDir("/uploads", "./uploads", iris.DirOptions{
		IndexName: "/index.html",
		Gzip:      false,
		ShowList:  false,
	})

	go func() {
		_ = app.Run(iris.Addr(":8080"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
}
