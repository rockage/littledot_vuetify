package main

import (
	"encoding/json"
	"mysql_con"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

type SubOrder struct {
	Amount        string `json:"amount"`
	Express       string `json:"freight_company"`
	ProductName   string `json:"item_describe"`
	SubOrderId    string `json:"id"`
	SubOrderPrice string `json:"price"`
	ProductId     string `json:"product_id"`
	ShippedDate   string `json:"shiped_date"`
	SubOrderState string `json:"state"`
	Tracking      string `json:"tracking_number"`
	Voltage       string `json:"voltage"`
}

type Order struct {
	Address    string     `json:"address"`
	Date       string     `json:"date"`
	Note       string     `json:"note"`
	OrderId    string     `json:"order_id"`
	Price      string     `json:"price"`
	OrderState string     `json:"state_id"`
	SubOrder   []SubOrder `json:"sub_order"`
	TbId       string     `json:"tb_id"`
	Vendor     string     `json:"vendor_id"`
}

func updateProductPosition(ctx iris.Context) { //更新产品显示位置
	type Position struct {
		Id       string `json:"id"`
		Position string `json:"position"`
	}

	var SQL string
	var new_position []Position
	str := ctx.FormValue("newPosition")
	json.Unmarshal([]byte(str), &new_position)
	SQL = "SELECT id,position FROM ld_products ORDER BY position"
	orgin_position, _ := mysql_con.Query(SQL)
	sql_update := ""
	sql_in := "WHERE id in ("
	for _, orgin := range orgin_position {
		for _, new := range new_position {
			if orgin["id"] == new.Id {
				if orgin["position"] != new.Position {
					sql_update += "WHEN " + orgin["id"] + " THEN " + new.Position + "\n"
					sql_in += orgin["id"] + ","
				} else {
					break
				}
			}
		}
	}
	if len(sql_update) > 0 {
		sql_in = string([]rune(sql_in)[:len(sql_in)-1]) // 消除最后一个逗号
		sql_in += ")"                                   // 加个)号
		SQL = "UPDATE ld_products \n" +
			"SET position = CASE ID\n" +
			sql_update +
			"END\n" +
			sql_in
		mysql_con.Exec(SQL)
	}
	ctx.Text("sucess")
}

func deleteProduct(ctx iris.Context) { //更新产品
	var SQL string
	SQL = "DELETE FROM ld_products WHERE id = " + ctx.FormValue("pid")
	mysql_con.Exec(SQL)
	SQL = "DELETE FROM ld_products_price WHERE pi = " + ctx.FormValue("pid")
	mysql_con.Exec(SQL)
	ctx.Text("sucess")
}

func updateProduct(ctx iris.Context) { //更新产品
	var SQL string
	pid := ctx.FormValue("pid")
	if pid == "" {
		SQL = "INSERT INTO ld_products " +
			"(name, tb_name, weight, class_id, position, note)" +
			" VALUES " +
			"('" + ctx.FormValue("pname") + "', '" + ctx.FormValue("tb_name") + "', '" + ctx.FormValue("weight") + "', " + ctx.FormValue("class") + ", 0 , " +
			"'" + ctx.FormValue("note") + "')"
		pid = mysql_con.Exec(SQL)
		SQL = "INSERT INTO ld_products_price (pi, vi, pr)  VALUES " + "(" + pid + ",1 , " + ctx.FormValue("price1") + ")"
		mysql_con.Exec(SQL)
		SQL = "INSERT INTO ld_products_price (pi, vi, pr)  VALUES " + "(" + pid + ",2 , " + ctx.FormValue("price2") + ")"
		mysql_con.Exec(SQL)

	} else {

		SQL = "UPDATE ld_products SET name = '" + ctx.FormValue("pname") + "'," +
			"tb_name = '" + ctx.FormValue("tb_name") + "'," +
			"weight = '" + ctx.FormValue("weight") + "'," +
			"class_id = " + ctx.FormValue("class") + "," +
			"note = '" + ctx.FormValue("note") + "'" +
			" WHERE id = " + pid
		mysql_con.Exec(SQL)
		SQL = "UPDATE ld_products_price SET pr =" + ctx.FormValue("price1") + " WHERE pi = " + ctx.FormValue("pid") + " AND vi = 1"
		mysql_con.Exec(SQL)
		SQL = "UPDATE ld_products_price SET pr =" + ctx.FormValue("price2") + " WHERE pi = " + ctx.FormValue("pid") + " AND vi = 2"
		mysql_con.Exec(SQL)
	}
	ctx.Text("sucess")
}

func getProductsClassList(ctx iris.Context) {
	var SQL string
	SQL = "select id, name, name_en, position from ld_products_class ORDER BY position"
	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}
}

