'use strict';

var WstCall = function(params) {
    this.params_ = params;
    this.roomServer_ = params.roomServer || '';

    this.channel_ = new WstSignalingChannel(params.wssUrl, params.wssPostUrl);
    this.channel_.onmessage = this.onRecvSignalingChannelMessage_.bind(this);

    this.pcClient_ = null;
    this.localStream_ = null;
    this.errorMessageQueue_ = [];
    this.startTime = null;

    // Public callbacks. Keep it sorted.
    this.oncallerstarted = null;
    this.onerror = null;
    this.oniceconnectionstatechange = null;
    this.onlocalstreamadded = null;
    this.onnewicecandidate = null;
    this.onremotehangup = null;
    this.onremotesdpset = null;
    this.onremotestreamadded = null;
    this.onsignalingstatechange = null;
    this.onstatusmessage = null;

    this.getMediaPromise_ = null;
    this.getIceServersPromise_ = null;
    this.requestMediaAndIceServers_();
};

WstCall.prototype.requestMediaAndIceServers_ = function() {
    this.getMediaPromise_ = this.maybeGetMedia_();
    this.getIceServersPromise_ = this.maybeGetIceServers_();
    trace('request media and ice servers.')
};

WstCall.prototype.isInitiator = function() {
    return this.params_.isInitiator;
};

WstCall.prototype.start = function(roomId) {
    this.connectToRoom_(roomId);
    if (this.params_.isLoopback) {
        setupLoopback(this.params_.wssUrl, roomId);
    }
};

WstCall.prototype.onRemoteHangup = function() {

};

// Asynchronously request user media if needed.
WstCall.prototype.maybeGetMedia_ = function() {
    // mediaConstraints.audio and mediaConstraints.video could be objects, so
    // check '!=== false' instead of '=== true'.
    var needStream = (this.params_.mediaConstraints.audio !== false ||
                      this.params_.mediaConstraints.video !== false);
    var mediaPromise = null;
    if (needStream) {
        var mediaConstraints = this.params_.mediaConstraints;

        mediaPromise = navigator.mediaDevices.getUserMedia(mediaConstraints)
            .catch(function(error) {
                if (error.name !== 'NotFoundError') {
                    throw error;
                }
                return navigator.mediaDevices.enumerateDevices()
                    .then(function(devices) {
                        var cam = devices.find(function(device) {
                            return device.kind === 'videoinput';
                        });
                        var mic = devices.find(function(device) {
                            return device.kind === 'audioinput';
                        });
                        var constraints = {
                            video: cam && mediaConstraints.video,
                            audio: mic && mediaConstraints.audio
                        };
                        return navigator.mediaDevices.getUserMedia(constraints);
                    });
            })
            .then(function(stream) {
                trace('Got access to local media with mediaConstraints:\n' +
                    ' \'' + JSON.stringify(mediaConstraints) + '\'');
                this.onUserMediaSuccess_(stream);
            }.bind(this))
            .catch(function(error){
                this.onError_('Error getting user media: ' + error.message);
                this.onUserMediaError_(error);
            }.bind(this));
    }
    else {
        mediaPromise = Promise.resolve();
    }
    return mediaPromise;
};

// Asynchronously request an ICE server if needed.
WstCall.prototype.maybeGetIceServers_ = function() {
    var shouldRequestIceServers = 
        (this.params_.iceServerRequestUrl &&
         this.params_.iceServerRequestUrl.length > 0 &&
         this.params_.peerConnectionConfig.iceServers &&
         this.params_.peerConnectionConfig.iceServers.length === 0);
    
    // var iceServerPromise = requestIceServers()
};

WstCall.prototype.onUserMediaSuccess_ = function(stream) {
    this.localStream_ = stream;
    if (this.onlocalstreamadded) {
        this.onlocalstreamadded(stream);
    }
};

WstCall.prototype.onUserMediaError_ = function(error) {
    var errorMessage = 'Failed to get access to local media. Error name was ' +
        error.name + '. Continuing without sending a stream.';
    this.onError_('getUserMedia error: ' + errorMessage);
    this.errorMessageQueue_.push(error);
    alert(errorMessage);
};


WstCall.prototype.onRecvSignalingChannelMessage_ = function(msg) {

};