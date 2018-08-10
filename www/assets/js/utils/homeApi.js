export default class HomeApi {
    static getList(status) {
        return fetch(`/tasks/my?status=${status}`, {
            method: 'GET',
            credentials: 'same-origin'
        }).then(res => {
            return res.json()
        })
    }

    static getProject() {
        return fetch("/project/my", {
            method: 'GET',
            credentials: 'same-origin'
        }).then(res => {
            return res.json()
        })
    }
}