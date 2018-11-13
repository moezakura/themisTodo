import Welcome from '@components/Welcome.vue'
import Login from '@components/Login.vue'
import DashBoard from '@components/Dashboard.vue'
import TaskBoard from '@components/TaskBoard.vue'
import NewProject from '@components/NewProject.vue'

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
        path: '/project/new',
        component: NewProject,
        name: 'newProject'
    }, {
        path: '/project/view/:projectId',
        component: TaskBoard,
        name: 'taskBoard',
        children: [
            {
                path: 'task/:taskId',
                component: TaskBoard,
                name: 'taskDetail'
            }, {
                path: 'settings',
                component: TaskBoard,
                meta: {isSettings: true},
                name: 'projectSettings',
            }
        ]
    }
]