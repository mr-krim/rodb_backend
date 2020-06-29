package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// Define monster struct
type monsterStruct struct {
	MobID       uint16 `json:"mobid"`
	Level       uint16 `json:"level"`
	MaxHP       uint16 `json:"maxhp"`
	Atk         uint16 `json:"atk"`
	Def         uint16 `json:"def"`
	Hit         uint16 `json:"hit"`
	Matk        uint16 `json:"matk"`
	Mdef        uint16 `json:"mdef"`
	Flee        uint16 `json:"flee"`
	Name        string `json:"name"`
	Race        string `json:"race"`
	Element     string `json:"element"`
	Size        string `json:"size"`
	Location    string `json:"location"`
	STR         uint16 `json:"str"`
	VIT         uint16 `json:"vit"`
	LUK         uint16 `json:"luk"`
	INT         uint16 `json:"int"`
	DEX         uint16 `json:"dex"`
	AGI         uint16 `json:"agi"`
	ShortMobURL string `json:"shortmoburl`
	FullMobURL  string `json:"fullmoburl`
}

// Define items struct
type itemStruct struct {
	ItemID       uint16 `json:"itemid"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Weight       uint16 `json:"weight"`
	Usage        string `json:"usage"`
	Description  string `json:"description"`
	ShortItemURL string `json:"shortitemurl`
	FullItemURL  string `json:"fullitemurl`
}

// func PingDB(db *sql.DB) {
// 	err := db.Ping()
// 	ErrorCheck(err)
// }

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

var dbUser, dbName, dbPass, dbHost string
var monster = monsterStruct{}

//var items = itemStruct{}

func checkEnvVars(key string, vars *string) {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("%s is not set\n", key)
		log.Fatal("Exception: value is not set")
	} else {
		*vars = val
	}
}

func checkReqEnvVars() {
	checkEnvVars("RODB_DB_HOST", &dbHost)
	checkEnvVars("RODB_DB_USERNAME", &dbUser)
	checkEnvVars("RODB_DB_PASSWORD", &dbPass)
	checkEnvVars("RODB_DB_NAME", &dbName)
}

func dbConn(dbHost, dbUser, dbPass, dbName *string) (db *sql.DB) {
	dbDriver := "mysql"
	db, connerr := sql.Open(dbDriver, *dbUser+":"+*dbPass+"@tcp("+*dbHost+":3306)/"+*dbName)
	ErrorCheck(connerr)
	return db
}

func showAllItems() (res []itemStruct) {
	db := dbConn(&dbHost, &dbUser, &dbPass, &dbName)
	rows, e := db.Query("select * from items")
	ErrorCheck(e)
	i := itemStruct{}
	for rows.Next() {
		e = rows.Scan(&i.ItemID, &i.Type, &i.Name, &i.Weight, &i.Usage, &i.Description, &i.ShortItemURL, &i.FullItemURL)
		ErrorCheck(e)
		res = append(res, i)
	}
	defer db.Close()
	return res
}

func showItem(itemID uint16) (i itemStruct) {
	db := dbConn(&dbHost, &dbUser, &dbPass, &dbName)
	rows, e := db.Query("select * from items where itemid=?", itemID)
	ErrorCheck(e)
	for rows.Next() {
		e = rows.Scan(&i.ItemID, &i.Type, &i.Name, &i.Weight, &i.Usage, &i.Description, &i.ShortItemURL, &i.FullItemURL)
		ErrorCheck(e)
	}
	defer db.Close()
	return i
}

func showAllMonsters() (res []monsterStruct) {
	db := dbConn(&dbHost, &dbUser, &dbPass, &dbName)
	rows, e := db.Query("select * from monsters")
	ErrorCheck(e)
	m := monsterStruct{}
	for rows.Next() {
		e = rows.Scan(&m.MobID, &m.Level, &m.MaxHP, &m.Atk, &m.Def, &m.Hit, &m.Matk, &m.Mdef,
			&m.Flee, &m.Name, &m.Race, &m.Element, &m.Size, &m.Location, &m.STR, &m.VIT, &m.LUK,
			&m.INT, &m.DEX, &m.AGI, &m.ShortMobURL, &m.FullMobURL)
		ErrorCheck(e)
		res = append(res, m)
	}
	defer db.Close()
	return res
}

func showMonster(mobID uint16) (m monsterStruct) {
	db := dbConn(&dbHost, &dbUser, &dbPass, &dbName)
	rows, e := db.Query("select * from monsters where mobid=?", mobID)
	ErrorCheck(e)
	for rows.Next() {
		e = rows.Scan(&m.MobID, &m.Level, &m.MaxHP, &m.Atk, &m.Def, &m.Hit, &m.Matk, &m.Mdef,
			&m.Flee, &m.Name, &m.Race, &m.Element, &m.Size, &m.Location, &m.STR, &m.VIT, &m.LUK,
			&m.INT, &m.DEX, &m.AGI, &m.ShortMobURL, &m.FullMobURL)
		ErrorCheck(e)
	}
	defer db.Close()
	return m
}

func main() {
	//simpleMob := make(map[uint16]monsterStruct)
	//simpleMob[1] = monsterStruct{1, 15, 570, 16, 3, 151, 9, 16, 8, "Rocker", "Insect", "Earth", "Medium", "Prontera Field"}
	//simpleMob[2] = monsterStruct{2, 14, 338, 22, 0, 142, 9, 13, 7, "Savage Babe", "Brute", "Earth", "Small", "Geffen East Field, Prontera Field"}

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
	app.Get("/monster/{monsterid:uint16 max(65534)}", func(ctx iris.Context) {
		mobID := ctx.Params().GetUint16Default("monsterid", 0)
		monster := showMonster(mobID)
		ctx.JSON(monster)
	})

	app.Get("/monster/all", func(ctx iris.Context) {
		allmobs := showAllMonsters()
		ctx.JSON(allmobs)
	})
	// Method:    GET; Resource:  http://localhost:8080/item/42
	app.Get("/item/{itemid:uint16 max(65534)}", func(ctx iris.Context) {
		itemID := ctx.Params().GetUint16Default("itemid", 0)
		item := showItem(itemID)
		ctx.JSON(item)
	})

	app.Get("/item/all", func(ctx iris.Context) {
		allitems := showAllItems()
		ctx.JSON(allitems)
	})

	// app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
	//	ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
	//})
	// app.Get("/u/{username:string}", func(ctx iris.Context) {
	//	ctx.Writef("username (string): %s", ctx.Params().Get("username"))
	//})

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
