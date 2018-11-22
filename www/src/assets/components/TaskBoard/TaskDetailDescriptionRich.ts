import Vue, {VNode} from 'vue'
import Task from "@scripts/model/api/task/Task"
import UserApi from "@scripts/api/UserApi"
import UserSearchRequest from "@scripts/model/api/UserSearchRequest"
import User from "@scripts/model/api/user/User"
import {Mutex} from 'async-mutex'

const mutex = new Mutex()
export default Vue.component('TaskDetailDescriptionRich', {
    name: "TaskDetailDescriptionRich",
    props: ['task'],
    data: () => {
        const userCache: Array<User> = []
        return {
            userCache: userCache,
        }
    },
    methods: {
        moveTask(taskId) {
            let task: Task = this.$store.getters.getCurrentTask
            if (task != undefined && task.projectId != undefined) {
                console.log(this.$router)
                this.$router.push({
                    name: "taskDetail",
                    params: {
                        projectId: task.projectId,
                        taskId: taskId
                    }
                })
            }
        },
        async getUser(name: string) {
            let u: User
            await mutex.runExclusive(async () => {
                u = this.userCache.find(user => {
                    return user.name == name
                })

                if (u == undefined) {
                    let searchRequest: UserSearchRequest = new UserSearchRequest()
                    searchRequest.name = name
                    searchRequest.max = 1
                    searchRequest.isInProject = true

                    let users = await UserApi.Search(searchRequest)
                    for (const user of users) {
                        this.userCache.push(user)
                    }
                }
            })
            return
        }
    },
    render: function (createElement) {
        let description = ""
        try {
            let task = this.$store.getters.getCurrentTask
            description = task.description
        } catch (e) {
            console.error(e)
        }

        let buff = ""
        let descriptionArray = description.split('')
        let createHTML: Array<VNode | string> = []
        let option = {
            isSharp: false,
            isAt: false,
        }
        for (const _i in descriptionArray) {
            const i = Number(_i)
            const c = descriptionArray[i]
            switch (c) {
                case " ":
                case "\n":
                    if (option.isSharp || option.isAt) {
                        if (option.isSharp) {
                            const id = buff.trim()
                            createHTML.push(createElement('span', {
                                attrs: {
                                    class: "task-id"
                                },
                                on: {
                                    click: () => this.moveTask(id)
                                },
                            }, `#${id}`))
                        } else if (option.isAt) {
                            const name = buff.trim()
                            this.getUser(name)
                            const user = this.userCache.find(user => {
                                return user.name == name
                            })
                            let userIcon = `/api/account/icon/noIcon`
                            if (user != undefined && user.iconPath != undefined) {
                                userIcon = `/api/account/icon/${user.iconPath}`
                            }

                            createHTML.push(createElement('span', {
                                attrs: {
                                    class: "task-user"
                                },
                            }, [
                                createElement('div', {
                                    attrs: {
                                        class: "task-user-icon",
                                    },
                                    style: {
                                        "background-image": `url('${userIcon}')`,
                                    }
                                }),
                                createElement('span', `@${name}`),
                            ]))
                        }

                        option.isSharp = false
                        option.isAt = false
                        if (descriptionArray.length > i + 1) {
                            const n = descriptionArray[i + 1]
                            option.isSharp = n === '#'
                            option.isAt = n === '@'
                        }

                        buff = ""
                    }

                    if (c == " ") {
                        createHTML.push(createElement('span', buff))
                        createHTML.push(createElement('span', " "))
                        buff = " "
                    } else if (c == "\n") {
                        if (!option.isSharp && !option.isAt && buff.length > 0) {
                            createHTML.push(buff)
                            buff = ""
                        }
                        createHTML.push(createElement('br'))
                    }

                    break
                case '#':
                    option.isSharp = true
                    break
                case '@':
                    option.isAt = true
                    break
                default:
                    buff += c
                    break
            }
        }
        if (buff.length > 0) {
            createHTML.push(buff)
        }

        return createElement('div', createHTML)
    }
})