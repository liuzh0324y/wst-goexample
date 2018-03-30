'use strict';

var startButton = document.getElementById('startButton');
var callButton = document.getElementById('callButton');
var hangupButton = document.getElementById('hangupButton');

callButton.disabled = true;
hangupButton.disabled = true;
startButton.onclick = start;
callButton.onclick = call;
hangupButton.onclick = hangup;

var startTime;
var localVideo = document.getElementById('localVideo');
var remoteVideo = document.getElementById('remoteVideo');

localVideo.addEventListener('loadedmetadata', function() {
    console.trace('Local video videowidth: ' + this.videoWidth + 'px, videoHeight: ' + this.videoHeight + 'px');
});

remoteVideo.addEventListener('loadedmetadata', function() {
    console.trace('Remote video videoWidth: ' + this.videoWidth + 'px, videoHeight: ' + this.videoHeight + 'px');
});

remoteVideo.onresize = function() {
    console.trace('Remote video size changed to ' + remoteVideo.videoWidth + 'x' + remoteVideo.videoHeight);
    if (startTime) {
        var elapsedTime = window.performance.now() - startTime;
        console.trace('Setup time: ' + elapsedTime.toFixed(3) + 'ms');
        startTime = null;
    }
};

var localStream;
var pc1;
var pc2;
var offerOptions = {
    offerToReceiveAudio: 1,
    offerToreceiveVideo: 1
};

function getName(pc) {
    return (pc === pc1) ? 'pc1' : 'pc2';
}

function getOtherPc(pc) {
    return (pc === pc1) ? pc2 : pc1;
}

function gotStream(stream) {
    console.trace('Received local stream');
    localVideo.srcObject = stream;
    localStream = stream;
    callButton.disabled = false;
}

function start() {
    console.trace('Requesting local stream');
    startButton.disabled = true;
    navigator.mediaDevices.getUserMedia({
        audio: true,
        video: true
    })
    .then(gotStream)
    .catch(function(e) {
        alert('getUserMedia() error: ' + e.name);
    });
}

function call() {
    callButton.disabled = true;
    hangupButton.disabled = false;
    console.trace('Starting call');
    startTime = window.performance.now();
    var videoTracks = localStream.getVideoTracks();
    var audioTracks = localStream.getAudioTracks();
    if (videoTracks.length > 0) {
        console.trace('Using video device: ' + videoTracks[0].label);
    }
    if (audioTracks.length > 0) {
        console.trace('Using audio device: ' + audioTracks[0].label);
    }
    
    var servers = null;
    
    pc1 = new RTCPeerConnection(servers);
    console.trace('Created local peer connection object pc1');
    pc1.onicecandidate = function(e) {
        onIceCandidate(pc1, e);
    };
    
    pc2 = new RTCPeerConnection(servers);
    console.trace('Create remote peer connection object pc2');
    pc2.onicecandidate = function(e) {
        onIceCandidate(pc2, e);
    };

    pc1.oniceconnectionstatechange = function(e) {
        onIceStateChange(pc1, e);
    }
    pc2.oniceconnectionstatechange = function(e) {
        onIceStateChange(pc2, e);
    }
    pc2.ontrack = gotRemoteStream;

    localStream.getTracks().forEach(
        function(track) {
            pc1.addTrack(
                track,
                localStream
            );
        }
    );
    console.trace('Added local stream to pc1');

    console.trace('pc1 createOffer start');
    pc1.createOffer(
        offerOptions
    ).then(
        onCreateOfferSuccess,
        onCreateSessionDescriptionError
    );
}

function onCreateSessionDescriptionError(error) {
    console.trace('Failed to create session description: ' + error.toString());
}

function onCreateOfferSuccess(desc) {
    console.trace('Offer from pc1\n' + desc.sdp);
    console.trace('pc1 setLocalDescription start');
    pc1.setLocalDescription(desc).then(
        function() {
            onSetLocalSuccess(pc1);
        },
        onSetSessionDescriptionError
    );
    console.trace('pc2 setRemoteDescription start');
    pc2.setRemoteDescription(desc).then(
        function() {
            onSetRemoteSuccess(pc2);
        },
        onSetSessionDescriptionError
    );
    console.trace('pc2 createAnswer start');
    pc2.createAnswer().then(
        onCreateAnswerSuccess,
        onCreateSessionDescriptionError
    );
}

function onSetLocalSuccess(pc) {
    console.trace(getName(pc) + ' setLocalDescription complete');
}

function onSetRemoteSuccess(pc) {
    console.trace(getName(pc) + ' setRemoteDescription complete');
}

function onSetSessionDescriptionError(error) {
    console.trace('Failed to set session description: ' + error.toString());
}

function gotRemoteStream(e) {
    if (remoteVideo.srcObject !== e.streams[0]) {
        remoteVideo.srcObject = e.streams[0];
        console.trace('pc2 received remote stream');
    }
}

function onCreateAnswerSuccess(desc) {
    console.trace('Answer from pc2:\n' + desc.sdp);
    console.trace('pc2 setLocalDescription start');
    pc2.setLocalDescription(desc).then(
        function() {
            onSetLocalSuccess(pc2);
        },
        onSetSessionDescriptionError
    );
    console.trace('pc1 setRemoteDescription start');
    pc1.setRemoteDescription(desc).then(
        function() {
            onAddIceCandidateSuccess(pc1);
        },
        function(err) {
            onAddIceCandidateError(pc, err);
        }
    );
}

function onIceCandidate(pc, event) {
    getOtherPc(pc).addIceCandidate(event.candidate)
    .then(
        function() {
            onAddIceCandidateSuccess(pc);
        },
        function(err) {
            onAddIceCandidateError(pc, err);
        }
    );
    console.trace(getName(pc) + ' ICE candidate: \n' + (event.candidate ? event.candidate.candidate : '(null)'));
}

function onAddIceCandidateSuccess(pc) {
    console.trace(getName(pc) + ' addIceCandidate success');
}

function onAddIceCandidateError(pc, error) {
    console.trace(getName(pc) + ' failed to add ICE Candidate: ' + error.toString());
}

function onIceStateChange(pc, event) {
    if (pc) {
        console.trace(getName(pc) + ' ICE state: ' + pc.iceConnectionState);
        console.log('ICE state change event: ', event);
    }
}

function hangup() {
    console.trace('Ending call');
    pc1.close();
    pc2.close();
    pc1 = null;
    pc2 = null;
    hangupButton.disabled = true;
    callButton.disabled = true;
}