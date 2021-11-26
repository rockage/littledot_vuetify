<template>
  <div
    class="text-caption"
    style="display:flex;flex-direction: row;width:100%;height: calc(91vh);overflow:hidden;"
  >
    <div
      style="display:flex;
              flex-direction: column;
              justify-content:flex-start;  
              align-items: center;
              background-color: #95a5a6;
              width:10%;height: 100%;"
    >
      <v-card
        style="width: 80%;margin-top: 10px;height: 90px;padding-left:10px;"
      >
        <v-radio-group mandatory v-model="server">
          <v-radio label="VPS" @click="ws_init('vps')" value="vps"></v-radio>
          <v-radio
            label="Local"
            @click="ws_init('localhost')"
            value="localhost"
          ></v-radio>
        </v-radio-group>
      </v-card>

      <v-btn
        dark
        color="teal"
        @click="ws.send('#rockage_cmds#status#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        状态
      </v-btn>

      <v-btn
        color="primary"
        @click="ws.send('#rockage_cmds#restart#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        重启
      </v-btn>
      <v-btn
        color="primary"
        @click="ws.send('#rockage_cmds#check#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        登录
      </v-btn>
      <v-btn
        color="primary"
        @click="ws.send('#rockage_cmds#start#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        抓单
      </v-btn>
      <v-btn
        color="secondary"
        @click="ws.send('#rockage_cmds#analyse#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        分析
      </v-btn>
      <v-btn
        color="secondary"
        @click="ws.send('#rockage_cmds#save#rockage_datas#null')"
        style="width: 80%;margin-top: 10px;"
      >
        存盘
      </v-btn>
      <v-btn
        dark
        color="indigo"
        @click="ws.send('#rockage_cmds#automatic#rockage_datas#open')"
        style="width: 80%;margin-top: 10px;"
      >
        自动
      </v-btn>

      <v-btn
        dark
        color="indigo"
        @click="ws.send('#rockage_cmds#automatic#rockage_datas#close')"
        style="width: 80%;margin-top: 10px;"
      >
        停止
      </v-btn>
      <v-btn @click="cls" style="width: 80%;margin-top: 10px;">
        CLS
      </v-btn>

      <input
        v-model="smscode"
        style="width: 80%;height:30px;margin-top: 50px;text-align:center"
        placeholder="短信验证码"
      />

      <v-btn @click="sendsms" style="width: 80%;margin-top: 5px;">
        发送!
      </v-btn>
    </div>
    <div
      style="background-color: #bdc3c7;width:45%;max-height: height: 500px%;overflow:scroll;"
      ref="msg"
    ></div>
    <div
      style="background-color: #ecf0f1;width:45%;max-height: height: 500px%;overflow:scroll;"
      ref="orders_msg"
    ></div>
  </div>
</template>
<script>
export default {
  name: 'crawler',
  data() {
    return {
      ws: '',
      smscode: '',
      page_orders: '',
      server: 'vps',
      //server: 'localhost',
    }
  },
  computed: {},

  methods: {
    cls: function() {
      this.$refs.msg.innerHTML = ''
      this.$refs.orders_msg.innerHTML = ''
    },

    sendsms: function() {
      this.ws.send('#rockage_cmds#sms#rockage_datas#' + this.smscode)
    },

    orders_msg_add: function(orders_msg) {
      let div = this.$refs.orders_msg
      setTimeout(() => {
        div.innerHTML += `<div>${orders_msg}</div>`
        div.scrollTop = div.scrollHeight // 自动翻滚到最后
      }, 0)
    },

    msg_add: function(msg) {
      let div = this.$refs.msg
      let myDate = new Date()
      let dt =
        myDate.getHours() +
        ':' +
        myDate.getMinutes() +
        ':' +
        myDate.getSeconds() +
        ':' +
        myDate.getMilliseconds() +
        '>'
      //此时必须异步执行滚动条滑动至底部
      setTimeout(() => {
        div.innerHTML += `<div>${dt} ${msg}</div>`
        div.scrollTop = div.scrollHeight // 自动翻滚到最后
      }, 0)
    },

    ws_onopen: function($event) {
      console.log('websocks connected.')
    },

    ws_onmessage: function($event) {
      let e = $event
      switch (e.data) {
        case 'cls':
          this.cls
          break
        default:
          let s = ''
          let regex = /#ORDERS#(.*)/g
          let m = regex.exec(e.data)
          if (m) {
            this.orders_msg_add(m[1])
          } else {
            this.msg_add(e.data)
          }
      }
    },
    ws_init: function(str) {
      switch (str) {
        case 'vps':
          this.ws = new WebSocket('wss://rockage.net/websocket')
          break
        case 'localhost':
          this.ws = new WebSocket('ws://127.0.0.1:6060')
          break
      }
      this.ws.onopen = this.ws_onopen
      this.ws.onmessage = this.ws_onmessage
      this.cls
    },
    send: function(msg) {
      this.ws.send(msg)
    },
  },
  mounted: function() {
    setTimeout(() => {
      this.ws_init(this.server)
    }, 500) // 延时500ms初始化
  },
}
</script>
<style></style>
