import store from './assets/scripts/store'
import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'

import './assets/styles/cssreset-min.css'
import './assets/fontawesome/web-fonts-with-css/css/fontawesome-all.css'
import './assets/styles/main.scss'

Vue.use(Vuex)
Vue.use(VueRouter)

new Vue({
    el: "#main",
    store,
})