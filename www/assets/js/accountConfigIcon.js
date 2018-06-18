let ProgressBar = require('progressbar.js');

class AccountConfigIcon {
    constructor() {
        this.iconUploadIconFile = document.querySelector('#accountSettingsIcon');
        if (this.iconUploadIconFile === undefined || this.iconUploadIconFile == null)
            return;

        this.iconUploadForm = document.querySelector('#accountSettingsIconFile');
        this.iconUploadFake = document.querySelector("#accountSettingsIconFake");
        this.iconUploadError = this.iconUploadIconFile.querySelector(".error");
        this.uploadXhr = null;

        let that = this;

        this.iconUploadFake.addEventListener('click', function (e) {
            that.iconUploadForm.click();
        });
        this.iconUploadFake.addEventListener('dragover', function (e) {
            that.fileDragOver(e);
        }, true);
        this.iconUploadFake.addEventListener('dragleave', function (e) {
            that.fileDragLeave(e);
        }, true);
        this.iconUploadFake.addEventListener('drop', function (e) {
            that.fileDrug(e);
        }, true);
        this.iconUploadForm.addEventListener('change', function (e) {
            that.fileChange(e);
        }, true);

        this.iconUploadFakeProgress = new ProgressBar.Circle(this.iconUploadFake, {
            strokeWidth: 2,
            easing: 'easeInOut',
            duration: 200,
            color: 'rgb(67, 160, 71)',
            trailColor: 'rgb(85, 85, 85)',
            trailWidth: 2,
            svgStyle: null
        });
    }

    fileDragOver(e) {
        e.preventDefault();
        this.iconUploadFake.classList.add('dragover');
    }

    fileDragLeave(e) {
        e.preventDefault();
        this.iconUploadFake.classList.remove('dragover');
    }

    fileDrug(e) {
        e.preventDefault();
        this.iconUploadFake.classList.remove('dragenter');
        if (e.dataTransfer.files.length > 1) {
            console.log("over files");
            return;
        }
        this.iconUploadFakeProgress.set(0);
        this.iconUploadForm.files = e.dataTransfer.files;
        this.fileChange();
    }

    fileChange() {
        let uploadFormData = new FormData(this.iconUploadIconFile);
        this.uploadFile(uploadFormData);
    }

    uploadFile(formData) {
        this.uploadXhr = new XMLHttpRequest();
        let upload = this.uploadXhr.upload;
        if (upload) {
            let that = this;
            this.uploadXhr.onreadystatechange = function (e) {
                if (this.readyState == 4) {
                    let json = JSON.parse(this.responseText);
                    if (json.success) {
                        that.iconUploadError.style.display = "none";
                        document.AccountConfig.changeSuccess.style.display = "block";
                        that.uploadedIcon();
                    } else {
                        that.iconUploadError.style.display = "block";
                        document.AccountConfig.iconUploadError.innerText = json.message;
                        that.changeSuccess.style.display = "none";
                    }
                }
            };
            upload.onerror = function (e) {
                that.uploadError(id);
            };
            upload.onabort = function (e) {
                that.uploadError(id);
            };
            upload.ontimeout = function (e) {
                that.uploadError(id);
            };
            upload.onloadstart = function (e) {
            };
            upload.onprogress = function (e) {
                console.log(e.loaded + " / " + e.total);
                let progress = e.loaded / e.total;
                that.iconUploadFakeProgress.animate(progress);
            };
        }

        this.uploadXhr.open("POST", "/account/updateIcon/" + accountUuid);
        this.uploadXhr.send(formData);
    }

    uploadedIcon() {
        let now = Math.floor((new Date()).getTime() / 1000);
        this.iconUploadFake.style.backgroundImage = "url('assets/accountIcon/" + accountUuid + ".png?t=" + now + "')";
    }
}

new AccountConfigIcon();