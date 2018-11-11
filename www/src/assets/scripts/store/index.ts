import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

export interface RootState {
    headerEnable: boolean,
    loadingCount: number
}

Vue.use(Vuex)

const store: StoreOptions<RootState> = {
    // データを保存するためのステートを作成
    state: {
        headerEnable: false,
        loadingCount: 0
    },
    getters: {
        isHeaderEnable: state => state.headerEnable,
        isLoadingShow: state => state.loadingCount > 0
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
        }
    },
}

export default new Vuex.Store<RootState>(store)
