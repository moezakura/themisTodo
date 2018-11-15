<template>
    <div id="accountSettings" class="basicForm">
        <icon-settings></icon-settings>
        <id-setting></id-setting>
        <display-name-setting></display-name-setting>
        <password-change></password-change>
    </div>
</template>

<script lang="ts">
    import IdSetting from "@components/AccountSettings/IdSetting"
    import DisplayNameSetting from "@components/AccountSettings/DisplayNameSetting"
    import PasswordChange from "@components/AccountSettings/PasswordChange"
    import User from "../scripts/model/api/user/User"
    import IconSettings from "./AccountSettings/IconSetting"

    export default {
        name: "AccountSettings",
        components: {
            IconSettings,
            IdSetting,
            DisplayNameSetting,
            PasswordChange
        },
        computed: {},
    }
</script>

<style lang="scss" scoped>

</style>