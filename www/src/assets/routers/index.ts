import Welcome from '@components/Welcome.vue'
import Login from '@components/Login.vue'
import DashBoard from '@components/Dashboard.vue'
import TaskBoard from '@components/TaskBoard.vue'
import NewProject from '@components/NewProject.vue'
import AccountSettings from '@components/AccountSettings.vue'
import AdminTop from '@components/Admin/AdminTop.vue'
import AdminDashboard from '@components/Admin/AdminDashboard.vue'
import ListAccount from '@components/Admin/Account/ListAccount.vue'
import NewAccount from '@components/Admin/Account/NewAccount.vue'

import NotFound from '@components/Errors/NotFound.vue'

export default [
    {
        path: '*',
        component: NotFound,
        meta: {hideHeader: true}
    },
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
    }, {
        path: '/settings',
        component: AccountSettings,
        name: 'settings'
    }, {
        path: '/admin',
        name: 'admin',
        components: {
            default: AdminTop,
            'admin-pages': AdminDashboard,
        },
        children: [
            {
                path: 'account/list',
                components: {
                    default: AdminTop,
                    'admin-pages': ListAccount,
                },
            },
            {
                path: 'account/new',
                components: {
                    default: AdminTop,
                    'admin-pages': NewAccount,
                },
                name: 'newAccount'
            }
        ]

    }
]