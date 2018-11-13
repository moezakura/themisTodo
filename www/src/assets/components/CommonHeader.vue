<template>
    <header>
        <router-link to="/"><img src="/assets/images/logo.png"></router-link>
        <router-link to="/home">
            <div>Home</div>
        </router-link>
        <router-link to="/project/new">
            <div>New Project</div>
        </router-link>
        <router-link to="/settings">
            <div>My Setting</div>
        </router-link>
        <router-link to="/account/add">
            <div>New Account</div>
        </router-link>
    </header>
</template>

<script lang="ts">
    export default {
        name: "CommonHeader"
    }
</script>

<style scoped>

</style>