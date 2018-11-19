<template>
    <div class="project-children">
        <div class="left">
            <h2>Menu</h2>
            <ul>
                <li>
                    <router-link :to="{name: 'taskBoard', params:{projectId: projectId}}">TaskBoard</router-link>
                </li>
                <li>
                    <router-link :to="{name: 'hiddenTasks', params:{projectId: projectId}}">HiddenTasks</router-link>
                </li>
                <li>
                    Reload
                </li>
            </ul>
        </div>
        <div class="content">
            <section class="hidden-task-list">
                <ul class="taskList">
                    <li v-for="task in hiddenTasks">
                        <task-line-with-check :task="task" :allowShowDetail="false"
                                              v-model="task.check"></task-line-with-check>
                    </li>
                </ul>
            </section>
            <section class="actions">
                <h2>選択項目の一括操作</h2>
                <ul class="action-buttons">
                    <li>Todo</li>
                    <li>Doing</li>
                    <li>PullRequest</li>
                    <li>Done</li>
                    <li class="delete">Delete</li>
                </ul>
            </section>
        </div>
    </div>
</template>

<script lang="ts">
    import TaskApi from "../scripts/api/TaskApi"
    import TaskSearchRequest from "../scripts/model/api/TaskSearchRequest"
    import TaskStatusConvert, {TaskStatus} from "../scripts/enums/TaskStatus"
    import TaskLineWithCheck from "./TaskBoard/TaskLineWithCheck"
    import TaskWithCheck from "../scripts/model/api/task/TaskWithCheck"

    export default {
        name: "TaskHiddenList",
        components: {TaskLineWithCheck},
        data: () => {
            const hiddenTasks: Array<TaskWithCheck> = []
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
    .project-children {
        display: flex;
        height: calc(100% - 75px);

        h2 {
            text-align: center;
            font-size: 18px;
            font-weight: 900;
            letter-spacing: 2px;
        }

        .left {
            width: 200px;
            margin-right: 20px;
        }
        .content {
            width: calc(100% - 200px - 20px);
            height: 100%;
            display: flex;

            .hidden-task-list {
                width: calc(100% - 240px);
                height: 100%;
                overflow-y: auto;

                .taskList {
                    display: flex;
                    flex-wrap: wrap;

                    li {
                        width: 100%;
                        box-sizing: border-box;
                    }
                }
            }

            .actions {
                width: 220px;
                margin: 0 10px;

                li {
                    text-align: center;
                    font-size: 16px;
                    height: 50px;
                    line-height: 50px;
                    border-bottom: solid 1px rgba(white, .3);
                    letter-spacing: 1px;
                    cursor: default;
                    transition: ease background-color .3s;

                    &:last-child{
                        border-bottom: 0;
                    }

                    &:hover {
                        background-color: rgba(67, 160, 71, .3);
                    }

                    &.delete:hover {
                        background-color: rgba(239, 83, 80, .3);
                    }
                }
            }
        }
    }

    @media screen and (min-width: 1100px) {
        .project-children .content .hidden-task-list .taskList li {
            width: calc(50% - 10px);
            margin: 0 5px;
        }
    }

    @media screen and (min-width: 1400px) {
        .project-children .content .hidden-task-list .taskList li {
            width: calc(33.33% - 10px);
            margin: 0 5px;
        }
    }
</style>