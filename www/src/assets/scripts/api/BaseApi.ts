import State from "@scripts/store"

export default class BaseApi {
    static getHeader() {
        let header = new Headers()
        header.append("x-access-token", State.getters.getToken)
        return header
    }
}