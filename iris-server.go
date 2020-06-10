package main

import (
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func checkReqEnvVars() {
	checkEnvVars := func(key string) {
		val, ok := os.LookupEnv(key)
		if !ok {
			fmt.Printf("%s is not set\n", key)
			// log.Fatal("Exception: value is not set")
		} else {
			fmt.Printf("%s=%s\n", key, val)
		}
	}
	checkEnvVars("RODB_DEBUG")
	checkEnvVars("RODB_DB_HOST")
	checkEnvVars("RODB_DB_USERNAME")
	checkEnvVars("RODB_DB_PASSWORD")
}

type mobStruct struct {
	Id       uint16 `json:"id"`
	Level    uint16 `json:"level"`
	Hp       uint16 `json:"hp"`
	Def      uint16 `json:"def"`
	Mdef     uint16 `json:"mdef"`
	Hit      uint16 `json:"hit"`
	Flee     uint16 `json:"flee"`
	Atk      uint16 `json:"atk"`
	Matk     uint16 `json:"matk"`
	Name     string `json:"name"`
	Race     string `json:"race"`
	Element  string `json:"element"`
	Size     string `json:"size"`
	Location string `json:"location"`
}

//func getMonsterInfo(id uint64) {
//	ctx.Writef("%+v\n", simpleMob[id])
//}

func main() {
	simpleMob := make(map[uint16]mobStruct)
	simpleMob[1] = mobStruct{1, 15, 570, 16, 3, 151, 9, 16, 8, "Rocker", "Insect", "Earth", "Medium", "Prontera Field"}
	simpleMob[2] = mobStruct{2, 14, 338, 22, 0, 142, 9, 13, 7, "Savage Babe", "Brute", "Earth", "Small", "Geffen East Field, Prontera Field"}
	// Check that all required variables are set in environments variables
	checkReqEnvVars()
	app := iris.New()
	app.OnErrorCode(iris.StatusNotFound, notFound)
	// Register custom handler for specific http errors.
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		// .Values are used to communicate between handlers, middleware.
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})
	app.Get("/", index)
	// Parse all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	tmpl := iris.HTML("./views", ".html")
	// Enable re-build on local template files changes.
	tmpl.Reload(true)
	// Register the view engine to the views,
	// this will load the templates.
	app.RegisterView(tmpl)
	// Default template funcs are:
	//
	// - {{ urlpath "myNamedRoute" "pathParameter_ifNeeded" }}
	// - {{ render "header.html" }}
	// and partial relative path to current page:
	// - {{ render_r "header.html" }}
	// - {{ yield }}
	// - {{ current }}
	// Register a custom template func:
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})

	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:    GET
	// Resource:  http://localhost:8080/monster/42
	// app.Get("/user/{id:string regexp(^[0-9]+$)}")
	app.Get("/monster/{uid:uint16 max(65534)}", func(ctx iris.Context) {
		mobID := ctx.Params().GetUint16Default("uid", 0)
		// ctx.JSON(getMonsterInfo)
		ctx.JSON(simpleMob[mobID])
	})

	// app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
	//	ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
	//})
	// app.Get("/u/{username:string}", func(ctx iris.Context) {
	//	ctx.Writef("username (string): %s", ctx.Params().Get("username"))
	//})

	// http://localhost:8080

	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutBodyConsumptionOnUnmarshal,
		iris.WithoutAutoFireStatusCode,
		iris.WithOptimizations,
		iris.WithTimeFormat("Mon, 01 Jan 2006 15:04:05 GMT"),
		iris.WithOtherValue("ServerName", "my amazing server"),
		iris.WithOtherValue("ServerOwner", "admin@example.com"),
		iris.WithCharset("UTF-8"),
	)
}

//[{"id": 2, "class": "archangel", "name": "Michael", "superpower": "resurrection", "description": "With their great offensive and defensive power and flying ability, angels and archangels are among the mightiest creatures in the game.", "image": "https://img.icons8.com/cute-clipart/64/000000/angel.png"}

func notFound(ctx iris.Context) {
	// when 404 then render the template
	// $views_dir/errors/404.html
	ctx.View("errors/404.html")

	// OR, if you had more paths and you want
	// to suggest relative paths to the user:
	// suggestPaths := ctx.FindClosest(3)
	// if len(suggestPaths) == 0 {
	//     ctx.WriteString("404 not found")
	//     return
	// }
	//
	// ctx.HTML("Did you mean?<ul>")
	// for _, s := range suggestPaths {
	//     ctx.HTML(`<li><a href="%s">%s</a></li>`, s, s)
	// }
	// ctx.HTML("</ul>")
}

func index(ctx iris.Context) {
	ctx.View("index.html")
}
