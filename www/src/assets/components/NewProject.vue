<template>
    <div>
        <form id="projectAdd" class="basicForm" method="post" enctype="application/x-www-form-urlencoded">
            <h2>New Project</h2>
            <div id="error">&nbsp;</div>
            <input type="text" placeholder="Add Project Name" name="name">
            <textarea placeholder="Project Description" name="description"></textarea>
            <input type="submit" value="ADD">
        </form>
    </div>
</template>

<script lang="ts">
    export default {
        name: "NewProject"
    }
</script>

<style lang="scss" scoped>

</style>