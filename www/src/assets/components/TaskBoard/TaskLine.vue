<template>
    <li>
        <div class="taskTitle">{{ task.name }}</div>
        <div class="taskId">#{{ task.taskId }}</div>
        <div class="taskAssign">
            <div class="taskAssignIcon" v-if="isShowAssign"
                 :style="{'background-image': `url('${backgroundImage}')`}"></div>
            <div class="taskAssignName" v-if="isShowAssign">{{ task.assignName }}</div>
        </div>
        <div class="taskLimit">
            <i class="fas fa-calendar-alt"></i>
            <span class="deadlineDate">{{ task.deadline }}</span>
            <span class="deadlineDate" v-if="!isCompleted">あと{{ todo.limitDate }}日</span>
            <span class="deadlineDate" v-else>Completed!</span>
        </div>
    </li>
</template>

<script lang="ts">
    export default {
        name: "TaskLine",
        props: ["task", "hideAssign"],
        data() {
            return {
                backgroundImage: `/assets/accountIcon/${this.task.assign}.png?t=${new Date().getTime()}`
            }
        },
        computed: {
            isCompleted() {
                return this.task.status === undefined || this.task.status === null || this.task.status !== 3
            },
            isShowAssign() {
                if (this.hideAssign === undefined || this.hideAssign == null) {
                    return true
                }
                return !this.hideAssign
            }
        }
    }
</script>