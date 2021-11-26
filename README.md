# Littledot Manager 要点记录

#### 整体结构

- 前端UI：vue + vutify 响应式布局，后端API： golang
- 主订单视图采用原生div table方式，没有采用vutify table
- 复习：vue之所以称之为响应式布局，与传统用JS操纵DOM去完成视图更新有着很大的区别，简单地说，传统上UI产生了变化，需要用代码去读这个DOM的ID，然后用DOM.innerText之类方法更新变量。反过来，代码改变了某个变量的值，需要用同样的方法给DOM的元素再次赋值，视图才会改变。VUE这类响应式布局则不用这么麻烦，只要UI发生了改变，其绑定的变量自动发生改变，反之，代码使变量发生了改变，VUE会自动将结果更新到视图上，简化了操作。用了VUE之后，我们只需要关心DOM和变量之间的绑定关系，至于他们是如何通信以及更新之类的事情，VUE替我们代劳了。

#### 1、dialog_data

- 此变量绑定在 ``<v-dialog v-model="dialog" width="80%"> ``控件上

- 关于子订单：

  子订单绑定在

```  html
  <v-card dense elevation="6" 
          outlined dense 
          style="padding-left:10px;padding-right:10px;padding-top:15px;margin-bottom:10px;" 
          v-for="item in dialog_data.sub_order">
```

  观察v-for，它将自动从dialog_data.sub_order数组中创建不同的card，一个card对应一个子订单

  - 一个主订单可以包含数个子订单，所谓子订单就是不同的机器，比如一个客户同时买了3台机器，那么主订单包含了客户的基本信息，子订单包含了着3台不同机器的信息。

  - 关于“自定义”， 凡是没有收录在产品库（SQL表：ld_products）中的产品，都称为自定义，在产品库里以ID=38代表。比方说一个客户买了“一条RCA信号线”、“返修LDP-1000”等，这些稀奇古怪的小东西不可能全部收录在产品库里。他们都以ID=38记为“自定义”。尽管自定义的id统一为38，但描述又是各不相同的，这个描述文字就放在子订单的item_describe字段中。

    现在问题来了，一般产品都会有ID和NAME（产品名）两个字段，是一对一的关系，很容易处理。但是当ID=38的时候，其文字描述要从子订单的另一个字段item_describe去对应，即同一个id=38，需要面对无数个不同的文字描述，是一对多的关系。

    解决方法：

    1. 读：

       子订单的Select控件（实际为一个v-autocomplete，相当于带自动补全的Select），
    
       ```html
        <v-autocomplete dense v-model="item.product_list.select" 
                        :items="item.product_list.items" 
                        item-text="name" 
                        item-value="id" 
                        label="产品" 
                        @blur="onBlurCustomized($event, item)" 
                        @change="updatePrice(dialog_data.vendor_class)" 
                        style="margin-right:5px;" 
                        return-object>
       </v-autocomplete>
       ```
    
       它绑定的是dialog_data.sub_order.product_list变量，这个变量又分为两部分
    
       1. 当前选中项：product_list.select
       2. 产品列表：product_list.items
       3. return-object是什么意思呢？它将id和name完整对应起来，如果没有return-object，我们无法从select剥离出id和name来，`${select.id}`代表ID, `${select.name}`代表此id的文字描述，如果不加上这个，则select属性只能解析出id，如果尝试读取它的name的话，会报undefined错误。
    
       只要是 “自定义”的子订单，尽管其id固定为38，但是其name则千差万别。A子订单的自定义：id=38 name="插座"，
    
       B子订单的自定义：id =38 name="同轴数据线"。简单地说，只要选择的不是常规产品，每一个子订单的Select都有可能是不同的，于是我们只能为每一个子订单动态创建不同的Select ：
    
       ```javascript
        //因为有“自定义”项目的存在，每个子订单的产品列表各不相同，需要修改并分开处理：
       let item_describe = ''
       for (i in me.dialog_data.product_list.items) {
       //如果子订单中产品选择的是'自定义'则显示自定义的值：
       	if (me.dialog_data.product_list.items[i].id == '38' && 
               me.dialog_data.sub_order[j].product_id == '38') {
               item_describe = order_details[j].item_describe
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
       me.dialog_data.sub_order[j].product_list.select.id = me.dialog_data.sub_order[j].product_id
       me.dialog_data.sub_order[j].product_list.select.name = me.dialog_data.sub_order[j].item_describe
       
       ```
    
       如果是常规产品item_describe则就是该产品从SQL读出来的name（产品名）字段。
    
       如果是自定义，则item_describe从子订单储存的item_describe读出这个自定义产品的文字描述。
    
       记住：两部分，一是list(列表)，一是select(选中项)，两者都要处理。
    
    2. 写：
    
       同样，在选择了“自定义”后，我们触发一个onBlurCustomized事件，将id = 38设为select控件的当前输入的text (e.target.value)
    
       ```javascript
               onBlurCustomized: function(e, item) {
                   let i
                   if (item.product_list.select.id == '38') {
                       for (i in item.product_list.items) {
                           if (item.product_list.items[i].id == '38') {
                               item.product_list.items[i].name = e.target.value
                           }
                       }
                   }
               },
       ```
    
       
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    







