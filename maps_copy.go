package main

import (
	"github.com/kataras/iris/v12"
)

// User is just a bindable object structure.
// type User struct {
// 	Username  string `json:"username"`
// 	Firstname string `json:"firstname"`
// 	Lastname  string `json:"lastname"`
// 	City      string `json:"city"`
// 	Age       int    `json:"age"`
// }
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

func main() {
	app := iris.New()
	// Method POST: http://localhost:8080/decode
	// app.Post("/decode", func(ctx iris.Context) {
	// 	var mob mobStruct
	// 	ctx.ReadJSON(&mob)
	//  ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
	// })

	// Method GET: http://localhost:8080/encode
	app.Get("/encode", func(ctx iris.Context) {
		simpleMob := make(map[uint8]mobStruct)
		simpleMob[1] = mobStruct{1, 15, 570, 16, 3, 151, 9, 16, 8, "Rocker", "Insect", "Earth", "Medium", "Prontera Field"}
		simpleMob[2] = mobStruct{2, 14, 338, 22, 0, 142, 9, 13, 7, "Savage Babe", "Brute", "Earth", "Small", "Geffen East Field, Prontera Field"}
		//doe := User{"Johndoe2222", "John", "Doe", "Neither FBI knows!!!", 25}
		// doe2 := User{
		// 	Username:  "Johndoe",
		// 	Firstname: "John",
		// 	Lastname:  "Doe",
		// 	City:      "Neither FBI knows!!!",
		// 	Age:       25,
		// }
		ctx.JSON(simpleMob[1])
	})
	// Listen for incoming HTTP/1.x & HTTP/2 clients on localhost port 8080.
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}
