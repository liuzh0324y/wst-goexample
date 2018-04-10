'use strict';

var leftVideo = document.getElementById('leftVideo')
var rightVideo = document.getElementById('rightVideo')

var stream;

function maybeCreateStream() {
    if (stream) {
        return;
    }
    if (leftVideo.captureStream) {
        stream = leftVideo.captureStream();
        rightVIdeo.srcObject = stream;
        console.log('Captured stream from leftVideo with captureStream', stream);
    }
    else if (leftVideo.mozCaptureStream) {
        stream = leftVideo.mozCaptureStream();
        rightVideo.srcObject = stream;
        console.log('Captured stream from leftVideo with mozCaptureStream()', stream);
    }
    else {
        console.log('captureStream() not supported');
    }
}

leftVideo.oncanplay = maybeCreateStream;
if (leftVideo.readyState >= 3) {
    maybeCreateStream();
}

leftVideo.play();