func getProducts(ctx iris.Context) {
	var SQL string
	SQL = "SELECT id, name, tb_name, weight,note, position," +
		"(select pr from ld_products_price where ld_products_price.pi = ld_products.id and ld_products_price.vi = 1) as price1," +
		"(select pr from ld_products_price where ld_products_price.pi = ld_products.id and ld_products_price.vi = 2) as price2," +
		"(select name from ld_products_class where ld_products_class.ID = ld_products.class_id) as class FROM ld_products ORDER BY position"
	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}

}

func getShippedSubOrders(ctx iris.Context) { //发货视图
	var SQL string
	now := time.Now()
	keyword := ctx.FormValue("keyword")
	date_start := ctx.FormValue("date_start")
	date_end := ctx.FormValue("date_end")
	if date_start == "" {
		date_start = "2017-01-01 00:00:00"
	}
	if date_end == "" {
		date_end = now.AddDate(0, 0, 0).Format("2006-01-02 23:59:59")
	}

	SQL = "SELECT order_id, id, item_describe, voltage, amount, price, shiped_date,tracking_number," +
		"(select nickname from ld_vendor where ld_vendor.id = ld_order_suborder.vendor_id) as vendor," +
		"(select name from ld_express where ld_express.id = ld_order_suborder.freight_company) as express," +
		"(select note from ld_order where ld_order.id = ld_order_suborder.order_id) as note," +
		"(select address from ld_order where ld_order.id = ld_order_suborder.order_id) as address " +
		"FROM ld_order_suborder " +
		"WHERE " +
		"(item_describe LIKE '%" + keyword + "%'" +
		" OR " +
		"vendor_id = (select ID from ld_vendor where nickname = '" + keyword + "'))" +
		" AND " +
		"shiped_date BETWEEN '" + date_start + "' AND '" + date_end + "' " +
		" ORDER BY shiped_date DESC"

	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}

}
func updateLogistics(ctx iris.Context) { //更新发货数据
	var SQL string
	order_id := ctx.FormValue("order_id")
	suborder_id := ctx.FormValue("suborder_id")
	fc := ctx.FormValue("freight_company")
	tn := ctx.FormValue("tracking_number")
	sd := ctx.FormValue("shiped_date")
	if fc == "" {
		fc = "17"

	}
	if sd == "" {
		sd = "1900-01-01 00:00:00"
	}
	var state string
	if fc != "17" && tn != "" && sd != "1900-01-01 00:00:00" {
		state = ",state = 3" //3 = 已发货
	}
	SQL = "UPDATE ld_order_suborder SET freight_company = '" + fc + "'," +
		"tracking_number = '" + tn + "'," +
		"shiped_date = '" + sd + "'" + state +
		" WHERE id = " + suborder_id
	mysql_con.Exec(SQL)

	SQL = "SELECT state FROM ld_order_suborder where order_id = " + order_id
	var err error
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)

	if err == nil {
		var pre_state string
		var change bool
		change = true
		pre_state = rst[0]["state"]
		for _, value := range rst {
			if pre_state != value["state"] {
				change = false
			}
		}
		if change == true {
			SQL = "UPDATE ld_order SET state_id = " + pre_state + " WHERE id = " + order_id
			mysql_con.Exec(SQL)
		}
	}

	ctx.Text("sucess")
}

