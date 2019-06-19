<template>
    <li class="task-timer-entry">
        <div class="label"></div>
        <div class="name">{{ task.name }}</div>
        <div class="note">{{ taskTimer.note }}</div>
        <div class="date">
            <div class="time-between">
                <span>{{ taskTimer.startDateHM }}</span>
                <span>〜</span>
                <span>{{ taskTimer.endDateHM }}</span>
            </div>
            <div class="time-total">
                <span class="time-total-text">{{ taskTimer.totalHM }}</span>
            </div>
        </div>
        <div class="actions">
            <i class="fas fa-stop" @click="stopTask(task)" v-if="taskTimer.endDateUnix === 0"></i>
            <i class="fas fa-trash" @click="deleteTask(taskTimer)" v-else></i>
            <i class="fas fa-edit"></i>
            <i class="fas fa-info"></i>
        </div>
    </li>
</template>

<script lang="ts">
    import Task from "@scripts/model/api/task/Task";
    import TaskTimerApi from "@scripts/api/TaskTimer";
    import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer";

    export default {
        name: "TaskTimerLine",
        props: {
            taskTimer: TaskTimer
        },
        computed: {
            task(): Task {
                return this.taskTimer.task
            }
        },
        methods: {
            async stopTask(task: Task): Promise<void> {
                this.$store.commit("incrementLoadingCount")
                const res = await TaskTimerApi.getTaskTimerStatus(task.createDate)

                if (res.start) {
                    const toggleRes = await TaskTimerApi.toggleTimer(task.createDate)
                    if (toggleRes.success) {
                        this.loadPage()
                    }
                }
                this.$store.commit("decrementLoadingCount")
            },
            async deleteTask(taskTimer: TaskTimer): Promise<void> {
                this.$store.commit("incrementLoadingCount")

                const toggleRes = await TaskTimerApi.deleteTaskTimer(taskTimer.id)
                if (toggleRes.success) {
                    this.loadPage()
                }
                this.$store.commit("decrementLoadingCount")
            }
        }
    }
</script>

<style scoped lang="scss">
    .task-timer-entry {
        $height: 55px;
        display: flex;
        height: $height;
        margin-bottom: 5px;
        padding-bottom: 5px;
        border-bottom: solid 1px rgba(white, .3);
        flex-wrap: wrap;

        &::after {
            position: relative;
            display: block;
            content: " ";
            border-bottom: solid 3px transparent;
            width: 100%;
            height: 0;
            margin: 3px 0 0 0;
        }

        &.active::after {
            @include animation(blinkBorderBottomAnimation 1s infinite, linear);
        }

        .label {
            width: 5px;
            background-color: red;
            margin-right: 10px;
            height: $height;
        }

        .name {
            width: 40%;
            line-height: $height;
            overflow: hidden;
            font-weight: bold;
            font-size: 16px;
        }

        .note {
            width: 25%;
            overflow-y: auto;
        }

        .date {
            margin: 0 auto;

            .time-between,
            .time-total {
                height: $height / 2;
                line-height: $height / 2;
            }

            .time-total {
                text-align: right;

                .time-total-text::before {
                    content: "total:";
                    margin-right: 5px;
                }
            }
        }

        .actions {
            display: flex;
            margin-left: 5px;

            i {
                $margin: 5px;
                width: $height;
                height: $height - $margin * 2;
                line-height: $height - $margin * 2;
                margin: $margin 0;
                font-size: 16px;
                text-align: center;
                box-sizing: border-box;
                background-color: transparent;
                border: solid 1px $accentColor;
                color: white;
                padding: 0 15px;
                letter-spacing: 1px;
                transition: ease background-color .3s;

                &:hover {
                    background-color: accentColor(.1);
                }

                &:not(:last-child) {
                    border-right: 0;
                }
            }
        }
    }
</style>