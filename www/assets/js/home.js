import Vue from 'vue';
import MyTaskList from './components/ListComponent.vue'

if(document.querySelector("#home")) {
    new Vue({
        delimiters: ['${', '}'],
        el: '#home',
        data: {
            todoList: [],
            doingList: [],
            projects: [],
        },
        components: {
            "parent": MyTaskList
        },
        created () {
            this.getList('TODO').then(data => this.todoList = data)
            this.getList('DOING').then(data => this.doingList = data)
            this.getProject().then(data => this.projects = data)
        },
        methods : {
            getList(status) {
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
            },
            getProject() {
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
        },
    })
}
