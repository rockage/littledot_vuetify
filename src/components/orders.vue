<template>
  <div class="text-caption">
    <!-- 订单弹窗 -->
    <v-dialog v-model="dialog" width="80%">
      <template v-slot:activator="{ on, attrs }">
        <div v-bind="attrs" v-on="on" ref="dialog"></div>
      </template>
      <v-card class="mx-auto">
        <v-card-title>
          <span class="headline"> 订单{{ this.dialog_data.order_id }} </span>
        </v-card-title>
        <v-card-text>
          <v-container class="text-caption">
            <v-row>
              <v-col :cols="3">
                <v-text-field
                  dense
                  v-model="dialog_data.date"
                  label="下单时间"
                ></v-text-field>
              </v-col>
              <v-col :cols="3">
                <v-text-field
                  dense
                  v-model="dialog_data.tb_id"
                  label="淘宝订单号"
                ></v-text-field>
              </v-col>
              <v-col :cols="2">
                <v-select
                  dense
                  v-model="dialog_data.vendor_id"
                  :items="dialog_data.vendor_list"
                  item-text="nickname"
                  item-value="id"
                  label="商户"
                  @change="onVendorChange($event)"
                ></v-select>
              </v-col>
              <v-col>
                <v-select
                  dense
                  v-model="dialog_data.state_id"
                  :items="dialog_data.state_list"
                  item-text="name"
                  item-value="id"
                  label="状态"
                ></v-select>
              </v-col>
              <v-col>
                <v-text-field
                  dense
                  v-model="dialog_data.price"
                  label="订单总价"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-textarea
                  dense
                  name="input-7-1"
                  filled
                  label="备注"
                  v-model="dialog_data.note"
                  rows="1"
                  row-height="20"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-textarea
                  dense
                  name="input-7-1"
                  filled
                  label="地址"
                  v-model="dialog_data.address"
                  rows="2"
                  row-height="30"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-btn
              icon
              color="#00B0FF"
              @click="onAddSubOrder(dialog_data.sub_order)"
            >
              <v-icon> mdi-plus </v-icon>
            </v-btn>
            <v-card
              dense
              elevation="6"
              outlined
              dense
              style="
                padding-left: 10px;
                padding-right: 10px;
                padding-top: 15px;
                margin-bottom: 10px;
              "
              v-for="item in dialog_data.sub_order"
            >
              <div class="box" style="flex-direction: row">
                <div>
                  <div class="box" style="flex-direction: row; flex-wrap: wrap">
                    <v-autocomplete
                      dense
                      v-model="item.product_list.select"
                      :items="item.product_list.items"
                      item-text="name"
                      item-value="id"
                      label="产品"
                      @blur="onProductListBlur($event, item)"
                      @change="updatePrice()"
                      style="margin-right: 5px"
                      return-object
                    ></v-autocomplete>
                    <v-select
                      dense
                      v-model="item.voltage"
                      :items="dialog_data.voltage_list"
                      item-text="name"
                      item-value="id"
                      label="电压"
                      style="max-width: 120px; margin-right: 5px"
                    ></v-select>
                    <v-text-field
                      dense
                      v-model="item.amount"
                      label="数量"
                      @input="updatePrice()"
                      style="max-width: 60px; margin-right: 5px"
                    ></v-text-field>
                    <v-text-field
                      dense
                      v-model="item.price"
                      label="价格"
                      @input="updateTotalPrice()"
                      style="max-width: 60px; margin-right: 5px"
                    ></v-text-field>
                    <v-text-field
                      v-model="item.shiped_date"
                      dense
                      label="发货日期"
                      style="max-width: 160px; margin-right: 5px"
                    ></v-text-field>
                    <v-select
                      dense
                      v-model="item.freight_company"
                      :items="dialog_data.express_list"
                      item-text="name"
                      item-value="id"
                      label="快递公司"
                      style="max-width: 120px; margin-right: 5px"
                    ></v-select>
                    <v-text-field
                      v-model="item.tracking_number"
                      dense
                      label="运单号"
                      style="max-width: 160px; margin-right: 5px"
                    ></v-text-field>
                    <v-select
                      dense
                      v-model="item.state"
                      :items="dialog_data.state_list"
                      item-text="name"
                      item-value="id"
                      label="状态"
                      style="max-width: 60px; margin-right: 5px"
                    ></v-select>
                  </div>
                </div>
                <div
                  class="box"
                  style="flex-direction: row; align-items: center"
                >
                  <v-btn icon color="pink" @click="onDelSubOrder(item)">
                    <v-icon> mdi-close-circle </v-icon>
                  </v-btn>
                  <div class="sm" style="color: #00897b">
                    {{ item.id }}
                  </div>
                </div>
              </div>
            </v-card>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-btn text color="pink" @click="onDelOrder()"> DELETE </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="dialog = false">
            Close
          </v-btn>
          <v-btn color="blue darken-1" text @click="updateOrder()">
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!-- 发货弹窗 -->
    <v-dialog v-model="dialog2" width="50%">
      <template v-slot:activator="{ on, attrs }">
        <div v-bind="attrs" v-on="on" ref="dialog2"></div>
      </template>
      <v-card class="mx-auto">
        <v-card-title>
          <span class="text-subtitle-1">
            订单{{ this.dialog2_data.sub_order_id }}:
            {{ this.dialog2_data.item_describe }}
          </span>
        </v-card-title>
        <v-card-text>
          <v-container class="text-caption">
            <v-row>
              <v-select
                dense
                v-model="dialog2_data.express"
                :items="dialog_data.express_list"
                item-text="name"
                item-value="id"
                label="快递公司"
                style="margin-right: 5px"
                full-width
              ></v-select>
            </v-row>
            <v-row>
              <v-text-field
                dense
                v-model="dialog2_data.tracking"
                label="快递单号"
                id="clipboard"
              ></v-text-field>
            </v-row>
            <v-row>
              <v-text-field
                dense
                v-model="dialog2_data.shipped_date"
                label="发货时间"
              ></v-text-field>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-btn text color="pink" @click="onDelOrder()"> DELETE </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="dialog2 = false">
            Close
          </v-btn>
          <v-btn
            color="blue darken-1"
            text
            @click="updateLogistics(dialog2_data)"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!-- 主界面 ####-->
    <div class="box">
      <v-menu offset-y style="max-width: 100px">
        <template v-slot:activator="{ attrs, on }">
          <v-btn
            small
            color="primary"
            class="white--text ma-1"
            v-bind="attrs"
            v-on="on"
          >
            视图
          </v-btn>
        </template>
        <v-list>
          <v-list-item link>
            <v-list-item-title @click="viewModel = 0"
              >订单视图</v-list-item-title
            >
          </v-list-item>

          <v-list-item link>
            <v-list-item-title @click="viewModel = 1"
              >发货视图</v-list-item-title
            >
          </v-list-item>
        </v-list>
      </v-menu>
      <v-spacer></v-spacer>
      <v-card style="max-height: 50px; margin: 5px">
        <v-card-actions>
          <v-autocomplete
            dense
            v-model="discount_product"
            :items="dialog_data.product_list.items"
            item-text="name"
            item-value="id"
            @change="updateDiscountPrice()"
            style="margin-right: 5px; max-width: 200px; font-size: xx-small"
            return-object
          ></v-autocomplete>
          <v-select
            dense
            v-model="discount_pcs"
            :items="discount_pcs_list"
            item-text="name"
            item-value="id"
            style="max-width: 50px; text-align: center; font-size: xx-small"
            :value="discount_pcs"
            @change="updateDiscountPrice()"
            small
            return-object
          ></v-select>
          <v-btn text style="margin-bottom: 10px"> {{ discount }} </v-btn>
        </v-card-actions>
      </v-card>
      <v-card style="max-height: 50px; margin: 5px">
        <v-card-actions>
          <v-menu
            v-model="date_menu1"
            :close-on-content-click="false"
            :nudge-right="0"
            transition="scale-transition"
            offset-y
            min-width="auto"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-text-field
                dense
                v-model="date_start"
                v-bind="attrs"
                v-on="on"
                style="max-width: 200px; font-size: xx-small"
                small
              ></v-text-field>
            </template>
            <v-date-picker
              v-model="date_start"
              locale="zh-cn"
              @input="date_menu1 = false"
            ></v-date-picker>
          </v-menu>
          <v-icon>mdi-chevron-right</v-icon>
          <v-menu
            v-model="date_menu2"
            :close-on-content-click="false"
            :nudge-right="-30"
            transition="scale-transition"
            offset-y
            min-width="auto"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-text-field
                dense
                v-model="date_end"
                v-bind="attrs"
                v-on="on"
                style="width: 200px; font-size: xx-small"
                small
              ></v-text-field>
            </template>
            <v-date-picker
              v-model="date_end"
              locale="zh-cn"
              @input="date_menu2 = false"
            ></v-date-picker>
          </v-menu>

          <v-text-field
            dense
            single-line
            small
            v-model="keyword"
            placeholder="关键字"
            style="width: 300px; font-size: xx-small"
          >
            <template v-slot:append-outer>
              <v-btn small outlined color="indigo" @click="initView()">
                <v-icon> mdi-magnify </v-icon>
              </v-btn>
            </template>
          </v-text-field>
          <v-menu offset-y style="max-width: 100px">
            <template v-slot:activator="{ attrs, on }">
              <v-btn
                class="mx-2"
                fab
                dark
                small
                text
                color="primary"
                style="margin-bottom: 5px"
                v-bind="attrs"
                v-on="on"
              >
                <v-icon>mdi-calendar-range</v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-item link>
                <v-list-item-title @click="range_default()">
                  默认时间段</v-list-item-title
                >
              </v-list-item>

              <v-list-item link>
                <v-list-item-title @click="range_whole()"
                  >所有时间段
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-card-actions>
      </v-card>
    </div>
    <!-- 订单视图 #### -->
    <div class="divTable" v-if="view_model == 0" style="font-size: small">
      <div class="divTableHeading">
        <div class="divTableCell center nowrap">
          <v-btn icon color="#00B0FF" @click="onAddNewOrder()" x-small>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </div>
        <div class="divTableCell center nowrap">日期/TB编号</div>
        <div class="divTableCell center nowrap">商家</div>
        <div class="divTableCell center nowrap">产品列表</div>
        <div class="divTableCell center nowrap">数量</div>
        <div class="divTableCell center nowrap">状态</div>
        <div class="divTableCell center nowrap">金额</div>
        <div class="divTableCell center nowrap">备注</div>
        <div class="divTableCell center nowrap">地址</div>
      </div>
      <div class="divTableBody">
        <div
          class="divTableRow"
          v-for="item in orders"
          :style="{ color: item.color }"
        >
          <div class="divTableCell">
            <v-btn
              x-small
              class="ma-2 divTableCell"
              outlined
              color="indigo"
              style="margin-left: 0px"
              @click="onOrderClick"
            >
              {{ item.id }}
            </v-btn>
          </div>
          <div class="divTableCell">{{ item.date }} {{ item.tb_id }}</div>
          <div class="divTableCell nowrap">
            {{ item.vendor }}
          </div>
          <div class="divTableCell nowrap" style="text-align: right">
            <span
              v-for="p in item.sub_order"
              style="white-space: pre-line"
              :style="{
                color: p.color,
              }"
              @mouseenter="onMouseHover"
              @mouseleave="onMouseLeave"
              @click="onMouseClick(p.sub_order_id, item.tb_id)"
            >
              {{ p.item_describe }}
              {{ p.voltage }}
            </span>
          </div>
          <div class="divTableCell center nowrap">
            <span
              v-for="p in item.sub_order"
              style="white-space: pre-line"
              :style="{
                color: p.color,
              }"
            >
              {{ p.amount }}
            </span>
          </div>
          <div class="divTableCell center nowrap">
            <span
              v-for="p in item.sub_order"
              style="white-space: pre-line; text-align: center"
              :style="{
                color: p.color,
              }"
            >
              {{ p.state }}
            </span>
          </div>
          <div class="divTableCell center nowrap">
            <span
              v-for="p in item.sub_order"
              style="white-space: pre-line; text-align: center"
              :style="{
                color: p.color,
              }"
            >
              {{ p.price }}
            </span>
          </div>
          <div class="divTableCell">
            {{ item.note }}
          </div>
          <div class="divTableCell">
            {{ item.address }}
          </div>
        </div>
      </div>
    </div>
    <!-- 发货视图 #####-->
    <div class="divTable" v-if="view_model == 1">
      <div class="divTableHeading">
        <div class="divTableCell center nowrap">
          <v-btn icon color="#00B0FF" @click="onAddNewOrder()" x-small>
            OID/PID
          </v-btn>
        </div>
        <div class="divTableCell center nowrap">商家</div>
        <div class="divTableCell center nowrap">产品</div>
        <div class="divTableCell center nowrap">电压</div>
        <div class="divTableCell center nowrap">数量</div>
        <div class="divTableCell center nowrap">金额</div>
        <div class="divTableCell center nowrap">发货时间</div>
        <div class="divTableCell center nowrap">快递</div>
        <div class="divTableCell center nowrap">单号</div>
        <div class="divTableCell center nowrap">备注</div>
        <div class="divTableCell center nowrap">地址</div>
      </div>
      <div class="divTableBody">
        <div class="divTableRow" v-for="item in shipped_orders">
          <div class="divTableCell">{{ item.order_id }}/{{ item.id }}</div>
          <div class="divTableCell nowrap">
            {{ item.vendor }}
          </div>
          <div class="divTableCell nowrap">
            {{ item.item_describe }}
          </div>

          <div class="divTableCell center nowrap">
            {{ item.voltage }}
          </div>
          <div class="divTableCell center nowrap">
            {{ item.amount }}
          </div>
          <div class="divTableCell center nowrap">
            {{ item.price }}
          </div>
          <div class="divTableCell">
            {{ item.shiped_date }}
          </div>
          <div class="divTableCell">
            {{ item.tracking_number }}
          </div>
          <div class="divTableCell">
            {{ item.express }}
          </div>
          <div class="divTableCell">
            {{ item.note }}
          </div>
          <div class="divTableCell">
            {{ item.address }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'orders',
  data() {
    return {
      discount_product: { name: 'MK2 6J1(普)', id: '2' },
      discount_pcs_list: [],
      discount_pcs: { name: '1', id: 1 },
      discount: '-160', //默认数据：一台MK2, 折扣160
      date_start: '',
      date_end: '',
      date_menu1: false,
      date_menu2: false,
      keyword: '',
      view_model: 0,
      localStorage: window.localStorage,
      old_element: Object,
      old_color: '',
      updateLogistics_tbid: '', // 打开淘宝页面时候的订单ID
      orders: [
        //主页面呈现
        {
          id: '',
          state_id: '',
          date: '',
          vendor: '',
          sub_order: [
            {
              sub_order_id: '',
              item_describe: '',
              amount: '',
              voltage: '',
              color: '',
              state: '',
              price: '',
              order_id: '',
            },
          ],
          note: '',
          address: '',
        },
      ],
      shipped_orders: [
        //已发货页面呈现，视图：已发货
        {
          order_id: '',
          id: '',
          vendor: '',
          item_describe: '',
          voltage: '',
          amount: '',
          price: '',
          shiped_date: '',
          tracking_number: '',
          tracking: '',
          note: '',
          address: '',
        },
      ],
      dialog: false, //订单编辑弹窗
      dialog2: false, //发货弹窗
      dialog_data: {
        //订单编辑弹窗数据
        address: '',
        date: '',
        note: '',
        order_id: '',
        price: '',
        state_id: '',
        tb_id: '',
        vendor_id: '',
        vendor_class: '',
        vendor_list: [],
        state_list: [],
        voltage_list: [],
        product_list: {
          select: { id: '', name: '' },
          items: [],
        },
        express_list: [],
        price_list: [],
        sub_order: [
          {
            amount: '',
            freight_company: '',
            freight_cost: '',
            item_describe: '',
            id: '',
            order_id: '',
            price: '',
            product_id: '38',
            shiped_date: '',
            state: '1',
            tracking_number: '',
            vendor_id: '',
            voltage: '',
            product_list: {
              select: {
                id: '',
                name: '',
              },
              items: [],
            },
          },
        ],
      },
      dialog2_data: {
        //发货弹窗数据
        express: '',
        tracking: '',
        shipped_date: '',
        sub_order_id: '',
        item_describe: '',
        order_id: '',
      },
    }
  },
  computed: {
    aMonthAgo: function () {
      let d = new Date()
      let str = ''
      str += d.getFullYear() + '-'
      str += d.getMonth() + '-'
      str += d.getDate()
      return str
    },
    Now: function () {
      let now = new Date()
      let YYYY = now.getFullYear()
      let MM = ('0' + (now.getMonth() + 1)).toString().slice(-2)
      let DD = ('0' + now.getDate()).toString().slice(-2)
      let hh = ('0' + now.getHours()).toString().slice(-2)
      let mm = ('0' + now.getMinutes()).toString().slice(-2)
      let ss = ('0' + now.getSeconds()).toString().slice(-2)
      // now = YYYY + "-" + MM + "-" + DD + " " + hh + ":" + mm + ":" + ss
      now = YYYY + '-' + MM + '-' + DD
      return now
    },
    NowWithTime: function () {
      let now = new Date()
      let YYYY = now.getFullYear()
      let MM = ('0' + (now.getMonth() + 1)).toString().slice(-2)
      let DD = ('0' + now.getDate()).toString().slice(-2)
      let hh = ('0' + now.getHours()).toString().slice(-2)
      let mm = ('0' + now.getMinutes()).toString().slice(-2)
      let ss = ('0' + now.getSeconds()).toString().slice(-2)
      now = YYYY + '-' + MM + '-' + DD + ' ' + hh + ':' + mm + ':' + ss
      return now
    },
    viewModel: {
      // getter
      get: function () {
        return this.view_model
      },
      // setter
      set: function (newValue) {
        this.view_model = newValue
        this.initView()
      },
    },
  },
  methods: {
    TEST1: function () {},

    TEST2: function () {
      localStorage.setItem('state_list', '')
    },
    TEST3: function () {
      //清除本地Cache
      this.localStorage.removeItem('state_list')
      this.localStorage.removeItem('vendor_list')
      this.localStorage.removeItem('voltage_list')
      this.localStorage.removeItem('product_list')
      this.localStorage.removeItem('express_list')
      this.localStorage.removeItem('price_list')
    },
    updateLogistics: function (d) {
      let param = new URLSearchParams()
      param.append('suborder_id', d.sub_order_id)
      param.append('freight_company', d.express)
      param.append('tracking_number', d.tracking)
      param.append('shiped_date', d.shipped_date)
      param.append('order_id', d.order_id)

      let clipboard = document.getElementById('clipboard')
      clipboard.select()
      document.execCommand('copy')

      if (this.updateLogistics_tbid != '') {
        window.open(
          'https://wuliu.taobao.com/user/consign.htm?trade_id=' +
            this.updateLogistics_tbid
        )
      }

      let me = this
      this.axios.post('updateLogistics', param).then((response) => {
        if (response.data == 'sucess') {
          me.initView()
          me.dialog2 = false //当所有子订单都发货完毕后，主订单修改为已发货
        }
      })
    },

    updateOrder: function () {
      let i = 0
      let re_state = true
      let send_data = JSON.stringify(this.dialog_data) //利用JSON的转换，也可以实现对象深度拷贝
      send_data = JSON.parse(send_data)

      //清空一些无用的list，减轻传输压力
      delete send_data.product_list
      delete send_data.express_list
      delete send_data.state_list
      delete send_data.vendor_list
      delete send_data.voltage_list
      delete send_data.price_list
      for (i in send_data.sub_order) {
        send_data.sub_order[i].product_id =
          send_data.sub_order[i].product_list.select.id
        send_data.sub_order[i].item_describe =
          send_data.sub_order[i].product_list.select.name
        delete send_data.sub_order[i].product_list //记录了当前所选项后，product_list 无用
        if (
          send_data.sub_order[i].state == 1 ||
          send_data.sub_order[i].state == 2 ||
          send_data.sub_order[i].state == 12
        ) {
            re_state = false // re_state 有什么用？ (2021-10-30 标注)
           
        }
      }

      // 2021-10-30 屏蔽：
      /*
      if (re_state && send_data.sub_order.length > 0) {
        send_data.state_id = send_data.sub_order[0].state
      } else {
        send_data.state_id = '2'
      }
     */
      send_data = JSON.stringify(send_data)
      let me = this
      let param = new URLSearchParams()
      param.append('data', send_data)

      console.log(send_data)

      this.axios.post('updateOrder', param).then((response) => {
        me.dialog = false
        me.initView()
      })
    },
    getOrderSubOrders: function (order_id, callback) {
      let param = new URLSearchParams()
      param.append('order_id', order_id)
      this.axios
        .post('getOrderSubOrders',param)
        .then((response) => {
          if (!response.data == '') {
            callback(response.data)
          } else {
            callback('null')
          }
        })
    },
    onOrderClick: function (e) {
      let me = this
      let i, j
      let order_id = e.srcElement.innerText
      this.getOrderSubOrders(order_id, function (ret) {
        //先从主订单(视图变量)里获取现成的数据：
        for (i in me.orders) {
          if (order_id === me.orders[i].id) {
            me.dialog_data.order_id = order_id
            me.dialog_data.date = me.orders[i].date
            me.dialog_data.tb_id = me.orders[i].tb_id
            me.dialog_data.state_id = me.orders[i].state_id
            me.dialog_data.vendor_id = me.orders[i].vendor_id
            me.dialog_data.price = me.orders[i].price
            me.dialog_data.note = me.orders[i].note
            me.dialog_data.address = me.orders[i].address
            break
          }
        }
        //获取子订单数据：
        let order_SubOrders
        ret == 'null'
          ? (order_SubOrders = [])
          : (order_SubOrders = JSON.parse(ret))
        me.dialog_data.sub_order = order_SubOrders //此处为浅拷贝，原有的sub_order结构将被冲刷(冲掉了sub_order独立的product_list)
        for (j in me.dialog_data.sub_order) {
          me.dialog_data.sub_order[j].product_list = {
            //在此处重建每一个sub_order独立的product_list
            select: {
              id: '',
              name: '',
            },
            items: [],
          }
          if (
            me.dialog_data.sub_order[j].shiped_date == '1900-01-01 00:00:00'
          ) {
            me.dialog_data.sub_order[j].shiped_date = ''
          }
          //因为有“自定义”项目的存在，每个子订单的产品列表各不相同，需要修改并分开处理：
          let item_describe = ''
          for (i in me.dialog_data.product_list.items) {
            //如果子订单中产品选择的是'自定义'则显示自定义的值：
            if (
              me.dialog_data.product_list.items[i].id == '38' &&
              me.dialog_data.sub_order[j].product_id == '38'
            ) {
              item_describe = order_SubOrders[j].item_describe
            } else {
              item_describe = me.dialog_data.product_list.items[i].name
            }
            //重建product_list必须逐条push而不能直接整体赋值, 参考资料：JS深浅拷贝
            me.dialog_data.sub_order[j].product_list.items.push({
              id: me.dialog_data.product_list.items[i].id,
              name: item_describe,
            })
          }
          //product_list选中项：
          me.dialog_data.sub_order[j].product_list.select.id =
            me.dialog_data.sub_order[j].product_id
          me.dialog_data.sub_order[j].product_list.select.name =
            me.dialog_data.sub_order[j].item_describe
        }
        i = 0 //手工指定vendor_class
        for (i in me.dialog_data.vendor_list) {
          if (me.dialog_data.vendor_id == me.dialog_data.vendor_list[i].id) {
            me.dialog_data.vendor_class = parseInt(
              me.dialog_data.vendor_list[i].class
            )
            break
          }
        }
        //me.dialog = true 如果通过模型打开dialog，则需要手动阻止事件传播。
        me.$refs.dialog.click()
      })
    },
    onAddNewOrder: function () {
      this.dialog_data.address = ''
      this.dialog_data.date = this.NowWithTime
      this.dialog_data.note = ''
      this.dialog_data.order_id = ''
      this.dialog_data.price = '0'
      this.dialog_data.state_id = '2'
      this.dialog_data.tb_id = ''
      this.dialog_data.vendor_id = '84'
      this.dialog_data.vendor_class = '1'
      this.dialog_data.sub_order = [
        {
          amount: '1',
          freight_company: '17',
          state: '2',
          voltage: '6',
          freight_cost: '',
          id: '',
          item_describe: '',
          id: '',
          order_id: '',
          price: '0',
          product_id: '38',
          shiped_date: '',
          tracking_number: '',
          vendor_id: '84',
          product_list: [],
        },
      ]
      //构建第一个子订单的时候，也必须深拷贝，不能直接赋值
      let temp = JSON.stringify(this.dialog_data.product_list) //利用JSON的转换，也可以实现对象深度拷贝
      temp = JSON.parse(temp)
      this.dialog_data.sub_order[0].product_list = temp
      this.$refs.dialog.click()
    },
    onDelOrder: function () {
      let param = new URLSearchParams()
      param.append('order_id', this.dialog_data.order_id)
      let me = this
      this.axios.post('deleteOrder', param).then((response) => {
        me.dialog = false
        me.initView()
      })
    },
    onDelSubOrder: function (item) {
      let i
      for (i in this.dialog_data.sub_order) {
        if (this.dialog_data.sub_order[i] == item) {
          this.dialog_data.sub_order.splice(i, 1)
        }
      }
    },
    onAddSubOrder: function (item) {
      let temp = JSON.stringify(this.dialog_data.product_list) //利用JSON的转换，也可以实现对象深度拷贝
      temp = JSON.parse(temp)
      item.push({
        product_list: temp,
        product_id: '38',
        amount: '1',
        freight_company: '17',
        state: '2',
        voltage: '6',
        price: '0',
      })
    },
    onProductListBlur: function (e, item) {
      let i
      if (item.product_list.select.id == '38') {
        for (i in item.product_list.items) {
          if (item.product_list.items[i].id == '38') {
            item.product_list.items[i].name = e.target.value
          }
        }
        item.product_list.select.name = e.target.value
      }
    },
    onMouseHover: function (e) {
      this.old_element = e
      this.old_color = e.srcElement.style.color
      e.srcElement.style.color = 'rgb(0, 30, 255)' //文字改为蓝色
      e.srcElement.style.cursor = 'pointer' //改变指针形状
    },
    onMouseLeave: function () {
      this.old_element.srcElement.style.color = this.old_color //还原颜色
    },

    onMouseClick: function (pid, tb_id) {
      let i, j
      for (i in this.orders) {
        for (j in this.orders[i].sub_order) {
          if (pid == this.orders[i].sub_order[j].sub_order_id) {
            this.dialog2_data.order_id = this.orders[i].id
            break
          }
        }
      }

      this.updateLogistics_tbid = tb_id // 点发货后打开的tb_id

      let param = new URLSearchParams()

      param.append('suborder_id', pid)
      let me = this
      this.axios.post('getLogistics', param).then((response) => {
        let data = JSON.parse(response.data)
        me.dialog2_data.item_describe = data[0].item_describe
        me.dialog2_data.sub_order_id = data[0].id
        me.dialog2_data.express = data[0].freight_company
        me.dialog2_data.tracking = data[0].tracking_number
        me.dialog2_data.shipped_date = data[0].shiped_date
        if (data[0].shiped_date == '1900-01-01 00:00:00') {
          this.dialog2_data.shipped_date = this.Now
        }
        this.$refs.dialog2.click()
      })
    },
    onVendorChange: function (e) {
      let i = 0
      let taobao = ''
      let contacts = ''
      let tel = ''
      let address = ''
      let vendor_id = e
      let vendor_class = 0
      //合成地址栏：
      for (i in this.dialog_data.vendor_list) {
        if (vendor_id == this.dialog_data.vendor_list[i].id) {
          !this.dialog_data.vendor_list[i].taobao == ''
            ? (taobao = '{' + this.dialog_data.vendor_list[i].taobao + '} ')
            : (taobao = '')
          !this.dialog_data.vendor_list[i].contacts == ''
            ? (contacts = this.dialog_data.vendor_list[i].contacts + ', ')
            : (contacts = '')
          !this.dialog_data.vendor_list[i].tel == ''
            ? (tel = this.dialog_data.vendor_list[i].tel + ', ')
            : (tel = '')
          !this.dialog_data.vendor_list[i].address == ''
            ? (address = this.dialog_data.vendor_list[i].address)
            : (address = '')
          this.dialog_data.address = taobao + contacts + tel + address
          this.dialog_data.vendor_class = parseInt(
            this.dialog_data.vendor_list[i].class
          )
          break
        }
      }
      this.updatePrice()
    },
    color: function (state_id) {
      switch (state_id) {
        case '1':
          return '#ff5252' //未：red
          break
        case '2':
          return '#1B5E20' //待：green
          break
        case '3':
          return '#202020' //已：black
          break
        case '4':
          return '#fb8c00' //包：orange
          break
        case '12':
          return '#e0e0e0' //关：grey
          break
        default:
          return '#202020' //black
      }
    },
    checkListCache: function () {
      let lists = [
        'state_list',
        'vendor_list',
        'voltage_list',
        'product_list',
        'express_list',
        'price_list',
      ]
      let re_build = false
      let i = 0
      for (i = 1; i <= 50; i++) {
        //折扣计算器创建1-50个元素
        this.discount_pcs_list.push({ name: i, id: i })
      }
      i = 0
      //localStorage里不存在对应的list或list未定义，都需要重建：
      for (i in lists) {
        if (!(lists[i] in this.localStorage)) {
          re_build = true
        } else {
          if (this.localStorage[lists[i]].length == 0) {
            re_build = true
          }
        }
      }
      if (re_build) {
        i = 0
        let me = this
        this.axios
          .get('getDefaultList', {
            params: {},
          })
          .then((response) => {
            for (i in lists) {
              this.localStorage.setItem(lists[i], response.data[i])
            }
            this.restoreList() //从远程恢复List
          })
      } else {
        this.restoreList() //直接从localStorage恢复List
      }
    },
    restoreList: function () {
      //将localStorage保存的list恢复到视图：
      this.dialog_data.state_list = JSON.parse(
        this.localStorage.getItem('state_list')
      )
      this.dialog_data.vendor_list = JSON.parse(
        this.localStorage.getItem('vendor_list')
      )
      this.dialog_data.voltage_list = JSON.parse(
        this.localStorage.getItem('voltage_list')
      )

      //product_list的恢复略有不同：
      let plist = JSON.parse(this.localStorage.getItem('product_list'))
      for (let i in plist) {
        this.dialog_data.product_list.items.push({
          id: plist[i].id,
          name: plist[i].name,
        })
      }
      this.dialog_data.product_list.select.id = '38'
      this.dialog_data.product_list.select.name = '自定义'

      this.dialog_data.express_list = JSON.parse(
        this.localStorage.getItem('express_list')
      )
      this.dialog_data.price_list = JSON.parse(
        this.localStorage.getItem('price_list')
      )
    },
    updateDiscountPrice: function () {
      let j = 0
      let retail_price = 0
      let wholesale_price = 0
      //计算价格：
      for (j in this.dialog_data.price_list) {
        if (this.discount_product.id == this.dialog_data.price_list[j].pi) {
          if (this.dialog_data.price_list[j].vi == '1') {
            retail_price =
              parseInt(this.dialog_data.price_list[j].pr) *
              parseInt(this.discount_pcs.id)
          }
          if (this.dialog_data.price_list[j].vi == '2') {
            wholesale_price =
              parseInt(this.dialog_data.price_list[j].pr) *
              parseInt(this.discount_pcs.id)
          }
        }
      }
      this.discount = wholesale_price - retail_price
    },
    updatePrice: function () {
      let i = 0
      let j = 0
      let total_price = 0
      //计算价格：
      for (i in this.dialog_data.sub_order) {
        for (j in this.dialog_data.price_list) {
          if (
            this.dialog_data.sub_order[i].product_list.select.id ==
            this.dialog_data.price_list[j].pi
          ) {
            if (
              this.dialog_data.vendor_class == this.dialog_data.price_list[j].vi
            ) {
              if (
                parseInt(this.dialog_data.price_list[j].pr) > 0 &&
                parseInt(this.dialog_data.sub_order[i].amount) > 0
              ) {
                this.dialog_data.sub_order[i].price =
                  parseInt(this.dialog_data.price_list[j].pr) *
                  parseInt(this.dialog_data.sub_order[i].amount)
                this.dialog_data.sub_order[i].price = String(
                  this.dialog_data.sub_order[i].price
                )
              } else {
                this.dialog_data.sub_order[i].price = '0'
              }
            }
          }
        }
        total_price =
          total_price + parseInt(this.dialog_data.sub_order[i].price)
      }
      this.dialog_data.price = String(total_price)
    },
    updateTotalPrice: function () {
      //手工输入单价时，不与自动计算模块联动，只更新总价
      let i = 0
      let total_price = 0
      for (i in this.dialog_data.sub_order) {
        total_price =
          total_price + parseInt(this.dialog_data.sub_order[i].price)
      }
      this.dialog_data.price = String(total_price)
    },
    range_default: function () {
      this.date_start = this.aMonthAgo
      this.date_end = this.Now
      this.keyword = ''
    },
    range_whole: function () {
      this.date_start = '2017-01-01'
      this.date_end = this.Now
    },
    initView: function () {
      let me = this
      let param = new URLSearchParams()
      switch (this.view_model) {
        case 0:
          param.append('keyword', me.keyword)
          param.append('date_start', me.date_start + ' 00:00:00')
          param.append('date_end', me.date_end + ' 23:59:59')
          this.axios
            .post('getOrders', param)
            .then((response) => {
              if (response.data == '') {
                me.orders = []
                return
              }
              me.orders = JSON.parse(response.data)
              let i, j, v, pl, p, c

              for (i in me.orders) {
                v = me.orders[i].vendor.split('|')
                me.orders[i].vendor_id = v[0]
                me.orders[i].vendor = v[1]

                me.orders[i].sub_order = []
                pl = me.orders[i].p_info.split(',')
                for (j in pl) {
                  p = pl[j].split('|')
                  if (
                    me.orders[i].state_id == '1' ||
                    me.orders[i].state_id == '12'
                  ) {
                    //两个特殊情况：当主订单状态为1(未)或12(关)的时候，产品列表沿用主订单状态的颜色
                    c = me.color(me.orders[i].state_id)
                  } else {
                    c = me.color(p[4])
                  }
                  me.orders[i].sub_order.push({
                    sub_order_id: p[0],
                    item_describe: p[1],
                    amount: p[2] + '\n',
                    voltage: p[3] + '\n',
                    color: c,
                    state: p[5] + '\n',
                    price: p[6] + '\n',
                    order_id: me.orders[i].order_id,
                  })
                }
                me.orders[i].color = me.color(me.orders[i].state_id)
                delete me.orders[i].p_info //这几个属性解析之后就没用了，可以删除
              }
            })
          break
        case 1:
          param.append('keyword', me.keyword)
          param.append('date_start', me.date_start + ' 00:00:00')
          param.append('date_end', me.date_end + ' 23:59:59')
          this.axios
            .post('getShippedSubOrders',param)
            .then((response) => {
              if (response.data == '') {
                me.shipped_orders = []
                return
              }
              me.shipped_orders = JSON.parse(response.data)
            })

          break
        default:
      }

      this.checkListCache() //检查list本地缓存
    },
  },
  mounted: function () {
    this.date_start = this.aMonthAgo
    this.date_end = this.Now
    this.viewModel = 0
  },
}
</script>
<style>
.box {
  display: flex;
}

/* 表格专属CSS来自DivTable.com */
.divTable {
  display: table;
  width: 100%;
}

.divTableRow {
  display: table-row;
}

.divTableRow:hover {
  background-color: #eeeeee;
}

.divTableHeading {
  background-color: #eee;
  display: table-header-group;
}

.divTableCell,
.divTableHead {
  border: 1px solid #e0e0e0;
  display: table-cell;
  padding: 3px 10px;
  vertical-align: middle;
}

.center {
  text-align: center;
}

.nowrap {
  white-space: nowrap;
}

.divTableHeading {
  background-color: #eee;
  display: table-header-group;
}

.divTableFoot {
  background-color: #eee;
  display: table-footer-group;
  font-weight: bold;
}

.divTableBody {
  display: table-row-group;
}
</style>
