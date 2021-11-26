import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import VueRouter from 'vue-router'
import axios from 'axios'

Vue.use(VueRouter)
Vue.config.productionTip = false

Vue.prototype.axios = axios

axios.defaults.withCredentials = true
 //axios.defaults.baseURL = 'https://rockage.net/api/'
 axios.defaults.baseURL = 'http://127.0.0.1:5050/'


new Vue({
    router,
    store,
    vuetify,
    render: function(h) {
        return h(App);
    }
}).$mount('#app')

/* 
方案1，ES5传统函数语法：
render: function (createElement) {
     return createElement(App);
}

方案2，ES6箭头函数语法：
render: createElement => createElement(App)


方案3，ES6语法 + h :
render: h => h(App) // h来自单词 hyperscript，等效于createElement

h 和 createElement 函数原型：
ƒ (a, b, c, d) { return createElement(vm, a, b, c, d, true); }
*/