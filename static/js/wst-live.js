
var local = document.getElementById("localvideo");
var joinBtn = document.getElementById("joinBtn");
var callBtn = document.getElementById("callBtn");
var quitBtn = document.getElementById("quitBtn");


local.addEventListener('loadedmetadata', function() {
    console.trace('Local video videowidth: ' + this.videoWidth + 'px, videoHeight: ' + this.videoHeight + 'px');
});

// WstWebRtcSendRecv({localVideo: localVideo}, function(){})

joinBtn.onclick = function () {
    options = {
        id: 1,
        localVideo: local
    }

    WstWebRtcSendRecv(options, function (stream) {
        console.log("callback")
    })
}

callBtn.onclick = function () {

}

quitBtn.onclick = function () {
    
}
