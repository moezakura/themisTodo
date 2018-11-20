<template>
    <div>
        <at-ta :members="target" name-key="searchKey" :ats="['@', '#']" :filterMatch="searchFilter"
               :allowSpaces="false">
            <template slot="item" scope="s">
                <span v-if="targetType === 'TASKS'" class="task-id">#{{ s.item.taskId }}</span>
                <span v-text="s.item.name"></span>
            </template>
            <textarea id="taskPopupDescription" class="editor" v-model="textAreaValue"
                      ref="task-popup-description"></textarea>
        </at-ta>
    </div>
</template>

<script lang="ts">
    import AtTa from 'vue-at/dist/vue-at-textarea.js'
    import UserSearchRequest from "../../scripts/model/api/UserSearchRequest"
    import UserApi from "../../scripts/api/UserApi"
    import ProjectApi from "../../scripts/api/ProjectApi"

    export default {
        name: "TaskDetailDescription",
        components: {AtTa},
        props: ['value'],
        data: () => {
            return {
                members: [],
                tasks: [],
                targetType: "MEMBERS",
                target: [],
                timerId: null
            }
        },
        computed: {
            projectId(): number | undefined {
                if (this.$store.getters.getCurrentProject === undefined) {
                    return
                }
                const project = this.$store.getters.getCurrentProject
                return project.uuid
            },
            textAreaValue: {
                get() {
                    return this.value
                },
                set(value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        },
        created() {
            this.$store.commit("incrementLoadingCount")
            this.$store.commit("incrementLoadingCount")

            const searchRequest = new UserSearchRequest()
            searchRequest.name = "%"
            searchRequest.displayName = "%"
            searchRequest.project = this.projectId
            searchRequest.isInProject = true
            searchRequest.max = 1000

            UserApi.Search(searchRequest).then(res => {
                this.members = res
                for (let i in this.members) {
                    this.members[i].searchKey = this.members[i].name
                }
                this.target = this.members.slice(0, this.members.length)
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })

            ProjectApi.getTasks(this.projectId).then(res => {
                this.tasks = res.task
                for (let i in this.tasks) {
                    this.tasks[i].searchKey = this.tasks[i].taskId.toString()
                }
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })

            let prevText = ""
            this.timerId = setInterval(() => {
                const v = this.$refs["task-popup-description"].value
                if (v != prevText) {
                    this.$emit('input', v)
                    this.$emit('change', v)
                    prevText = v
                }
            }, 30)
        },
        destroyed() {
            clearInterval(this.timerId)
        },
        methods: {
            searchFilter(search, chunk, at) {
                if (at == "@" && this.targetType != "MEMBERS") {
                    this.target = this.members.slice(0, this.members.length)
                    this.targetType = "MEMBERS"
                } else if (at == "#" && this.targetType != "TASKS") {
                    this.target = this.tasks.slice(0, this.tasks.length)
                    this.targetType = "TASKS"
                }
                return search.toLowerCase().indexOf(chunk.toLowerCase()) > -1
            }
        }
    }
</script>

<style lang="scss" scoped>
</style>