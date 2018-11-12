<template>
    <form id="taskboardAddForm" class="basicForm">
        <h3>New Task</h3>
        <div id="taskboardAddClose" @click="hideTaskAdd"><i class="fas fa-angle-double-right"></i></div>
        <div class="error" v-show="errorMessage !== undefined && errorMessage.length > 0">{{ errorMessage }}</div>
        <p>Name</p>
        <input type="text" name="name">
        <p class="label">Creator:<span>{{ displayName }}</span></p>
        <p class="label">ID:<span id="createTaskId"></span></p>
        <p>Assign</p>
        <input type="text" value="" name="assign" id="assign" autocomplete="off">
        <p>Deadline</p>
        <input type="date" value="" name="deadline">
        <p>Description</p>
        <textarea name="description"></textarea>
        <input type="submit" value="ADD">
    </form>
</template>

<script lang="ts">
    export default {
        name: "TaskAdd",
        props: ['value'],
        data: () => {
            return {
                errorMessage: "",
            }
        },
        computed: {
            displayName(): string {
                return ""
            }
        },
        methods: {
            hideTaskAdd() {
                this.$emit("input", false)
                this.$emit("change", false)
            }
        }
    }
</script>

<style scoped>

</style>