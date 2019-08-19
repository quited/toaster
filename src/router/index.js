import Vue from '../../node_modules/vue/dist/vue.js'
import Router from 'vue-router'
import indexT from '../components/indexT'
import mainPageT from '../components/mainPageT'

Vue.use(Router)

export default new Router({
    routes: [
        {
            path: '/',
            name: 'index',
            component: indexT
        },
        {
            path: '/mainPage',
            name: 'mainPage',
            component: mainPageT
        }
    ]
})