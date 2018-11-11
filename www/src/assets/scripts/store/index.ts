import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

export interface RootState {
    headerEnable: boolean
}

Vue.use(Vuex)

const store: StoreOptions<RootState> = {
    // データを保存するためのステートを作成
    state: {
        headerEnable: false
    },
    getters: {
        isHeaderEnable: state => state.headerEnable
    },
    mutations: {
        setHeaderEnable(state, value) {
            state.headerEnable = value
        }
    },
}

export default new Vuex.Store<RootState>(store)
