var iconUploadIconFile = document.querySelector('#accountSettingsIcon'),
    iconUploadForm = document.querySelector('#accountSettingsIconFile'),
    iconUploadFake = document.querySelector("#accountSettingsIconFake"),
    iconUploadError = iconUploadIconFile.querySelector(".error"),
    uploadXhr = null;

iconUploadFake.addEventListener('click', function (e) {
    iconUploadForm.click();
});

iconUploadFake.addEventListener('dragover', fileDragOver, true);
iconUploadFake.addEventListener('dragleave', fileDragLeave, true);
iconUploadFake.addEventListener('drop', fileDrug, true);
iconUploadForm.addEventListener('change', fileChange, true);

let iconUploadFakeProgress = new ProgressBar.Circle(iconUploadFake, {
    strokeWidth: 2,
    easing: 'easeInOut',
    duration: 200,
    color: 'rgb(67, 160, 71)',
    trailColor: 'rgb(85, 85, 85)',
    trailWidth: 2,
    svgStyle: null
});

function fileDragOver(e) {
    e.preventDefault();
    iconUploadFake.classList.add('dragover');
}

function fileDragLeave(e) {
    e.preventDefault();
    iconUploadFake.classList.remove('dragover');
}

function fileDrug(e) {
    e.preventDefault();
    iconUploadFake.classList.remove('dragenter');
    if (e.dataTransfer.files.length > 1) {
        console.log("over files");
        return;
    }
    iconUploadFakeProgress.set(0);
    iconUploadForm.files = e.dataTransfer.files;
    fileChange();
}

function fileChange() {
    let uploadFormData = new FormData(iconUploadIconFile);
    console.log(uploadFormData.get("icon"));
    uploadFile(uploadFormData);
}

function uploadFile(formData) {
    uploadXhr = new XMLHttpRequest();
    let upload = uploadXhr.upload;
    if (upload) {

        uploadXhr.onreadystatechange = function(e) {
            if (this.readyState == 4) {
                let json = JSON.parse(this.responseText);

                if(json.success){
                    iconUploadError.style.display = "none";
                    changeSuccess.style.display = "block";
                    uploadedIcon();
                }else{
                    iconUploadError.style.display = "block";
                    iconUploadError.innerText = json.message;
                    changeSuccess.style.display = "none";
                }
            }
        };
        upload.onerror = function (e) {
            uploadError(id);
        };
        upload.onabort = function (e) {
            uploadError(id);
        };
        upload.ontimeout = function (e) {
            uploadError(id);
        };
        upload.onloadstart = function (e) {
        };
        upload.onprogress = function (e) {
            console.log(e.loaded + " / " + e.total);
            let progress = e.loaded / e.total;
            iconUploadFakeProgress.animate(progress);
        };
    }

    uploadXhr.open("POST", "/account/updateIcon/" + accountUuid);
    uploadXhr.send(formData);
}

function uploadedIcon() {
    let now = Math.floor((new Date()).getTime() / 1000);
    iconUploadFake.style.backgroundImage = "url('assets/accountIcon/" + accountUuid + ".png?t=" + now + "')";
}