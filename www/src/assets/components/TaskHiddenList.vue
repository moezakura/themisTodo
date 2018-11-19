<template>
    <div>
        <div>
            <h2>Menu</h2>
            <ul>
                <li>
                    <router-link :to="{name: 'hiddenTasks', params:{projectId: projectId}}">HiddenTasks</router-link>
                </li>
            </ul>
        </div>
        <ul class="taskList">
            <task-line v-for="task in hiddenTasks" :task="task" :allowShowDetail="false"></task-line>
        </ul>
    </div>
</template>

<script lang="ts">
    import Task from "../scripts/model/api/task/Task"
    import TaskApi from "../scripts/api/TaskApi"
    import TaskSearchRequest from "../scripts/model/api/TaskSearchRequest"
    import TaskStatusConvert, {TaskStatus} from "../scripts/enums/TaskStatus"
    import TaskLine from "./TaskBoard/TaskLine"

    export default {
        name: "TaskHiddenList",
        components: {TaskLine},
        data: () => {
            const hiddenTasks: Array<Task> = []
            return {
                hiddenTasks: hiddenTasks,
            }
        },
        computed: {
            projectId(): number {
                if (this.$route.params["projectId"] == undefined) {
                    return 0
                }
                return this.$route.params["projectId"]
            }
        },
        created() {
            let searchRequest = new TaskSearchRequest()
            searchRequest.projectId = this.projectId
            searchRequest.status = TaskStatusConvert.toNumber(TaskStatus.HIDE)

            this.$store.commit("incrementLoadingCount")
            TaskApi.search(searchRequest).then(res => {
                this.hiddenTasks = res.task
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })
        }
    }
</script>

<style lang="scss" scoped>

</style>