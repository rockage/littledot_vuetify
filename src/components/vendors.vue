<template>
  <!--主界面-->
  <div class="text-caption">
    <div
      class="flex-row"
      style="align-items: flex-start; justify-content: flex-start"
    >
      <v-btn dense @click="onAddClick()" class="pa-1 ma-2" small color="teal">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
      <v-btn dense @click="sortPosition" class="pa-1 ma-2" small>SORT</v-btn>
      <v-btn dense @click="updatePosition" class="pa-1 ma-2" small>APPLY</v-btn>
    </div>
    <div class="flex-column" style="align-items: center; margin-right: 20px">
      <div class="flex-column" style="width: 100%">
        <v-card width="100%" v-for="vendor in vendors" class="pa-1 ma-3">
          <div class="flex-row" style="justify-content: space-between">
            <div class="flex-row" style="width: 100%">
              <div
                class="flex-column"
                style="justify-content: center; font-size: medium"
              >
                <div
                  contenteditable="true"
                  style="
                    width: 30px;
                    height: 22px;
                    background-color: #ecf0f1;
                    text-align: center;
                    border: 0px;
                    ime-mode: Disabled;
                  "
                  ref="position"
                >
                  {{ vendor.position }}
                </div>
              </div>
              <div class="flex-column" style="margin-left: 10px; width: 100%">
                <div class="flex-row">
                  <div style="font-size: large; width: 100%" ref="nickname">
                    {{ vendor.nickname }}
                  </div>
                  <div ref="id">{{ vendor.ID }}</div>
                </div>
                <div style="color: #7f8c8d; width: 100%">
                  <div style="font-size: medium; width: 100%" ref="sale">
                    {{ vendor.sale }}
                  </div>
                </div>
              </div>
            </div>
            <div class="flex-column" style="justify-content: center">
              <v-btn icon @click="vendor.show = !vendor.show">
                <v-icon>{{
                  vendor.show ? 'mdi-chevron-up' : 'mdi-chevron-down'
                }}</v-icon>
              </v-btn>
            </div>
          </div>
          <div v-show="vendor.show" class="pa-0 ma-0">
            <v-divider></v-divider>
            <div class="flex-row">
              <div class="item-title">分类</div>
              <div ref="vendor_class" class="item-text">
                {{ vendor.class.name }}
              </div>
            </div>
            <div class="flex-row">
              <div class="item-title">全称</div>
              <div ref="fullname" class="item-text">
                {{ vendor.fullname }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">联系人</div>
              <div ref="contacts" class="item-text">
                {{ vendor.contacts }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">电话</div>
              <div ref="tel" class="item-text">
                {{ vendor.tel }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">地址</div>
              <div ref="address" class="item-text">
                {{ vendor.address }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">Email</div>
              <div ref="email" class="item-text">
                {{ vendor.email }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">微信</div>
              <div ref="wechat" class="item-text">
                {{ vendor.wechat }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">旺旺</div>
              <div ref="taobao" class="item-text">
                {{ vendor.taobao }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">QQ</div>
              <div ref="qq" class="item-text">
                {{ vendor.qq }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">支付宝</div>
              <div ref="alipay" class="item-text">
                {{ vendor.alipay }}
              </div>
            </div>
            <div style="display: flex; flex-direction: row">
              <div class="item-title">银行账号</div>
              <div ref="bank" class="item-text">
                {{ vendor.bank }}
              </div>
            </div>
            <div class="flex-row" style="justify-content: center">
              <v-btn
                color="green lighten-2"
                style="width: 100%"
                text
                @click="onEditClick(vendor.index)"
              >
                EDIT
              </v-btn>
            </div>
          </div>
        </v-card>
      </div>
    </div>
    <div>
      <v-snackbar v-model="alert" :timeout="3000">
        产品位置已按当前视图保存
        <template v-slot:action="{ attrs }">
          <v-btn color="blue" text v-bind="attrs" @click="alert = false">
            Close
          </v-btn>
        </template>
      </v-snackbar>
    </div>
    <!-- 编辑弹窗 -->
    <v-dialog v-model="dialog_show" max-width="90%">
      <v-card>
        <div class="flex-column" style="padding: 20px">
          <div class="flex-row" style="flex-wrap: wrap; width: 100%">
            <v-text-field
              v-model="dialog_data.nickname"
              label="简称"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>

            <v-select
              label="分类"
              v-model="dialog_data.class"
              :items="class_list.items"
              item-text="name"
              item-value="id"
              return-object
              style="flex-grow: 0; margin-right: 10px"
            ></v-select>

            <v-text-field
              v-model="dialog_data.fullname"
              label="全称"
              style="flex-grow: 1; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row">
            <v-text-field
              v-model="dialog_data.sale"
              label="客户描述"
              style="flex-grow: 1; width: 40%; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row">
            <v-text-field
              v-model="dialog_data.address"
              label="地址"
              style="flex-grow: 1; width: 40%; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row" style="flex-wrap: wrap; width: 100%">
            <v-text-field
              v-model="dialog_data.contacts"
              label="联系人"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>
            <v-text-field
              v-model="dialog_data.tel"
              label="电话"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>
            <v-text-field
              v-model="dialog_data.email"
              label="Email"
              style="flex-grow: 1; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row" style="flex-wrap: wrap; width: 100%">
            <v-text-field
              v-model="dialog_data.taobao"
              label="淘宝ID"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>

            <v-text-field
              v-model="dialog_data.alipay"
              label="支付宝"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>

            <v-text-field
              v-model="dialog_data.bank"
              label="银行账号"
              style="flex-grow: 1; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row" style="flex-wrap: wrap; width: 100%">
            <v-text-field
              v-model="dialog_data.wechat"
              label="微信"
              style="flex-grow: 0; margin-right: 10px"
            ></v-text-field>
            <v-text-field
              v-model="dialog_data.qq"
              label="QQ"
              style="flex-grow: 1; margin-right: 10px"
            ></v-text-field>
          </div>

          <div class="flex-row" style="justify-content: space-between">
            <div>
              <v-btn color="pink" text v-show="del_button"  @click="onDeleteClick()"> DELETE </v-btn>
            </div>
            <div class="flex-row">
              <v-btn color="blue darken-1" text @click="dialog_show = false">
                CLOSE
              </v-btn>

              <v-btn color="success" text @click="updateVendor()"> SAVE </v-btn>
            </div>
          </div>
        </div>
      </v-card>
    </v-dialog>
    <!-- MessageBox 弹窗 -->
    <v-dialog v-model="msgbox_show" width="50%">
      <v-card
        class="flex-column"
        style="
          height: 200px;
          justify-content: space-between;
          align-items: center;
        "
      >
        <div></div>
        <!-- 占位DIV -->
        <div>{{ msgbox_text }}</div>
        <div class="flex-row" style="justify-content: flex-end">
          <v-btn color="blue darken-1" text @click="msgbox_show = false">
            CLOSE
          </v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
<script>
export default {
  name: 'products',
  data() {
    return {
      vendors: new Array(),
      dialog_show: '',
      dialog_data: {},
      class_list: {
        items: [
          {
            id: '0',
            name: '供应商',
          },
          {
            id: '1',
            name: '零售商',
          },
          {
            id: '2',
            name: '批发商',
          },
        ],
      },
      msgbox_text: '',
      msgbox_show: false,

      del_button: false,
      alert: '',
      show: false,
    }
  },
  methods: {
    test2: function () {
      this.vendors = []
      this.initView()
    },
    onAddClick: function () {
      this.op = 'addnew'
      this.dialog_data = {
        // 清空dialog视图
        class: {
          id: '0',
          name: '供应商',
        },
      }

      this.del_button = false
      this.dialog_show = true
    },
    onEditClick: function (index) {
      this.op = 'edit'
      let vendor_class = this.$refs.vendor_class[index].innerText
      for (let i in this.class_list.items) {
        if (vendor_class == this.class_list.items[i].name) {
          // name 转 id
          vendor_class = {
            id: this.class_list.items[i].id,
            name: this.class_list.items[i].name,
          }
          break
        }
      }
      this.dialog_data = {
        // 更新dialog视图
        id: this.$refs.id[index].innerText,
        nickname: this.$refs.nickname[index].innerText,
        sale: this.$refs.sale[index].innerText,
        fullname: this.$refs.fullname[index].innerText,
        contacts: this.$refs.contacts[index].innerText,
        tel: this.$refs.tel[index].innerText,
        address: this.$refs.address[index].innerText,
        email: this.$refs.email[index].innerText,
        wechat: this.$refs.wechat[index].innerText,
        taobao: this.$refs.taobao[index].innerText,
        qq: this.$refs.qq[index].innerText,
        alipay: this.$refs.alipay[index].innerText,
        bank: this.$refs.bank[index].innerText,
        class: vendor_class,
      }
      this.del_button = true
      this.dialog_show = true
    },
    onDeleteClick: function () {
      let param = new URLSearchParams()
      param.append('id', this.dialog_data.id)
      param.append('op', "delete")
      let me = this
      this.axios.post('updateVendors', param).then((response) => {
        if (response.data == 'sucess') {
          me.initView()
          me.dialog_show = false
        }
      })
    },
    updateVendor: function () {
      let isEmpty = (exp) => {
        if (exp === undefined) {
          console.log('is undefined')
          return true
        }
        if (exp === null) {
          console.log('is null')
          return true
        }
        if (!exp) {
          console.log('is !exp')
          return true
        }
        if (exp.length === 0) {
          console.log('length === 0')
          return true
        }
        if (exp === '') {
          console.log("exp === ''")
          return true
        }
        return false
      }

      if (isEmpty(this.dialog_data.nickname) == true) {
        this.msgbox_text = '记录保存失败，请至少指定 “简称” 和 “分类” 。 '
        this.msgbox_show = true
        return
      }
      let param = new URLSearchParams()
      param.append('datas', JSON.stringify(this.dialog_data))
      param.append('op', this.op)
      console.log(JSON.stringify(this.dialog_data))
      let me = this
      this.axios.post('updateVendors', param).then((response) => {
        if (response.data == 'sucess') {
          me.initView()
          me.dialog_show = false
        }
      })
    },
    sortPosition() {
      let i = 0
      for (i in this.vendors) {
        this.$refs.position[i].innerText = parseInt(i) + 1
      }
    },
    updatePosition() {
      let i = 0
      let vendor_position = new Array()
      for (i in this.vendors) {
        vendor_position.push({
          id: this.vendors[i].ID,
          position: this.$refs.position[i].innerText,
        })
      }
      let me = this
      let param = new URLSearchParams()
      param.append('newPosition', JSON.stringify(vendor_position))
      param.append('op', 'updatePosition')
      this.axios.post('updateVendors', param).then((response) => {
        if (response.data == 'sucess') {
          me.alert = !me.alert
          this.initView()
        }
      })
    },
    initView: function () {
      let me = this
      this.axios
        .get('getVendors', {
          params: {},
        })
        .then((response) => {
          if (response.data == '') {
            alert('ERROR')
          }
          let datas = response.data

          for (let i in datas) {
            datas[i].index = i
            datas[i].show = false
            for (let j in me.class_list.items) {
              if (datas[i].class == me.class_list.items[j].id) {
                // datas[i].class = me.class_list.items[j].name // id 转 name
                datas[i].class = {
                  name: me.class_list.items[j].name,
                  id: me.class_list.items[j].id,
                }
                break
              }
            }
          }
          me.vendors = []
          this.$nextTick(() => {
            // 将回调延迟到下次 DOM 更新循环之后执行，否则视图无法正常更新
            me.vendors = [...datas] // 通过3个点的方式实现深拷贝（仅适用于数组）
          })
        })
    },
  },
  mounted: function () {
    this.initView()
  },
}
</script>
<style>
.flex-column {
  display: flex;
  flex-direction: column;
}

.flex-row {
  display: flex;
  flex-direction: row;
}

.item-title {
  height: 20px;
  width: 80px;
  background: #ecf0f1;
  text-align: center;
}

.item-text {
  margin-left: 5px;
  width: 100%;
}

.alert {
  position: fixed;
  left: 30%;
  top: 30%;
  z-index: 11;
}
</style>