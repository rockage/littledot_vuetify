<template>
	<!--主界面-->
	<div class="text-caption">
		<div class="table">
			<div class="head">
				<div class="headcell" style="max-width:5%;">
					<v-btn icon color="#00B0FF" @click="addItem" x-small>
						<v-icon>mdi-plus</v-icon>
					</v-btn>
				</div>
				<div class="headcell" style="max-width:9%;">产品名称</div>
				<div class="headcell" style="max-width:3%;">零售</div>
				<div class="headcell" style="max-width:3%;">批发</div>
				<div class="headcell" style="max-width:26%;">淘宝显示名称</div>
				<div class="headcell" style="max-width:8%;">分类</div>
				<div class="headcell">备注</div>
				<div class="headcell" style="max-width:3%;">
					<v-btn x-small @click="sortProductPosition()">重量</v-btn>
				</div>
				<div class="headcell" style="max-width:4%;">
					<v-btn x-small @click="updateProductPosition()">位置</v-btn>
				</div>
			</div>

			<div class="tr" v-for="item in products" :key="item.id">
				<div class="td" style="max-width:5%;" :id="item.id">
					<v-btn
						x-small
						class="ma-2"
						outlined
						color="indigo"
						style="margin-left: 0px"
						@click="editItem"
					>
						{{ item.id }}
					</v-btn>
				</div>
				<div class="td" style="max-width:9%;">{{ item.name }}</div>
				<div class="td" style="max-width:3%;">{{ item.price1 }}</div>
				<div class="td" style="max-width:3%;">{{ item.price2 }}</div>
				<div class="td" style="max-width:26%;">{{ item.tb_name }}</div>
				<div class="td" style="max-width:8%;">{{ item.class }}</div>
				<div class="td">{{ item.note }}</div>
				<div class="td" style="max-width:3%;">{{ item.weight }}</div>
				<div class="td" style="max-width:4%;justify-content: center;">
					<div
						contenteditable="true"
						style="width:30px;
							 height:22px;
							 background-color:#ECF0F1;
							 text-align:center;border:0px;
							 ime-mode:Disabled;
							 "
					>
						{{ item.position }}
					</div>
				</div>
			</div>
		</div>
		<!-- 编辑弹窗 -->
		<v-dialog v-model="dialog" max-width="50%">
			<v-card>
				<v-card-title>
					<span class="headline">编辑产品</span>
				</v-card-title>
				<v-card-text>
					<v-container>
						<v-row>
							<v-col cols="12" sm="12" md="12">
								<v-text-field
									v-model="product_dialog.name"
									label="产品名称"
								></v-text-field>
							</v-col>
						</v-row>
						<v-row>
							<v-col cols="12" sm="12" md="12">
								<v-text-field
									v-model="product_dialog.tb_name"
									label="淘宝显示名称"
								></v-text-field>
							</v-col>
						</v-row>
						<v-row>
							<v-col cols="12" sm="6" md="6">
								<v-text-field
									v-model="product_dialog.price1"
									label="零售价"
								></v-text-field>
							</v-col>
							<v-col cols="12" sm="6" md="6">
								<v-text-field
									v-model="product_dialog.price2"
									label="批发价"
								></v-text-field>
							</v-col>
						</v-row>
						<v-row>
							<v-col cols="12" sm="6" md="6">
								<v-select
									label="分类"
									v-model="product_dialog.class_selected"
									:items="products_class_list"
									item-text="name"
									item-value="id"
									return-object
								></v-select>
							</v-col>
							<v-col cols="12" sm="6" md="6">
								<v-text-field
									v-model="product_dialog.weight"
									label="重量"
								></v-text-field>
							</v-col>
						</v-row>
						<v-row>
							<v-col cols="12" sm="12" md="12">
								<v-textarea
									v-model="product_dialog.note"
									filled
									label="备注"
								></v-textarea>
							</v-col>
						</v-row>
					</v-container>
				</v-card-text>
				<v-card-actions>
					<v-btn :disabled="del_button" color="pink" text @click="deleteProduct"
						>DELETE</v-btn
					>
					<v-spacer></v-spacer>
					<v-btn color="blue darken-1" text @click="dialog = false">
						Cancel
					</v-btn>
					<v-btn color="blue darken-1" text @click="updateProduct()">
						Save
					</v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>
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
	</div>
