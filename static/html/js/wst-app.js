'use strict';

// Keep this in sync with the HTML element id attributes. Keep it sorted.
var UI_CONSTANTS = {
    localVideo: '#local-video',
    remoteVideo: '#remote-video',
    joinButton: '#join-button',
    callButton: '#call-button',
    quitButton: '#quit-button',
}

var loadingParams = {
    roomServer: ':8090',
    roomId: 'R1001',
    clientId: 'C1001',
    wssUrl: 'wss://10.33.48.14:8090/wst',
    mediaConstraints: {
        audio: false,
        video: true
    }
}

var WstApp = function(params) {
    trace('Initializing; server=' + params.roomServer + '.');
    trace('Initializing; room=' + params.roomId + '.');
    trace('Initializing; client=' + params.clientId + '.');

    this.localVideo_ = $(UI_CONSTANTS.localVideo);
    this.remoteVideo_ = $(UI_CONSTANTS.remoteVideo);
    this.joinButton_ = $(UI_CONSTANTS.joinButton);
    this.callButton_ = $(UI_CONSTANTS.callButton);
    this.quitButton_ = $(UI_CONSTANTS.quitButton);

    this.joinButton_.addEventListener('click', this.onJoinClick_.bind(this), false);
    this.callButton_.addEventListener('click', this.onCallClick_.bind(this), false);
    this.quitButton_.addEventListener('click', this.onQuitClick_.bind(this), false);

    this.loadingParams_ = params;
    this.loadUrlParams_();

    this.localStream_ = null;

}

WstApp.prototype.loadUrlParams_ = function() {
    var DEFAULT_VIDEO_CODEC = 'VP9';
    // var urlParams = queryStringToDictinary(window.location.search);
}

WstApp.prototype.onJoinClick_ = function() {
    trace('click join button.');
    if (this.loadingParams_.roomId) {

        this.createCall_();

        this.finishCallSetup_(this.loadingParams_.roomId);
    }
    
};

WstApp.prototype.onCallClick_ = function() {
    trace('click call button.');
    this.call_.sendSignalingMessage_('hello');
};

WstApp.prototype.onQuitClick_ = function() {
    trace('click quit button.');
};

WstApp.prototype.createCall_ = function() {
    trace('create call');

    this.call_ = new WstCall(this.loadingParams_);

    this.call_.onremotehangup = this.onRemoteHangup_.bind(this);
    this.call_.onremotesdpset = this.onRemoteSdpSet_.bind(this);
    this.call_.onremotestreamadded = this.onRemoteStreamAdded_.bind(this);
    this.call_.onlocalstreamadded = this.onLocalStreamAdded_.bind(this);

    this.call_.onerror = this.displayError_.bind(this);
    this.call_.onstatusmessage = this.displayStatus_.bind(this);
    this.call_.oncallerstarted = this.displaySharingInfo_.bind(this);
};

WstApp.prototype.finishCallSetup_ = function(roomId) {
    this.call_.start(roomId);
}
WstApp.prototype.onRemoteHangup_ = function() {
    this.displayStatus_('The remote side hangup up.');
    // Stop waiting for remote video.
    this.remoteVideo_.oncanplay = undefined;

    

    this.call_.onRemoteHangup();
};

WstApp.prototype.onLocalStreamAdded_ = function(stream) {
    trace(' on local stream added.');
    this.localStream_ = stream;

    this.localVideo_.srcObject = this.localStream_;
}

WstApp.prototype.onRemoteSdpSet_ = function(hasRemoteVideo) {
    if (hasRemoteVideo) {
        trace('Waiting for remote video.');
        this.waitForRemoteVideo_();
    }
};

WstApp.prototype.onRemoteStreamAdded_ = function(stream) {
    trace('on remote stream added.');
}


WstApp.prototype.displayError_ = function(error) {
    trace(error);
}

WstApp.prototype.displayStatus_ = function(status) {

}

WstApp.prototype.displaySharingInfo_ = function(roomId, roomLink) {

}