func updatePackets(ctx iris.Context) { //打包
	var SQL string
	suborder_id := ctx.FormValue("suborder_id")
	state_id := ctx.FormValue("state_id")

	SQL = "UPDATE ld_order_suborder SET state = " + state_id +
		" WHERE id = " + suborder_id
	mysql_con.Exec(SQL) // state = 4 : 已打包 state = 2 : 解除打包

	SQL = "SELECT state FROM ld_order_suborder where order_id = (select order_id from ld_order_suborder where id =" + suborder_id + ")"
	var err error
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)

	if err == nil {
		var pre_state string
		var change bool
		change = true
		pre_state = rst[0]["state"]
		for _, value := range rst {
			if pre_state != value["state"] {
				change = false
			}
		}
		if change == true {
			SQL = "UPDATE ld_order SET state_id = " + pre_state + " WHERE id = (select order_id from ld_order_suborder where id =" + suborder_id + ")"
			mysql_con.Exec(SQL) // 全部打包完毕
		}
	}

	ctx.Text("sucess")
}

func getLogistics(ctx iris.Context) { //获取发货数据
	var SQL string
	suborder_id := ctx.FormValue("suborder_id")
	SQL = "SELECT id, item_describe, freight_company, tracking_number, shiped_date FROM ld_order_suborder where id = " + suborder_id
	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}
}

func deleteOrder(ctx iris.Context) { //删除订单
	var SQL string
	order_id := ctx.FormValue("order_id")
	SQL = "DELETE FROM ld_order WHERE id = " + order_id
	mysql_con.Exec(SQL)
	SQL = "DELETE FROM ld_order_suborder WHERE order_id = " + order_id
	mysql_con.Exec(SQL)
	ctx.Text("sucess")
}

func updateOrder(ctx iris.Context) { //更新订单（包括所属子订单）
	var SQL string
	var order Order
	data := ctx.FormValue("data")
	json.Unmarshal([]byte(data), &order)
	if order.OrderId == "" {
		SQL = "INSERT INTO ld_order " +
			"(date, vendor_id, tb_id, price, state_id, note, address)" +
			" VALUES " +
			"('" + order.Date + "', " + order.Vendor + ", '" + order.TbId + "', " + order.Price + ", " +
			order.OrderState + ", '" + order.Note + "', '" + order.Address + "')"
		order.OrderId = mysql_con.Exec(SQL)
		for _, value := range order.SubOrder {
			if value.ShippedDate == "" {
				value.ShippedDate = "1900-01-01 00:00:00" //设置一个具体的值，否则MySQL会报： Incorrect datetime value 错误
			}
			SQL = "INSERT INTO ld_order_suborder " +
				"(order_id, product_id, vendor_id, item_describe, voltage, amount, price, shiped_date, freight_company, tracking_number, state)" +
				" VALUES " +
				"(" + order.OrderId + ", " + value.ProductId + ", " + order.Vendor + ", '" + value.ProductName + "', " +
				value.Voltage + ", " + value.Amount + ", " + value.SubOrderPrice + ", '" + value.ShippedDate + "', " +
				value.Express + ", '" + value.Tracking + "', " + value.SubOrderState + ")"
			ret := mysql_con.Exec(SQL)
			ctx.Text(ret)
		}
	} else {
		SQL = "UPDATE ld_order SET date = '" + order.Date + "'," +
			"vendor_id = " + order.Vendor + "," +
			"tb_id = '" + order.TbId + "'," +
			"price = " + order.Price + "," +
			"state_id = " + order.OrderState + "," +
			"note = '" + order.Note + "'," +
			"address = '" + order.Address + "' WHERE id = " + order.OrderId
		mysql_con.Exec(SQL)
		SQL = "DELETE FROM ld_order_suborder WHERE order_id = " + order.OrderId //重构子订单：先删除与主订单的所有连接
		mysql_con.Exec(SQL)

		//重建子订单：
		for _, value := range order.SubOrder {
			if value.ShippedDate == "" {
				value.ShippedDate = "1900-01-01 00:00:00"
			}
			SQL = "INSERT INTO ld_order_suborder " +
				"(order_id, product_id, vendor_id, item_describe, voltage, amount, price, shiped_date, freight_company, tracking_number, state)" +
				" VALUES " +
				"(" + order.OrderId + ", " + value.ProductId + ", " + order.Vendor + ", '" + value.ProductName + "', " +
				value.Voltage + ", " + value.Amount + ", " + value.SubOrderPrice + ", '" + value.ShippedDate + "', " +
				value.Express + ", '" + value.Tracking + "', " + value.SubOrderState + ")"
			ret := mysql_con.Exec(SQL)
			ctx.Text(ret)
		}

	}

}

