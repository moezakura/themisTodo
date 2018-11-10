<template>
    <li :id="'singleTask-'+todo.createDate" :data-id="todo.createDate" :data-task-id="todo.taskId">
        <div class="taskTitle">{{todo.name}}</div>
        <div class="taskId">#{{todo.taskId}}</div>
        <div class="taskAssign">
            <div class="taskAssignIcon" :style="imageData"></div>
            <div class="taskAssignName">{{todo.assignName}}</div>
        </div>
        <div class="taskLimit">
            <i class="fas fa-calendar-alt"></i>
            <span class="deadlineDate">{{todo.deadline}}</span>
            <span class="deadlineDate" v-if="todo.status === undefined || todo.status === null || todo.status !== 3">あと{{todo.limitDate}}日</span>
            <span class="deadlineDate" v-else>Completed!</span>
        </div>
    </li>
</template>

<script>
    export default {
        props: ['todo'],
        data () {
            return {
                imageData: {
                    backgroundImage: "url(\"/assets/accountIcon/" + this.todo.assign + ".png?t=" + new Date().getTime() + "\")"
                }
            }
        }
    }
</script>