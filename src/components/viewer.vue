<template>
    <div class="text-caption">
        <!-- 订单视图 -->
        <div class="divTable">
            <div class="divTableHeading">
 
                <div class="divTableCell center nowrap">日期</div>
                <div class="divTableCell center nowrap">商家</div>
                <div class="divTableCell center nowrap">产品列表</div>
                <div class="divTableCell center nowrap">数量</div>
                <div class="divTableCell center nowrap">状态</div>
                <div class="divTableCell center nowrap">备注</div>
                <div class="divTableCell center nowrap">地址</div>
            </div>
            <div class="divTableBody">
                <div class="divTableRow" v-for="item in orders" :style="{ color: item.color }">

                    <div class="divTableCell">{{ item.date }} {{ item.tb_id }}</div>
                    <div class="divTableCell nowrap">
                        {{ item.vendor }}
                    </div>
                    <div class="divTableCell nowrap" style="text-align: right">
                        <span v-for="p in item.sub_order" style="white-space:pre-line;" :style="{
                color: p.color,
              }" >
                            {{ p.item_describe }}
                            {{ p.voltage }}
                        </span>
                    </div>
                    <div class="divTableCell center nowrap">
                        <span v-for="p in item.sub_order" style="white-space: pre-line;" :style="{
                color: p.color,
              }">
                            {{ p.amount }}
                        </span>
                    </div>
                    <div class="divTableCell center nowrap">
                        <span v-for="p in item.sub_order" style="white-space: pre-line;text-align:center;" :style="{
                color: p.color,
              }">
                            {{ p.state }}
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
    </div>
</template>
<script>
export default {
    name: 'orders',
    data() {
        return {
            orders: [
                //主页面呈现
                {
                    id: '',
                    state_id: '',
                    date: '',
                    vendor: '',
                    sub_order: [{
                        sub_order_id: '',
                        item_describe: '',
                        amount: '',
                        voltage: '',
                        color: '',
                        state: '',
                        price: '',
                        order_id: '',
                    }, ],
                    note: '',
                    address: '',
                },
            ],
        }
    },
    computed: {

    },
    methods: {
        TEST1: function() {},

        color: function(state_id) {
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
    },
    mounted: function() {
        console.log("1")
        let me = this
        this.axios
            .get('getOrdersForViewer', {
                params: {
                },
            })
            .then((response) => {
                console.log("2")
                if (response.data == '') {
                    me.orders = []
                    return
                }
                me.orders = JSON.parse(response.data)
                let i, j, v, pl, p, c

                console.log("3")
                for (i in me.orders) {
  
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
                    console.log("4")
                }
                console.log("5")
                console.log(me.orders)

            })

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
    font-size:large;
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
    font-weight: bold;
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