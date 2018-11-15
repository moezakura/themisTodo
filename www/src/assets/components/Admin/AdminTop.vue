<template>
    <div>
        <div>
            <ul>
                <li>
                    <div>Accounts</div>
                    <ul>
                        <li><router-link to="/admin/account/new">New Account</router-link></li>
                    </ul>
                </li>
            </ul>
        </div>
        <div>
            <router-view name="admin-pages"></router-view>
        </div>
    </div>
</template>

<script>
    export default {
        name: "AdminTop"
    }
</script>

<style scoped>

</style>