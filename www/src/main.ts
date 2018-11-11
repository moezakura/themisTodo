import store from './assets/scripts/store'
import Router from './assets/routers/index'
import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'

import CommonHeader from '@components/CommonHeader.vue'
import LoadingOverlay from '@components/Overlay/LoadingOverlay.vue'

import './assets/styles/cssreset-min.css'
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
        CommonHeader,
        LoadingOverlay
    },
    computed: {
        isHeaderEnable(): boolean {
            return this.$store.getters.isHeaderEnable
        }
    },
    watch: {
        '$route'(to, from) {
            this.headerController()
        }
    },
    methods: {
        headerController() {
            let showHeader = true
            if (this.$route.meta !== undefined && this.$route.meta.hideHeader !== undefined) {
                showHeader = !this.$route.meta.hideHeader
            }
            this.$store.commit("setHeaderEnable", showHeader)
        }
    },
    created() {
        this.headerController()
    }
})