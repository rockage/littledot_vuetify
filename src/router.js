import VueRouter from "vue-router";
import orders from './components/orders'
import products from './components/products'
import crawler from './components/crawler'
import charts from './components/charts'
import vendors from './components/vendors'



const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => err)
}

export default new VueRouter({
    mode: 'history',
    routes: [
        {path: '/', component: orders},
        {path: '/orders', name: 'orders', component: orders},
        {path: '/products', name: 'products', component: products},
        {path: '/crawler', name: 'crawler', component: crawler},
        {path: '/charts', name: 'charts', component: charts},
        {path: '/vendors', name: 'vendors', component: vendors},
    ]
}) 
