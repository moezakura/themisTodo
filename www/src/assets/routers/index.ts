import Welcome from '@components/Welcome.vue'
import Login from '@components/Login.vue'

export default [
    {
        path: '/',
        component: Welcome,
        name: 'welcome',
        meta: { hideHeader: true }
    }, {
        path: '/login',
        component: Login,
        name: 'login',
        meta: { hideHeader: true }
    }
]