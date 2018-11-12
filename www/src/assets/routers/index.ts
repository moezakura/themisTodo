import Welcome from '@components/Welcome.vue'
import Login from '@components/Login.vue'
import DashBoard from '@components/DashBoard.vue'
import TaskBoard from '@components/TaskBoard.vue'

export default [
    {
        path: '/',
        component: Welcome,
        name: 'welcome',
        meta: {hideHeader: true}
    }, {
        path: '/login',
        component: Login,
        name: 'login',
        meta: {hideHeader: true}
    }, {
        path: '/home',
        component: DashBoard,
        name: 'dashboard'
    }, {
        path: '/project/view/:projectId',
        component: TaskBoard,
        name: 'taskBoard'
    }, {
        path: '/project/view/:projectId/task/:taskId',
        component: TaskBoard,
        name: 'taskDetail'
    }
]