</template>
<script>
export default {
	name: 'products',
	components: {
		
	},
	data() {
		return {
			alert: false,
			drag: false,
			del_button: false,
			//定义要被拖拽对象的数组
			products: [],
			products_class_list: [],
			dialog: false,
			product_dialog: {
				id: '',
				name: '',
				price1: '',
				price2: '',
				tb_name: '',
				class_selected: {
					id: '',
					name: '',
				},
				note: '',
				weight: '',
				position: '',
			},
			localStorage: window.localStorage,
			editedItem: '',
		}
	},
	methods: {
		sortProductPosition() {
			// 按当前视图生成排序，不存盘
			let p = 1 // 从1开始
			let i = 2
			let id,
				position = ''

			while (true) {
				id = document
					.evaluate(
						'/html/body/div/div/main/div/div/div[1]/div[' +
							i +
							']/div[1]/button/span',
						document
					)
					.iterateNext()
				if (id == null) {
					break
				} else {
					position = document
						.evaluate(
							'/html/body/div/div/main/div/div/div[1]/div[' +
								i +
								']/div[9]/div',
							document
						)
						.iterateNext()
					i++
				}
				position.innerText = p
				p++
			}
		},
		updateProductPosition(e, item) {
			// 将排序（注意不是视图）存盘
			let i = 2
			let id,
				position = ''
			let newPosition = new Array()
			while (true) {
				id = document
					.evaluate(
						'/html/body/div/div/main/div/div/div[1]/div[' +
							i +
							']/div[1]/button/span',
						document
					)
					.iterateNext()
				if (id == null) {
					break
				} else {
					position = document
						.evaluate(
							'/html/body/div/div/main/div/div/div[1]/div[' +
								i +
								']/div[9]/div',
							document
						)
						.iterateNext()
					i++
				}
				newPosition.push({
					id: id.innerText,
					position: position.innerText,
				})
			}
			let me = this
			let param = new URLSearchParams()
			param.append('newPosition', JSON.stringify(newPosition))
			console.log(newPosition)
			this.axios.post('updateProductPosition', param).then((response) => {
				if (response.data == 'sucess') {
					me.alert = !me.alert
					this.initView()
				}
			})
		},
		updateProduct: function(e) {
			let param = new URLSearchParams()
			param.append('pid', this.product_dialog.id)
			param.append('pname', this.product_dialog.name)
			param.append('tb_name', this.product_dialog.tb_name)
			param.append('weight', this.product_dialog.weight)
			param.append('class', this.product_dialog.class_selected.id)
			param.append('note', this.product_dialog.note)
			param.append('price1', this.product_dialog.price1)
			param.append('price2', this.product_dialog.price2)
			let me = this
			this.axios.post('updateProduct', param).then((response) => {
				if (response.data == 'sucess') {
					me.initView()
					me.dialog = false
				}
			})
		},
		deleteProduct: function(e) {
			let param = new URLSearchParams()
			param.append('pid', this.product_dialog.id)
			let me = this
			this.axios.post('deleteProduct', param).then((response) => {
				if (response.data == 'sucess') {
					me.initView()
					me.dialog = false
				}
			})
		},
		addItem: function(e) {
			this.product_dialog.name = ''
			this.product_dialog.id = ''
			this.product_dialog.note = ''
			this.product_dialog.price1 = 0
			this.product_dialog.price2 = 0
			this.product_dialog.tb_name = ''
			this.product_dialog.weight = ''
			this.product_dialog.position = 0
			this.product_dialog.class_selected = {
				name: this.products_class_list[0].name,
				id: this.products_class_list[0].id,
			}
			this.del_button = true
			this.dialog = true
		},
		editItem: function(e) {
			let pid = e.srcElement.innerText
			let item = {}

			let i = 0
			for (i in this.products) {
				if (this.products[i].id == pid) {
					console.log(pid)
					item = this.products[i]
					break
				}
			}
			this.product_dialog.name = item.name
			this.product_dialog.id = item.id
			this.product_dialog.note = item.note
			this.product_dialog.price1 = item.price1
			this.product_dialog.price2 = item.price2
			this.product_dialog.tb_name = item.tb_name
			this.product_dialog.weight = item.weight
			this.product_dialog.position = item.position
			i = 0
			for (i in this.products_class_list) {
				if (this.products_class_list[i].name == item.class) {
					this.product_dialog.class_selected = {
						name: this.products_class_list[i].name,
						id: this.products_class_list[i].id,
					}
					break
				}
			}
			this.del_button = false
			this.dialog = true
		},
		checkListCache: function() {
			let lists = ['products_class_list']
			let re_build = false
			let i = 0
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
					.get('getProductsClassList', {
						params: {},
					})
					.then((response) => {
						for (i in lists) {
							this.localStorage.setItem(lists[i], response.data)
						}
					})
			}
			this.restoreList() //直接从localStorage恢复List
		},
		restoreList: function() {
			this.products_class_list = JSON.parse(
				this.localStorage.getItem('products_class_list')
			)
		},
		initView: function() {
			this.axios
				.get('getProducts', {
					params: {},
				})
				.then((response) => {
					if (response.data == '') {
						alert('ERROR')
					}
					this.products = JSON.parse(response.data)
					//this.products.splice(0, 2)
					this.checkListCache() //检查list本地缓存
				})
		},
	},
	mounted: function() {
		this.initView()
	},
}
</script>
<style>
.box {
	display: flex;
}
/*定义要拖拽元素的样式*/
.ghost {
	background-color: #90caf9 !important;
}
/* DIV Table: */
.table {
	display: flex;
	flex-direction: column;
	margin: 0px;
}
.tr {
	display: flex;
}
.tr:hover {
	background-color: #e0e0e0;
}
.td {
	padding: 5px;
	display: flex;
	color: #616161;
	align-items: center;
	flex: 1;
	border-width: 1px;
	border-style: ridge;
}
.head {
	display: flex;
}
.headcell {
	display: flex;
	color: black;
	background: #e0e0e0;
	font-size: small;
	padding: 5px;
	align-items: center;
	justify-content: center;
	flex: 1;
	border-width: 1px;
	border-style: ridge;
}

.alert {
	position: fixed;
	left: 30%;
	top: 30%;
	z-index: 11;
}
</style>
