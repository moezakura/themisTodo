import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'
import Project from "@scripts/model/api/project/Project"
import Task from "@scripts/model/api/task/Task"
import ProjectSettingsStore from "@scripts/model/ProjectSettingsStore"
import {ProjectDetailStatus} from "@scripts/enums/ProjectDetailStatus"
import User from "@scripts/model/api/user/User"

export interface RootState {
    accessToken: string,
    isLogin: boolean,
    headerEnable: boolean,
    loadingCount: number,
    profile: User | undefined,
    currentProject: Project | undefined
    currentTask: Task | undefined
    projectSettings: ProjectSettingsStore
    projectDetailStatus: ProjectDetailStatus
}

Vue.use(Vuex)

const store: StoreOptions<RootState> = {
    // データを保存するためのステートを作成
    state: {
        accessToken: "",
        isLogin: false,
        headerEnable: false,
        loadingCount: 0,
        profile: undefined,
        currentProject: undefined,
        currentTask: undefined,
        projectSettings: new ProjectSettingsStore(),
        projectDetailStatus: ProjectDetailStatus.HIDE
    },
    getters: {
        getToken: state => state.accessToken,
        isLogin: state => state.isLogin,
        isHeaderEnable: state => state.headerEnable,
        isLoadingShow: state => state.loadingCount > 0,
        getMyProfile: state => state.profile,
        getCurrentProject: state => state.currentProject,
        getCurrentTask: state => state.currentTask,
        getProjectSettings: state => state.projectSettings,
        getProjectDetailStatus: state => state.projectDetailStatus,
    },
    mutations: {
        setToken(state, value) {
            state.accessToken = value
        },
        setIsLogin(state, value) {
            state.isLogin = value
        },
        setHeaderEnable(state, value) {
            state.headerEnable = value
        },
        incrementLoadingCount(state) {
            ++state.loadingCount
        },
        decrementLoadingCount(state) {
            if (state.loadingCount > 0) {
                --state.loadingCount
            }
        },
        setMyProfile(state, value) {
            state.profile = value
        },
        setCurrentProject(state, value) {
            state.currentProject = value
        },
        setCurrentTask(state, value) {
            state.currentTask = value
        },
        setProjectSettings(state, value) {
            state.projectSettings = value
        },
        setProjectSettingsProps(state, _value) {
            const key = _value["key"]
            const value = _value["value"]

            Vue.set(state.projectSettings, key, value)
        },
        setProjectDetailStatus(state, value) {
            state.projectDetailStatus = value
        },
    },
}

export default new Vuex.Store<RootState>(store)
