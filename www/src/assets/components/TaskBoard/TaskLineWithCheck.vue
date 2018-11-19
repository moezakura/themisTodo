<template>
    <div class="task-line-with-check" @click="taskClick">
        <task-line :task="task" :hide-assign="hideAssign" :full-deadline="fullDeadline"
                   :allow-show-detail="allowShowDetail" class="task-line"></task-line>
        <div class="check fas" :class="{ checked: value }"></div>
    </div>
</template>

<script lang="ts">
    import TaskLine from "./TaskLine"

    export default {
        name: "TaskLineWithCheck",
        components: {TaskLine},
        props: ["task", "hideAssign", "fullDeadline", "allowShowDetail", "value"],
        methods: {
            taskClick() {
                this.$emit('input', !this.value)
                this.$emit('change', !this.value)
            }
        }
    }
</script>
<style lang="scss" scoped>
    .task-line-with-check {
        cursor: default;

        div.task-line {
            padding-left: 60px;

            &::after {
                width: calc(100% + 60px + 5px);
                margin-left: -60px;
            }
        }

        .check {
            position: relative;
            display: block;
            width: 25px;
            height: 25px;
            margin: -55px 0 28px 9px;
            border: solid 1px white;

            &::before {
                content: "\f00c";
                opacity: 0;
                text-align: center;
                margin-left: 3px;
                line-height: 25px;
                font-size: 20px;
                transition: ease opacity .1s;
            }

            &.checked::before {
                opacity: 1;
            }
        }
    }
</style>