func getOrderSubOrders(ctx iris.Context) { //获取子订单数据
	var SQL string
	order_id := ctx.FormValue("order_id")
	SQL = "SELECT * FROM ld_order_suborder where order_id = " + order_id
	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}
}

func getOrders(ctx iris.Context) { //获取订单数据
	var SQL string
	now := time.Now()
	keyword := ctx.FormValue("keyword")
	date_start := ctx.FormValue("date_start")
	date_end := ctx.FormValue("date_end")

	if date_start == "" {
		date_start = "2017-01-01 00:00:00"
	}
	if date_end == "" {
		date_end = now.AddDate(0, 0, 0).Format("2006-01-02 23:59:59")
	}

	var search string
	if len(keyword) > 0 {
		search = "(address like '%" + keyword + "%' OR note like '%" + keyword + "%' OR tb_id like '%" + keyword + "%') AND "

	} else {
		search = "(state_id = 2 OR state_id = 4) OR "
	}

	//state_id： 2，4 不受时间条件限制
	//state_id显示优先级 : 2(待) > 4(包) > 1(未) [3(已)和12(关)不优先显示, 以时间排序]

	SQL = "SELECT id, date, hide, price, tb_id, note, address, state_id," +
		"concat(ld_order.vendor_id,'|',(Select `nickname` from `ld_vendor` where ld_order.vendor_id = ld_vendor.id)) As vendor," +
		"(SELECT group_concat(id,'|'," +
		"item_describe,'|'," +
		"amount,'|', " +
		"(Select `name` from `ld_voltage` where ld_voltage.id = ld_order_suborder.voltage),'|', " +
		"state,'|', " +
		"(Select `name` from `ld_order_state` where ld_order_state.id = ld_order_suborder.state),'|', " +
		"price) " +
		"FROM ld_order_suborder " +
		"where ld_order_suborder.order_id = ld_order.id) as p_info " +
		"FROM `ld_order` WHERE " +
		search +
		"date between '" + date_start + "' and '" + date_end + "' " +
		"ORDER BY field(`state_id`,1,4,2) DESC, `date` DESC"

	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}
}

func getOrdersForViewer(ctx iris.Context) { // 为 “今日发货列表” 提供数据，该页面独立存在，并不是Vue主程序的一部分
	var SQL string
	SQL = `SELECT id, date,  note, address, state_id, 
(Select nickname from ld_vendor where ld_order.vendor_id = ld_vendor.id) As vendor,
(SELECT group_concat(id,'|',item_describe,'|',amount,'|', (Select name from ld_voltage where ld_voltage.id = ld_order_suborder.voltage),'|', state,'|', 
(Select name from ld_order_state where ld_order_state.id = ld_order_suborder.state),'|', price) FROM ld_order_suborder where ld_order_suborder.order_id = ld_order.id) as p_info 
FROM ld_order WHERE state_id = 2 OR state_id = 4 ORDER BY field(state_id, 4, 2) DESC, date DESC`
	var err error
	var b []byte
	var rst []map[string]string
	rst, err = mysql_con.Query(SQL)
	if err == nil {
		b, err = json.Marshal(rst)
		if err == nil {
			ctx.JSON(string(b))
		}
	}
}

