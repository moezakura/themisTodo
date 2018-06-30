import Vue from 'vue';
import MyTaskList from '../components/ListComponent.vue';
import * as api from '../utils/api'

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
            api.getList('TODO').then(data => this.todoList = data)
            api.getList('DOING').then(data => this.doingList = data)
            api.getProject().then(data => this.projects = data)
        },
        methods : { },
    })
}
