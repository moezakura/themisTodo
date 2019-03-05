<template>
	<div id="welcome">
		<h2>Themis&nbsp;Todo</h2>
		<h3>Themis&nbsp;Todoは簡単なTODO管理アプリケーションです</h3>
		<router-link to="/login">ログイン</router-link>
	</div>
</template>

<script lang="ts">
    export default {
        name: "Welcome",
        computed: {
            isLogin(): Boolean {
                return this.$store.getters.isLogin
            }
        },
        created(): void {
            if (this.isLogin) {
                this.$router.push("/home")
            }
        }
    }
</script>

<style scoped>

</style>