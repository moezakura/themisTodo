import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'
import Project from "@scripts/model/api/project/Project"
import Task from "@scripts/model/api/task/Task"
import ProjectSettingsStore from "@scripts/model/ProjectSettingsStore"
import {ProjectDetailStatus} from "@scripts/enums/ProjectDetailStatus"

export interface RootState {
    headerEnable: boolean,
    loadingCount: number,
    currentProject: Project | undefined
    currentTask: Task | undefined
    projectSettings: ProjectSettingsStore
    projectDetailStatus: ProjectDetailStatus
}

Vue.use(Vuex)

const store: StoreOptions<RootState> = {
    // データを保存するためのステートを作成
    state: {
        headerEnable: false,
        loadingCount: 0,
        currentProject: undefined,
        currentTask: undefined,
        projectSettings: new ProjectSettingsStore(),
        projectDetailStatus: ProjectDetailStatus.HIDE
    },
    getters: {
        isHeaderEnable: state => state.headerEnable,
        isLoadingShow: state => state.loadingCount > 0,
        getCurrentProject: state => state.currentProject,
        getCurrentTask: state => state.currentTask,
        getProjectSettings: state => state.projectSettings,
        getProjectDetailStatus: state => state.projectDetailStatus,
    },
    mutations: {
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
