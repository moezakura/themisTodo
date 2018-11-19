<template>
    <div :class="[limitAddClass]" class="task-item" @click="showTaskDetail">
        <div class="taskTitle">{{ task.name }}</div>
        <div class="task-info">
            <div class="taskId">#{{ task.taskId }}</div>
            <div class="taskAssign" v-if="isShowAssign">
                <div class="taskAssignIcon" :style="{'background-image': `url('${backgroundImage}')`}"></div>
                <div class="taskAssignName">{{ task.assignName }}</div>
            </div>
            <div class="taskLimit" :class="{ 'no-assign': !isShowAssign }">
                <i class="fas fa-calendar-alt"></i>
                <span class="deadlineDate">{{ isFullDeadline ? task.deadline : task.deadlineMD }}</span>
                <span class="deadlineDate deadline-limit" v-if="!isCompleted">{{ task.limitDate }}日</span>
                <span class="deadlineDate deadline-limit" v-if="isCompleted">DONE!</span>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
    export default {
        name: "TaskLine",
        props: ["task", "hideAssign", "fullDeadline", "allowShowDetail"],
        data() {
            return {
                backgroundImage: `/api/account/icon/${this.task.assignIconPath}`
            }
        },
        computed: {
            isCompleted(): boolean {
                if (this.task === undefined || this.task.status === undefined || this.task.status === null) {
                    return false
                }
                return this.task.status === 3
            },
            isFullDeadline(): boolean {
                if (this.fullDeadline === undefined || this.fullDeadline == null) {
                    return false
                }
                return this.fullDeadline
            },
            isShowAssign(): boolean {
                if (this.hideAssign === undefined || this.hideAssign == null) {
                    return true
                }
                return !this.hideAssign
            },
            limitAddClass(): string {
                if (this.task.status == 3) {
                    return "normal"
                }
                if (this.task.limitDate <= 0) {
                    return "over"
                }
                else if (this.task.limitDate <= 1) {
                    return "limit1"
                }
                else if (this.task.limitDate <= 2) {
                    return "limit2"
                }
                else if (this.task.limitDate <= 3) {
                    return "limit3"
                }

                return "normal"
            }
        },
        methods: {
            showTaskDetail() {
                if (this.allowShowDetail == undefined || this.allowShowDetail == true) {
                    this.$router.push({
                        name: "taskDetail", params: {
                            projectId: this.task.projectId,
                            taskId: this.task.taskId,
                        }
                    })
                }
            }
        }
    }
</script>