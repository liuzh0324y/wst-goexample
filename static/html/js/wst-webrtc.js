

// defaults media constraints
var MEDIA_CONSTRAINTS = {
    audio: false,
    video: {
        width: 640,
        height: 480,
        framerate: 15
    }
}

var ua = (window && window.navigator) ? window.navigator.userAgent : ''

function WstGotUserMediaSuccess(stream) {

}

function WstGotUserMediaFailed(err) {
    
}

function WstWebRtc(mode, options, callback) {
    // if (!(this instanceof WstWebRtc)) {
    //     console.log("WebRtc instance.")
    //     return new WstWebRtc(mode, options, callback)
    // }
    console.log("WebRtc instance.")

    var id = options.id;
    var localVideo = options.localVideo;
    // var remoteVideo = options.remoteVIdeo 

    window.navigator.getUserMedia(MEDIA_CONSTRAINTS, 
        function (stream) { 
            var videoTracks = stream.getVideoTracks();
            // console.log('Got stream with constraints:', constraints);
            console.log('Using video device: ' + videoTracks[0].label);
            stream.oninactive = function() {
                console.log('Stream inactive');
            };
            window.stream = stream;
            localVideo.srcObject = stream;
        }, 
        function (err) {
            console.log(err)
        })
}

var WstWebRtcRecvOnly = function(options, callback) {
    return new WstWebRtc("recvonly", options, callback)
}

var WstWebRtcSendOnly = function(options, callback) {
    return new WstWebRtc("sendonly", options, callback)
}

var WstWebRtcSendRecv = function(options, callback) {
    return new WstWebRtc("sendrecv", options, callback)
}

// function WstWebRtcRecvOnly(options, callback) {
//     if (!(this instanceof WstWebRtcRecvOnly)) {
//         console.log("WebRtcRecv instance.")
//         return new WstWebRtcRecvOnly(options, callback)
//     }
//     // WstWebRtcRecvOnly.super_.call(this, "recvonly", options, callback)
// }

// function WstWebRtcSendOnly(options, callback) {
//     if (!(this instanceof WstWebRtcSendOnly)) {
//         console.log("WebRtcSend instance.")
//         return new WstWebRtcSendOnly(options, callback)
//     }
//     WstWebRtcSendOnly.super_.call(this, "sendonly", options, callback)
// }

// function WstWebRtcSendRecv(options, callback) {
//     if (!(this instanceof WstWebRtcSendRecv)) {
//         console.log("WebRtcSendRecv instance.")
//         return new WstWebRtcSendRecv(options, callback)
//     }
//     WstWebRtcSendRecv.super_.call(this, "sendrecv", options, callback)
// }
