<template>
  <div class="text-caption">
    <div style=" display: flex;  flex-direction: column; width:100% ">
      <div
        style=" display: flex;  flex-direction: row; align-items:center; margin-left: 10px;"
      >
        <v-menu offset-y style="max-width: 100px">
          <template v-slot:activator="{ attrs, on }">
            <v-btn
              small
              color="primary"
              class="white--text ma-1"
              v-bind="attrs"
              v-on="on"
            >
              {{ s_year }}
            </v-btn>
          </template>
          <v-list>
            <v-list-item link v-for="start in s_years">
              <v-list-item-title
                @click="chartsYearChange('set_start', start)"
                >{{ start }}</v-list-item-title
              >
            </v-list-item>
          </v-list>
        </v-menu>

        <div><v-icon>mdi-chevron-right</v-icon></div>

        <v-menu offset-y style="max-width: 100px">
          <template v-slot:activator="{ attrs, on }">
            <v-btn
              small
              color="primary"
              class="white--text ma-1"
              v-bind="attrs"
              v-on="on"
            >
              {{ e_year }}
            </v-btn>
          </template>
          <v-list>
            <v-list-item link v-for="end in e_years">
              <v-list-item-title @click="chartsYearChange('set_end', end)">{{
                end
              }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
      <div style="margin-top:20px;margin-bottom:20px">
        <highcharts :options="charts"></highcharts>
      </div>
      <hr />
      <div
        style=" width:100%;  display: flex;  flex-direction: column "
        v-for="bar in bars"
      >
        <div
          style=" width:100%;margin-top: 20px;font-size:medium; margin-left:10px;"
        >
          {{ bar.title }} 销售额：{{ bar.total }}（元）{{ bar.waitsend }}
        </div>
        <div style=" width:100%;  display: flex;  flex-direction: row ">
          <div style=" width:50%;">
            <highcharts :options="bar.left_bar"></highcharts>
          </div>
          <div style=" width:50%;">
            <highcharts :options="bar.right_bar"></highcharts>
          </div>
        </div>
        <hr />
      </div>
    </div>
  </div>
</template>
<script>
import { Chart } from 'highcharts-vue'

export default {
  name: 'charts',
  data() {
    return {
      chinese_months: [
        '1月',
        '2月',
        '3月',
        '4月',
        '5月',
        '6月',
        '7月',
        '8月',
        '9月',
        '10月',
        '11月',
        '12月',
      ],
      s_year: '',
      e_year: '',
      s_years: [],
      e_years: [],
      charts: {},
      bars: new Array(),
    }
  },

  components: {
    highcharts: Chart,
  },

  computed: {
    Now: function() {
      let now = new Date()
      return this.timeFormat(now)
    },
  },

  methods: {
    timeFormat: function(date) {
      let YYYY = date.getFullYear()
      let MM = ('0' + (date.getMonth() + 1)).toString().slice(-2)
      let DD = ('0' + date.getDate()).toString().slice(-2)
      let hh = ('0' + date.getHours()).toString().slice(-2)
      let mm = ('0' + date.getMinutes()).toString().slice(-2)
      let ss = ('0' + date.getSeconds()).toString().slice(-2)
      date = YYYY + '-' + MM + '-' + DD + ' ' + hh + ':' + mm + ':' + ss
      return date
    },

    DrawBar: function(datas) {
      let i = 0
      let bars = new Array()
      this.bars = [] // 先清空，注意 this.bars 相当于vue的全局变量,与局部变量 bars 是两个不同的变量
      for (i in datas) {
        bars.push({
          title: datas[i].Title,
          total: datas[i].Total,
          waitsend: datas[i].WaitSend,
          left_bar: {
            chart: {
              type: 'bar',
            },
            title: {
              text: '',
            },
            subtitle: {
              text: '经销商TOP 15',
            },
            xAxis: {
              categories: datas[i].Bars_L.name,
              title: {
                text: null,
              },
            },
            yAxis: {
              min: 0,
              title: {
                text: '',
              },
            },
            tooltip: {
              valueSuffix: '元',
            },
            plotOptions: {
              bar: {
                dataLabels: {
                  enabled: true,
                },
              },
            },
            legend: {
              enabled: false,
            },
            series: [
              {
                name: '销售额',
                data: datas[i].Bars_L.data,
              },
            ],
          },

          right_bar: {
            chart: {
              type: 'bar',
            },
            title: {
              text: '',
            },
            subtitle: {
              text: '产品TOP 15',
            },
            xAxis: {
              categories: datas[i].Bars_R.name,
              title: {
                text: null,
              },
            },
            yAxis: {
              enabled: false,
              min: 0,
              title: {
                text: '',
              },
            },
            tooltip: {
              valueSuffix: '元',
            },
            plotOptions: {
              bar: {
                dataLabels: {
                  enabled: true,
                },
              },
            },
            legend: {
              enabled: false,
            },

            series: [
              {
                name: '销售额',
                data: datas[i].Bars_R.data,
              },
            ],
          }, // -> right_bar
        }) // -> this.bars.push()
      } // -> for 循环
      this.bars = [...bars] // 深度拷贝数组
    },

    DrawChart: function(series, categories) {
      let me = this
      this.charts = {
        chart: {
          type: 'spline',
          inverted: false,
        },
        title: {
          text: '',
        },
        subtitle: {
          text: '', //subtitle
        },
        xAxis: {
          categories: categories,
        },
        yAxis: {
          title: {
            text: '销售额 (K)',
          },
        },
        legend: {
          layout: 'vertical',
          align: 'right',
          verticalAlign: 'middle',
        },
        plotOptions: {
          line: {
            dataLabels: {
              enabled: false,
            },
            enableMouseTracking: true,
          },
          series: {
            cursor: 'pointer',
            events: {
              click: function(event) {
                let s1 = event.point.series.name
                let s2 = event.point.category
                s2 = s2.replace(/月/, '') // 范例s1: “2021” s2: “9月”
                let start = s1 + '-' + s2 + '-01 00:00:00'
                let temp = new Date(start) // 字符串转时间
                start = me.timeFormat(temp) // 格式化为：YYYY-MM-DD hh:mm:ss
                temp = new Date(s1, s2, 0).getDate() // 得到当月的最后一天
                let end = s1 + '-' + s2 + '-' + temp + ' 23:59:59'
                temp = new Date(end)
                end = me.timeFormat(temp)
                me.barMonthChange(start, end)
              },
            },
          },
        }, // =>plotOptions
        series: series,
        responsive: {
          rules: [
            {
              condition: {
                maxWidth: 500,
              },
              chartOptions: {
                legend: {
                  layout: 'horizontal',
                  align: 'center',
                  verticalAlign: 'bottom',
                },
              },
            },
          ],
        }, // =>responsive
      } // =>Highcharts.chart
    }, // =>charts
    barMonthChange: function(s1, s2) {
      // 绘制当月柱状视图

      let start = s1
      let end = s2

      let me = this
      let param = new URLSearchParams()
      param.append('op', 'bars_single_month')
      param.append('start', start)
      param.append('end', end)
      this.axios.post('getStatistics', param).then((response) => {
        let datas = response.data
        let temp = new Date(datas.Title)
        let YYYY = temp.getFullYear()
        let MM = ('0' + (temp.getMonth() + 1)).toString().slice(-2)
        datas.Title = YYYY + '年' + MM + '月'

        let now = new Date() // 如点击本月，额外统计“等待发货的宝贝”
        if (
          temp.getFullYear() == now.getFullYear() &&
          temp.getMonth() == now.getMonth()
        ) {
          datas.WaitSend = ' 未发货的订单合计：' + datas.WaitSend + '（元）'
        } else {
          datas.WaitSend = ''
        }
        let data_arr = new Array() // 将对象变成一个下标为0的对象数组，以兼容DrawBar()函数
        data_arr.push(datas)
        me.DrawBar(data_arr)
      })
    },

    chartsYearChange: function(op, arg) {
      switch (op) {
        case 'set_start':
          this.s_year = arg
          break
        case 'set_end':
          this.e_year = arg
          break
      }
      let s = parseInt(this.s_year)
      let e = parseInt(this.e_year)
      let temp
      if (s > e) {
        // 如：起始时间 > 结束时间，相互交换，防止报错
        temp = s
        s = e
        e = temp
      }

      let s_YYYY = String(s)
      let e_YYYY = String(e)
      let start = s_YYYY + '-01-01 00:00:00'
      let end = e_YYYY + '-12-31 23:59:59'
      temp = new Date(start)
      start = this.timeFormat(temp)
      temp = new Date(end)
      end = this.timeFormat(temp)

      this.makeView(start, end, 0)
    },

    makeView: function(start, end, op) {
      let me = this
      let param = new URLSearchParams()

      param.append('op', 'charts')
      param.append('start', start)
      param.append('end', end)

      this.axios.post('getStatistics', param).then((response) => {
        // 绘制曲线：
          console.log( "1" )
        me.DrawChart(response.data, me.chinese_months)
console.log( "2" )

        if (start.substring(0, 7) == end.substring(0, 7) || op == 1) {
          // 点击到本月方块 或 op=1 (首次加载)
          // 首次加载的特殊性：
          // 1. 曲线是上一年+本年度
          // 2. 柱状图又需要是本月度的，因此这里需要重新对start进行赋值
console.log( "3" )          
          let temp = new Date() 
          let YYYY = temp.getFullYear()
          let MM = ('0' + (temp.getMonth() + 1)).toString().slice(-2)
          let start = YYYY + '-' + MM + '-01 00:00:00' 
          temp = new Date(start) 
          start = this.timeFormat(temp) 
console.log( "4" )          
          this.barMonthChange(start, end)
console.log( "5" )  
console.log( start,"   ",end )        
        } else {
          // 绘制柱状图：
          let param = new URLSearchParams() // 清空参数
          param.append('op', 'bars_years')
          param.append('start', start)
          param.append('end', me.Now)
          me.axios.post('getStatistics', param).then((response) => {
          
            let datas = response.data
            let temp = new Array() // 柱状图逆序
            let n = datas.length
            let i = 0
            for (i in datas) {
              temp.push(datas[n - i - 1])
            }
            me.DrawBar(temp)
          })
        } // -> start.substring(0, 7) == end.substring(0, 7)
      }) // -> axios.post()
    },

    initView: function() {
      // 处理开始时间：
      let temp = new Date() // 获取当前时间
      let YYYY = temp.getFullYear()
      YYYY = String(parseInt(YYYY) - 1) // 年份减1
      this.s_year = YYYY // 视图的起始年份
      let start = YYYY + '-01-01 00:00:00' // 拼装字符串
      temp = new Date(start) // 将字符串又转为时间
      start = this.timeFormat(temp) // 格式化为：YYYY-MM-DD hh:mm:ss
      // 处理结束时间：
      temp = new Date() // 获取当前时间
      YYYY = temp.getFullYear()
      this.e_year = YYYY // 视图的结束年份
      let end = YYYY + '-12-31 23:59:59'

      for (let i = 2017; i <= YYYY; i++) {
        // 生成列表：从2017至今的年份
        this.s_years.push(i)
        this.e_years.push(i)
      }
      this.makeView(start, end, 1) // 第三个参数置1，表示首次加载，需要绘制本月销售柱状图
    },
  },

  mounted: function() {
    this.initView()
  },
}
</script>
<style></style>
