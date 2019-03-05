import store from './assets/scripts/store'
import Router from './assets/routers/index'
import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'
import AsyncComputed from 'vue-async-computed'

import CommonHeader from '@components/CommonHeader.vue'
import LoadingOverlay from '@components/Overlay/LoadingOverlay.vue'

import './assets/styles/cssreset-min.css'
import './assets/styles/main.scss'
import AuthApi from "@scripts/api/AuthApi"

Vue.use(AsyncComputed)
Vue.use(Vuex)
Vue.use(VueRouter)

const router = new VueRouter({
        mode: 'history',
        routes: Router
    })

;(async () => {
    const token = localStorage.getItem("accessToken")
    if (token != null) {
        store.commit("setToken", token)
    }

    const res = await AuthApi.auth()
    store.commit("setIsLogin", res.success)

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
            if (!res.success) {
                this.$router.push({name: "login"})
            }
            this.headerController()
        }
    })
})()