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

                    u = this.userCache.find(user => {
                        return user.name == name
                    })
                    if (u == undefined) {
                        let user = new User()
                        user.name = name
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
        descriptionArray.push("\n")
        let createHTML: Array<VNode> = []
        let option = {
            isSharp: false,
            isAt: false,
            tmpStartWithH: "",
            isUrl: false,
            isNoBuffer: false,
            StartWithCheck: "",
        }
        for (let i = 0; i < descriptionArray.length; i++) {
            const c = descriptionArray[i]
            switch (c) {
                case " ":
                case "\n":
                    if (option.isUrl) {
                        option.tmpStartWithH = ""
                        createHTML.push(createElement('a', {
                            attrs: {
                                href: buff,
                                target: "_blank"
                            },
                        }, buff))
                        buff = ""
                        option.isSharp = false
                        option.isAt = false
                    } else if (option.isSharp || option.isAt) {
                        if (option.isSharp) {
                            const id = buff.trim().slice(1)
                            createHTML.push(createElement('span', {
                                attrs: {
                                    class: "task-id"
                                },
                                on: {
                                    click: () => this.moveTask(id)
                                },
                            }, `#${id}`))
                        } else if (option.isAt) {
                            const name = buff.trim().slice(1)
                            this.getUser(name)
                            const user = this.userCache.find(user => {
                                return user.name == name
                            })
                            let userIcon = `/api/account/icon/noIcon`
                            let addClass = ""
                            let userNameText = "@" + name
                            if (user != undefined && user.iconPath != undefined) {
                                userIcon = `/api/account/icon/${user.iconPath}`
                                addClass = "task-user"
                                userNameText = name


                                createHTML.push(createElement('span', {
                                    attrs: {
                                        class: addClass
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
                                    createElement('span', userNameText),
                                ]))
                            } else {
                                createHTML.push(createElement('span', userNameText))
                            }
                        }

                        buff = ""
                    }

                    option.isSharp = false
                    option.isAt = false
                    option.isUrl = false
                    if (descriptionArray.length > i + 1 && c == " ") {
                        const n = descriptionArray[i + 1]
                        option.isSharp = n === '#'
                        option.isAt = n === '@'
                    }

                    if (c == " ") {
                        createHTML.push(createElement('span', buff))
                        createHTML.push(createElement('span', " "))
                        buff = " "
                    } else if (c == "\n") {
                        if (!option.isSharp && !option.isAt && buff.length > 0) {
                            createHTML.push(createElement('span', buff))
                            buff = ""
                        }
                        createHTML.push(createElement('br'))
                    }

                    break
                case '#':
                    option.isSharp = true
                    buff += c
                    break
                case '@':
                    option.isAt = true
                    buff += c
                    break
                default:
                    if (!option.isNoBuffer) {
                        buff += c
                    }

                    if (option.tmpStartWithH.length > 0 || c.toLowerCase() == "h") {
                        option.tmpStartWithH += c

                        if (option.tmpStartWithH == "http://" || option.tmpStartWithH == "https://") {
                            option.isUrl = true
                            buff = buff.slice(0, -option.tmpStartWithH.length)
                            createHTML.push(createElement('span', buff))
                            buff = option.tmpStartWithH
                        }
                    }
                    break
            }
        }
        if (buff.length > 0) {
            createHTML.push(createElement('span', buff))
        }

        return createElement('div', createHTML)
    }
})