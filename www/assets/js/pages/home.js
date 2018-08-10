import Vue from 'vue';
import MyTaskList from '../components/TaskListComponent.vue';
import HomeApi from '../utils/homeApi'

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
            HomeApi.getList('TODO').then(json => {
                if(!json.success) {
                    console.error("API ERROR")
                    return
                }
                this.todoList = json.task
            })
            HomeApi.getList('DOING').then(json => {
                if(!json.success) {
                    console.error("API ERROR")
                    return
                }
                this.doingList = json.task
            })
            HomeApi.getProject().then(json => {
                if(!json.success) {
                    console.error("API ERROR")
                    return
                }
                this.projects = json.project
            })
        },
        methods : { },
    })
}
