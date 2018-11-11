import Welcome from '@components/Welcome.vue'
import Login from '@components/Login.vue'
import DashBoard from '@components/DashBoard.vue'

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
    }, {
        path: '/home',
        component: DashBoard,
        name: 'dashboard'
    }
]