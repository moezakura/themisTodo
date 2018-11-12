import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'
import Project from "@scripts/model/api/project/Project"
import Task from "@scripts/model/api/task/Task"

export interface RootState {
    headerEnable: boolean,
    loadingCount: number,
    currentProject: Project | undefined
    currentTask: Task | undefined
}

Vue.use(Vuex)

const store: StoreOptions<RootState> = {
    // データを保存するためのステートを作成
    state: {
        headerEnable: false,
        loadingCount: 0,
        currentProject: undefined,
        currentTask: undefined
    },
    getters: {
        isHeaderEnable: state => state.headerEnable,
        isLoadingShow: state => state.loadingCount > 0,
        getCurrentProject: state => state.currentProject,
        getCurrentTask: state => state.currentTask
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
        }
    },
}

export default new Vuex.Store<RootState>(store)