func getDefaultList(ctx iris.Context) { //获取默认列表
	var returnValue [6]string
	returnValue[0] = getStateList()
	returnValue[1] = getVendorList()
	returnValue[2] = getVoltageList()
	returnValue[3] = getProductList()
	returnValue[4] = getExpressList()
	returnValue[5] = getPriceList()
	ctx.JSON(returnValue)
}

func getStateList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT id,name FROM ld_order_state")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getVendorList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT id,nickname,class,address,taobao,contacts,tel FROM ld_vendor WHERE class = 1 OR class = 2 OR class = 3 OR class = 4 ORDER BY position")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getVoltageList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT id,name FROM ld_voltage")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getProductList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT id,name FROM ld_products ORDER BY position")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getExpressList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT id,name FROM ld_express ORDER BY position")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getPriceList() string {
	var rst []map[string]string
	rst, _ = mysql_con.Query("SELECT pi, vi, pr FROM ld_products_price")
	b, _ := json.Marshal(rst)
	return string(b)
}

func getStatistics(ctx iris.Context) { // 为销售图表模块提供数据
	type Chart struct {
		Name string `json:"name"`
		Data []int  `json:"data"`
	}

	type Bar struct {
		Name []string `json:"name"`
		Data []int    `json:"data"`
	}

	type Bar_L_R struct {
		Title    string
		Total    string
		WaitSend string
		Bars_L   Bar
		Bars_R   Bar
	}

	queryBars := func(start string, end string) (Bar, Bar) {
		/*
			Bar的数据结构：
			{
				Name:["雷音","易潇","尼古"]
				Data:[1,2,3]
			}
		*/

		sql_vendors := `
			SELECT SUM(price) AS PRICE,
			   (SELECT nickname FROM ld_vendor WHERE ld_order_suborder.vendor_id = ld_vendor.id) AS VENDOR
			FROM ld_order_suborder WHERE shiped_date BETWEEN '` + start + `' AND '` + end + `'
	        GROUP BY vendor_id
	        ORDER BY PRICE DESC
	        LIMIT 15`

		// 左侧柱状图：Bars_L (经销商)
		var rst []map[string]string
		rst, _ = mysql_con.Query(sql_vendors) // 查询经销商
		var one_vendor Bar
		one_vendor.Name = make([]string, len(rst))
		one_vendor.Data = make([]int, len(rst))
		for k, v := range rst {
			one_vendor.Name[k] = v["VENDOR"]
			value_int, _ := strconv.Atoi(v["PRICE"])
			one_vendor.Data[k] = value_int
		}
		// 右侧柱状图：Bars_R (产品)
		sql_products := `
		    SELECT SUM(price) AS PRICE,
			   (SELECT name FROM ld_products WHERE ld_order_suborder.product_id = ld_products.id) AS PRODUCT
			FROM ld_order_suborder WHERE shiped_date BETWEEN '` + start + `' AND '` + end + `'
			GROUP BY product_id
			ORDER BY PRICE DESC
			LIMIT 15`
		rst, _ = mysql_con.Query(sql_products)
		var one_product Bar
		one_product.Name = make([]string, len(rst))
		one_product.Data = make([]int, len(rst))
		for k, v := range rst {
			one_product.Name[k] = v["PRODUCT"]
			value_int, _ := strconv.Atoi(v["PRICE"])
			one_product.Data[k] = value_int
		}

		return one_vendor, one_product
	}

	op := ctx.FormValue("op")
	start := ctx.FormValue("start")
	end := ctx.FormValue("end")
	var start_year, end_year int
	switch op {
	case "charts": // 主曲线：单位（年）
		//var charts []Chart
		var rst []map[string]string
		sql_charts := `SELECT DATE_FORMAT(shiped_date, '%Y-%m') as DATE,
			 				  sum(price) as PRICE 
					   FROM  ld_order_suborder 
					   WHERE state = 3 AND shiped_date BETWEEN '` + start + `' AND '` + end + `'
					   GROUP BY DATE_FORMAT(shiped_date, '%Y-%m')
					   ORDER BY DATE_FORMAT(shiped_date, '%Y') DESC, DATE_FORMAT(shiped_date, '%m')`
		// GROUP BY: 年-月
		// ORDER BY: 年逆序, 月顺序
		rst, _ = mysql_con.Query(sql_charts)
		/*
			chart的数据结构：
			[{
				name:"2020",
				data:[1,2,3,4,5,6,7,8,9,10,11,12]
			}, ...]
		*/
		var charts []Chart
		var one Chart
		var pre_date string = ""
		var last_rec int = len(rst) - 1
		for k, v := range rst {
			date := string([]rune(v["DATE"])[:4]) // DATE左取4位：2021-09 -> 2021
			value, _ := strconv.Atoi(v["PRICE"])

			if k == 0 { // 特殊情况：k=0只设置name
				one.Name = date
			}

			if k == last_rec { // 特殊情况2：k=最后一条数据，将最后一条数据串入后跳出循环
				one.Data = append(one.Data, value)
				charts = append(charts, one)
				break
			}

			if pre_date != date && k != 0 { // 正常情况：年份不同 && 不是首次进入
				charts = append(charts, one) // 串入本年度各月份数据
				one.Data = make([]int, 0)    // 清空切片
				one.Name = date              // 开启一个新的年度
			}

			one.Data = append(one.Data, value)
			pre_date = date

		}
		b, _ := json.Marshal(charts)
		s := string(b)
		ctx.Text(s)

	case "bars_years": // 统计柱状图：单位（年）
		var bars []Bar_L_R
		start_year, _ = strconv.Atoi(string([]rune(start)[:4]))
		end_year, _ = strconv.Atoi(string([]rune(end)[:4]))
		for i := start_year; i <= end_year; i++ {
			s := strconv.Itoa(i) + "-01-01 00:00:00"
			e := strconv.Itoa(i) + "-12-31 23:59:59"
			var one Bar_L_R
			one.Title = strconv.Itoa(i)
			sql_total := `SELECT SUM(price) AS PRICE
						FROM ld_order_suborder WHERE shiped_date BETWEEN '` + s + `' AND '` + e + `'`
			rst, _ := mysql_con.Query(sql_total) // 查询全年统计
			one.Total = rst[0]["PRICE"]
			one.Bars_L, one.Bars_R = queryBars(s, e) // 查询左：经销商销量，右：产品销量
			bars = append(bars, one)
		}
		b, _ := json.Marshal(bars)
		s := string(b)
		ctx.Text(s)

	case "bars_single_month": // 统计柱状图：单位（月）
		var one Bar_L_R
		one.Title = start
		one.Bars_L, one.Bars_R = queryBars(start, end)
		sql_total := `SELECT SUM(price) AS PRICE
						FROM ld_order_suborder WHERE shiped_date BETWEEN '` + start + `' AND '` + end + `'`
		rst, _ := mysql_con.Query(sql_total) // 查询当月统计
		one.Total = rst[0]["PRICE"]

		sql_wait_send := "SELECT sum(price) as WAITSEND from ld_order WHERE state_id = 2"
		rst, _ = mysql_con.Query(sql_wait_send) // 查询全年统计
		one.WaitSend = rst[0]["WAITSEND"]

		b, _ := json.Marshal(one)
		s := string(b)
		ctx.Text(s)

	default:
	}

}

