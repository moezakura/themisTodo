<template>
    <div>
        <p class="no-tasks" v-if="taskTimers.length === 0">There is no Doing task.</p>
        <div class="tasks" v-else v-for="i in taskTimers">
            <ul>
                <task-timer-line :task-timer="i" :class="{ active: i.endDateUnix === 0 }"
                                 @load-page="loadList"></task-timer-line>
            </ul>
        </div>
    </div>
</template>

<script lang="ts">
    import TaskTimerApi from "@scripts/api/TaskTimer";
    import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer";
    import TaskTimerLine from "@components/TaskTimer/TaskTimerLine.vue"

    interface TaskTimerUserDashboardData {
        taskTimersLoadTimer: any
        taskTimers: Array<TaskTimer>
    }

    export default {
        name: "TaskTimerUserDashboard",
        components: {TaskTimerLine},
        data(): TaskTimerUserDashboardData {
            return {
                taskTimersLoadTimer: 0,
                taskTimers: [],
            }
        },
        methods: {
            loadList(silent?: boolean): void {
                if (!silent) {
                    this.$store.commit("incrementLoadingCount");
                }
                TaskTimerApi.getMyDoingList().then(res => {
                    if (res.success) {
                        this.taskTimers = res.list;
                    }
                }).finally(() => {
                    if (!silent) {
                        this.$store.commit("decrementLoadingCount");
                    }
                });
            }
        },
        created(): void {
            this.loadList();

            this.taskTimersLoadTimer = setInterval(() => {
                this.loadList(true);
            }, 3000);
        },
        beforeDestroy(): void {
            clearInterval(this.taskTimersLoadTimer);
        }
    }
</script>

<style scoped lang="scss">
    .no-tasks {
        text-align: center;
        margin: 20px 0;
        font-weight: 900;
        opacity: 0.4;
        letter-spacing: 1.2px;
        font-size: 16px;
    }

    .tasks {
        margin: 0 0 15px 0;
    }
</style>