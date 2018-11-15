<template>
    <form id="accountSettingsIcon" ref="account-settings-icon" autocomplete="off" @dragover.prevent="dragOver"
          @dragleave.prevent="dragLeave" @drop.prevent="drop" :class="{ 'drag-over': dragging }">
        <div class="success" v-show="changeSuccess" @click="clearMessages">Success changed</div>
        <div class="error" v-show="errorMessage.length > 0" @click="clearMessages">{{ errorMessage }}</div>
        <div id="accountSettingsIconFake" @click="clickIconFake"
             :style="{'background-image': `url('${userIconPath}')`}">
            <div :style="{'height': `${uploadProgress}%`}" class="uploading-image">
                <div :style="{'background-image': `url('${uploadingBlob}')`}" class="uploading-image-body"></div>
            </div>
        </div>
        <input type="file" id="accountSettingsIconFile" ref="icon-file" name="icon" @change="change">
        <p>To change the icon Drag and drop the file or click the icon</p>
        <transition>
            <div class="uploading-progress-bar" v-if="uploadProgress < 100">
                <div class="uploading-progress-bar-body" :style="{'width': `${uploadProgress}%`}"></div>
            </div>
        </transition>
    </form>
</template>

<script lang="ts">
    import User from "../../scripts/model/api/user/User"
    import AccountApi from "../../scripts/api/AccountApi"

    export default {
        name: "IconSettings",
        data: () => {
            return {
                changeSuccess: false,
                errorMessage: "",
                dragging: false,
                uploadingBlob: "",
                progress: 0
            }
        },
        computed: {
            myProfile(): User | undefined {
                return this.$store.getters.getMyProfile
            },
            userIconPath(): string {
                if (this.myProfile == undefined || this.myProfile.iconPath == undefined ||
                    this.myProfile.iconPath.length <= 0) {
                    return ""
                }
                return `/api/account/icon/${this.myProfile.iconPath}`
            },
            uploadProgress(): number {
                if (this.uploadingBlob.length <= 0) {
                    return 100
                }
                return this.progress
            }
        },
        created() {

        },
        methods: {
            clearMessages() {
                this.changeSuccess = false
                this.errorMessage = ''
            },
            clickIconFake() {
                this.$refs["icon-file"].click()
            },
            dragOver() {
                this.dragging = true
            },
            dragLeave() {
                this.dragging = false
            },
            drop(e) {
                this.dragLeave()
                if (e.dataTransfer.files.length > 1) {
                    this.errorMessage = "Too many files"
                    return
                }
                this.$refs["icon-file"].files = e.dataTransfer.files
                this.change()
            },
            change() {
                this.clearMessages()
                const files = this.$refs['icon-file'].files
                if (files.length <= 0) {
                    this.errorMessage = "non file"
                    return
                }
                const file = files[0]
                const reader = new FileReader()
                reader.onload = () => {
                    this.uploadingBlob = reader.result
                }
                reader.readAsDataURL(file)

                const uploadFormData = new FormData(this.$refs['account-settings-icon'])
                AccountApi.uploadImage(uploadFormData, progress => {
                    this.progress = progress * 100
                }).then(res => {
                    if (!res.success) {
                        this.errorMessage = res.message
                    } else {
                        this.success = true
                        let profile = this.$store.getters.getMyProfile
                        profile.iconPath = res.fileId
                        this.$store.commit("setMyProfile", profile)
                    }
                }).finally(() => {
                    this.uploadingBlob = ""
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    p {
        text-align: center;
    }

    .v-enter,
    .v-leave-to {
        opacity: 0;
        div, ul {
            transform: translateY(-5px);
        }
    }

    .v-enter-active,
    .v-leave-active {
        transition: ease all .2s;
        div, ul {
            transition: ease all .2s;
        }
    }
</style>