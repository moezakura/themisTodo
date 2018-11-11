import store from './assets/scripts/store'
import Router from './assets/routers/index'
import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'
import CommonHeader from '@components/CommonHeader.vue'

import './assets/styles/cssreset-min.css'
import './assets/fontawesome/web-fonts-with-css/css/fontawesome-all.css'
import './assets/styles/main.scss'

Vue.use(Vuex)
Vue.use(VueRouter)

const router = new VueRouter({
    mode: 'history',
    routes: Router
})

new Vue({
    el: "#main",
    store,
    router,
    components: {
        CommonHeader
    },
    computed: {
        isHeaderEnable(): boolean {
            return this.$store.getters.isHeaderEnable
        }
    },
    watch: {
        '$route'(to, from) {
            this.$store.commit("setHeaderEnable", to.path !== "/")
        }
    }
})