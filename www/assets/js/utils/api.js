export function getList(status) {
    return fetch(`/tasks/my?status=${status}`, {
        method: 'GET',
        credentials: 'same-origin'
    }).then(res => {
        return res.json()
    }).then(json => {
        if(!json.success) {
            console.log(json.message)
        } else {
            return json.task
        }
    })
}

export function getProject() {
    return fetch("/project/my", {
        method: 'GET',
        credentials: 'same-origin'
    }).then(res => {
        return res.json()
    }).then(json => {
        if(!json.success) {
            console.log(json.message)
        } else {
            return json.project
        }
    })
}