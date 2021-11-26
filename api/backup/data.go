package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"mysql_con"
)

//getGender
func getGender(ctx iris.Context) {
	var rst []map[string]string
	rst, _ = mysql_con.Query("select id as value,gender as label from pre_gender")
	if rst != nil {
		b, err := json.Marshal(rst)
		if err == nil {
			_, _ = ctx.JSON(string(b))
		}
	} else {
		ctx.Text("error")
	}
}

//getLocation
func getLocation(ctx iris.Context) {
	var rst []map[string]string
	rst, _ = mysql_con.Query("select id as value,location as label from pre_location")
	if rst != nil {
		b, err := json.Marshal(rst)
		if err == nil {
			_, _ = ctx.JSON(string(b))
		}
	} else {
		ctx.Text("error")
	}
}

//getLevel
func getLevel(ctx iris.Context) {

	posts := ctx.FormValue("posts")
fmt.Println("select level from pre_level where posts > " + posts + " ORDER BY posts LIMIT 1")
	var rst []map[string]string
	rst, _ = mysql_con.Query("select level from pre_level where posts > " + posts + " ORDER BY posts LIMIT 1")
	if rst != nil {
		b, err := json.Marshal(rst)
		if err == nil {
			_, _ = ctx.JSON(string(b))
		}
	} else {
		ctx.Text("error")
	}
}
