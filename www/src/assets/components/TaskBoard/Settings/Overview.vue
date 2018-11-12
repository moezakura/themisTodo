<template>
    <form id="projectConfigForm" class="basicForm" autocomplete="off">
        <p>Name</p>
        <input type="text" name="name" value="">
        <p>Description</p>
        <textarea name="description"></textarea>
        <input type="submit" value="SAVE">
    </form>
</template>

<script lang="ts">
    export default {
        name: "Overview"
    }
</script>

<style scoped>

</style>