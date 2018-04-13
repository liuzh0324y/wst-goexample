
'use strict';

var localVideo = document.getElementById('localvideo');
var remoteVideo = document.getElementById('remotevideo');
var joinBtn = document.getElementById('joinBtn');
var callBtn = document.getElementById('callBtn');
var quitBtn = document.getElementById('quitBtn');

joinBtn.onclick = join;
callBtn.onclick = call; 
quitBtn.onclick = quit;

var localStream;
var pc1;
var pc2;
var offerOptions = {
    offerToReceiveAudio: 0,
    offerToReceiveVideo: 1
}

var constraints = window.constraints = {
    audio:false,
    video:true
};

localVideo.addEventListener('loadedmetadata', function() {
        console.log('Local video width: ' + this.videoWidth + 'px, height: ' + this.videoHeight + 'px');
});

remoteVideo.addEventListener('loadedmetadata', function(){
    console.log('Remote video width: ' + this.videoWidth + 'px, height: ' + this.videoHeight + 'px');
});

function getName(pc) {
    return (pc === pc1) ? 'pc1' : 'pc2';
}

function getOtherPc(pc) {
    return (pc === pc1) ? pc2 : pc1;
}

function join() {
    navigator.mediaDevices.getUserMedia(constraints).then(WstGotStreamSuccess).catch(WstGotStreamFailed);
    console.log('join btn.');
}

function call() {

    var videoTracks = localStream.getVideoTracks();
    if (videoTracks.length > 0) {
        console.log('Using video device: ' + videoTracks[0].label)
    };

    var servers = null;
    pc1 = new RTCPeerConnection(servers);
    pc1.onicecandidate = function(err) {
        pc1.addIceCandidate(event.candidate).then(
            function() {
                console.log('pc1 addIceCandidate success');
            },
            function(err) {
                console.log('pc1 failed to add ICE Candidate: ' + err.toString());
            }
        );
        console.log('pc1 ICE candidate: \n' + (event.candidate ? event.candidate.candidate : '(null)'));
    };
    
    pc2 = new RTCPeerConnection(servers);
    pc2.onicecandidate = function(err) {
        pc2.addIceCandidate(event.candidate).then(
            function() {
                console.log('pc2 addIceCandidate success');
            },
            function(err) {
                console.log('pc2 failed to add ICE Candidate: ' + err.toString());
            }
        );
        console.log('pc2 ICE candidate: \n' + (event.candidate ? event.candidate.candidate : '(null)'));
    };
    
    pc1.oniceconnectionstatechange = function(err) {
        console.log('pc1 ICE state change event: ', err);
    };
    pc2.oniceconnectionstatechange = function(err) {
        console.log('pc2 ICE state change event: ', err);
    };
    pc2.ontrack = WstGotRemoteStream;

    localStream.getTracks().forEach(function(track) {
        pc1.addTrack(track, localStream);
    });

    pc1.createOffer(offerOptions).then(onCreateOfferSuccess, onCreateSessionDescriptionError);
    console.log('call btn.');
}

function quit() {
    console.log('quit btn.');
}

function WstGotStreamSuccess(stream) {
    
    var tracks = stream.getVideoTracks();
    console.log('Got stream with constraints:', constraints);
    console.info('Using video device: ' + tracks[0].label);
    stream.oninactive = function() {
        console.log('Stream inactive');
    };
    window.stream = stream;
    localVideo.srcObject = stream;
    localStream = stream;
    console.log('get user media device success.');
}

function WstGotStreamFailed(err) {
    console.error(err);
}

function WstGotRemoteStream(err) {
    if (remoteVideo.srcObject !== err.streams[0]) {
        remoteVideo.srcObject = err.streams[0];
        console.log('pc2 received remote stream');
    }
}

function onIceCandidate(pc, event) {
    getOtherPc(pc).addIceCandidate(event.candidate).then(
        function() {
            console.log(getName(pc) + ' addIceCandidate success');
        },
        function(err) {
            console.log(getName(pc) + ' failed to add ICE Candidate: ' + err.toString());
        }
    );
    console.log(getName(pc) + ' ICE candidate: \n' + (event.candidate ? event.candidate.candidate : '(null)'));
}

function onIceStateChange(pc, event) {
    if (pc) {
        console.log('ICE state change event: ', event);
    }
}

function onCreateOfferSuccess(desc) {
    console.log('Offer from pc1\n' + desc.sdp);
    console.log('pc1 setLocalDescription start');
    pc1.setLocalDescription(desc).then(
        function() {
            onSetLocalSuccess(pc1);
            console.log('pc1 setLocalDescription complete.');
        },
        onSetSessionDescriptionError
    );

    pc2.setRemoteDescription(desc).then(
        function() {
            onSetRemoteSuccess(pc2);
        },
        onSetSessionDescriptionError
    );

    pc2.createAnswer().then(
        onCreateAnswerSuccess,
        onCreateSessionDescriptionError
    );
}

function onSetLocalSuccess(pc) {
    console.log('setLocalDescription complete.');
}

function onSetRemoteSuccess(pc) {
    console.log('setRemoteDescription complete.');
}
function onSetSessionDescriptionError(err) {
    console.log('Failed to set session description: ' + err.toString());
}

function onCreateAnswerSuccess(desc) {
    console.log('Answer from pc2:\n' + desc.sdp);
    console.log('pc2 setLocalDescription start');
    pc2.setLocalDescription(desc).then(
        function() {
            onSetLocalSuccess(pc2);
        },
        onSetSessionDescriptionError
    );

    console.log('pc1 setRemoteDescription start');
    pc1.setRemoteDescription(desc).then(
        function() {
            onSetRemoteSuccess(pc1);
        },
        onSetSessionDescriptionError
    );
}

function onCreateSessionDescriptionError(err) {
    console.log('Failed to create session description: ' + err.toString());
}

// var ws = new WebSocket('wss://'+ location.host + '/wst');
// ws.onopen = function(evt) {
//     console.log('Connection open...');
//     ws.send('Hello js websocket.');
// };

// ws.onmessage = function(evt) {
//     console.log('Recevied Message: ' + evt.data);
//     ws.close();
// };

// ws.onclose = function(evt) {
//     console.log('Connection closed.');
// }
