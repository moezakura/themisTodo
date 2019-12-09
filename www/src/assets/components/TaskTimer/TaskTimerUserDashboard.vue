<template>
    <div>
        <p class="no-tasks" v-if="timerTasks.length === 0">There is no Doing task.</p>
    </div>
</template>

<script lang="ts">
    import Task from "@scripts/model/api/task/Task";

    interface TaskTimerUserDashboardData {
        timerTasks: Array<Task>
    }

    export default {
        name: "TaskTimerUserDashboard",
        data(): TaskTimerUserDashboardData {
            return {
                timerTasks: [],
            }
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
</style>