'use strict';

var errorElement = document.querySelector('errorMsg')
var localvideo = document.querySelector('video')
var joinBtn = document.getElementById('join')

joinBtn.onclick = join

var constraints = {
    audio: false,
    video: true
};

function handleSuccess(stream) {
    var videoTracks = stream.getVideoTracks();
    console.log('Got stream with constraints:', constraints);
    console.log('Using video device: ' + videoTracks[0].label);
    stream.oninactive = function() {
        console.log('Stream inactive');
    };
    window.stream = stream;
    localvideo.srcObject = stream;
}

function handleError(error) {
    if (error.name == 'ConstrainNotSatisfiedError') {
        errorMsg('The resolution ' + constraints.video.exact + 'x' + constraints.video.width.exact + ' px is not supported by your device.')
    }
    else if (error.name == 'PermissionDeniedError') {
        errorMsg('Permission have not been granted to use you camera and ' + 'microphone, you need to allow the page access to your devices in ' + 'order for the demo to work.')
    }
    errorMsg('getUserMedia error: ' + error.name, error);
}

function errorMsg(msg, error) {
    errorElement.innerHTML += '<p>' + msg + '</p>';
    if (typeof error !== 'undefined') {
        // console.error(error);
        console.log(error);
    }
}

function join() {
    navigator.mediaDevices.getUserMedia(constraints).then(handleSuccess).catch(handleError);
}