func getVendors(ctx iris.Context) { // 获取经销商列表
	SQL := `SELECT * FROM ld_vendor ORDER BY position`
	rst, _ := mysql_con.Query(SQL) // 查询当月统计
	b, _ := json.Marshal(rst)
	s := string(b)
	ctx.Text(s)
}

func updateVendors(ctx iris.Context) { // 修改经销商数据
	type Position struct {
		Id       string `json:"ID"`
		Position string `json:"position"`
	}

	type UserClass struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	type UserInfo struct {
		Id       string    `json:"id"`
		Nickname string    `json:"nickname"`
		Sale     string    `json:"sale"`
		Fullname string    `json:"fullname"`
		Contacts string    `json:"contacts"`
		Tel      string    `json:"tel"`
		Address  string    `json:"address"`
		Email    string    `json:"email"`
		Wechat   string    `json:"wechat"`
		Taobao   string    `json:"taobao"`
		QQ       string    `json:"qq"`
		Alipay   string    `json:"alipay"`
		Bank     string    `json:"bank"`
		Class    UserClass `json:"class"`
	}

	var userinfo UserInfo

	switch ctx.FormValue("op") {
	case "updatePosition":
		var new_position []Position
		str := ctx.FormValue("newPosition")
		json.Unmarshal([]byte(str), &new_position)
		SQL := "SELECT ID,position FROM ld_vendor ORDER BY position"
		orgin_position, _ := mysql_con.Query(SQL)
		sql_update := ""
		sql_in := "WHERE ID in ("
		for _, orgin := range orgin_position {
			for _, new := range new_position {
				if orgin["ID"] == new.Id {
					if orgin["position"] != new.Position {
						sql_update += "WHEN " + orgin["ID"] + " THEN " + new.Position + "\n"
						sql_in += orgin["ID"] + ","
					} else {
						break
					}
				}
			}
		}
		if len(sql_update) > 0 {
			sql_in = string([]rune(sql_in)[:len(sql_in)-1]) // 消除最后一个逗号
			sql_in += ")"                                   // 加个)号
			SQL = "UPDATE ld_vendor \n" +
				"SET position = CASE ID\n" +
				sql_update +
				"END\n" +
				sql_in
			mysql_con.Exec(SQL)
		}
		ctx.Text("sucess")
		return
	case "edit":
		str := ctx.FormValue("datas")
		json.Unmarshal([]byte(str), &userinfo)

		SQL := `UPDATE ld_vendor SET ` +
			` nickname =  '` + userinfo.Nickname + `',` +
			` sale =  '` + userinfo.Sale + `',` +
			` fullname =  '` + userinfo.Fullname + `',` +
			` contacts =  '` + userinfo.Contacts + `',` +
			` tel =  '` + userinfo.Tel + `',` +
			` address =  '` + userinfo.Address + `',` +
			` email =  '` + userinfo.Email + `',` +
			` wechat =  '` + userinfo.Wechat + `',` +
			` taobao =  '` + userinfo.Taobao + `',` +
			` qq =  '` + userinfo.QQ + `',` +
			` alipay =  '` + userinfo.Alipay + `',` +
			` bank =  '` + userinfo.Bank + `',` +
			` class =  '` + userinfo.Class.Id + `' WHERE id = ` + userinfo.Id
		mysql_con.Exec(SQL)
		ctx.Text("sucess")
		return
	case "addnew":
		str := ctx.FormValue("datas")
		json.Unmarshal([]byte(str), &userinfo)
		SQL := "INSERT INTO ld_vendor " +
			"(nickname, sale, fullname, contacts, tel, address, email, wechat, taobao, qq, alipay, bank, position, class)" +
			" VALUES " +
			"('" + userinfo.Nickname + "'," +
			"'" + userinfo.Sale + "'," +
			"'" + userinfo.Fullname + "'," +
			"'" + userinfo.Contacts + "'," +
			"'" + userinfo.Tel + "'," +
			"'" + userinfo.Address + "'," +
			"'" + userinfo.Email + "'," +
			"'" + userinfo.Wechat + "'," +
			"'" + userinfo.Taobao + "'," +
			"'" + userinfo.QQ + "'," +
			"'" + userinfo.Alipay + "'," +
			"'" + userinfo.Bank + "'," +
			" 0," + // 默认位置为0
			"'" + userinfo.Class.Id + "')"
		mysql_con.Exec(SQL)
		ctx.Text("sucess")
		return
	case "delete":
		SQL := "DELETE FROM ld_vendor WHERE id = " + ctx.FormValue("id")
		mysql_con.Exec(SQL)
		ctx.Text("sucess")
		return
	default:

	}

}
