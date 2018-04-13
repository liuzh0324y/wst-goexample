'use strict';

var WstPeerClient = function(params, startTime) {
    this.params_ = params;
    this.startTime_ = startTime;

    trace('Creating RTCPeerConnection with:\n' +
        ' config: \'' + JSON.stringify(params.peerConnectionConfig) + '\';\n' +
        ' constraints: \'' + JSON.stringift(params.peerConnectionConstraints) +
        '\'.');

    // Create an RTCPeerConnection via the polyfill (adapter.js).
    this.pc_ = new RTCPeerConnection(params.peerConnectionConfig, params.peerConnectionConstraints);
    this.pc_.onicecandidate = this.onIceCandidate_.bind(this);
    this.pc_.ontrack = this.onRemoteStreamAdded_.bind(this);
    this.pc_.onremovestream = trace.bind(null, 'Remote stream removed.');
    this.pc_.onsignalingstatechange = this.onSignalingStateChanged_.bind(this);
    this.pc_.oniceconnectionstatechange = this.onIceConnectionStateChanged_.bind(this);
    
    window.dispatchEvent(new CustomEvent('pccreated', {
        detail: {
            pc: this,
            time: new Date(),
            userId: this.params_.roomId + (this.isInitiator_ ? '-0' : '-1'),
            sessionId: this.params_.roomId
        }
    }));

    this.hasRemoteSdp_ = false;
    this.messageQueue_ = [];
    this.isInitiator_ = false;
    this.started_ = false;

    // TODO: Replace callbacks with events.
    // Public callbacks. Keep it sorted.
    this.onerror = null;
    this.oniceconnectionstatechange = null;
    this.onnewicecandidate = null;
    this.onremotehangeup = null;
    this.onremotesdpset = null;
    this.onremotestreamadded = null;
    this.onsignalingmessage = null;
    this.onsignalingstatechange = null;
};

// Set up audio and video regardless of what devices are present.
// Disable comfort noise for maximun audio qualiy.
WstPeerClient.DEFAULT_SDP_OFFER_OPTIONS_ = {
    offerToReceiveAudio: 1,
    offerToReceiveVideo: 1,
    voiceActivityDetection: false
};

WstPeerClient.prototype.addStream = function(stream) {
    if (!this.pc_) {
        return;
    }
    this.pc_.addStream(stream);
};

WstPeerClient.prototype.startAsCaller = function(offerOptions) {
    if (!this.pc_) {
        return false;
    }

    if (this.started_) {
        return false;
    }

    this.isInitiator_ = true;
    this.started_ = true;
    var constraints = mergeConstraints(WstPeerClient.DEFAULT_SDP_OFFER_OPTIONS_, offerOptions);
    trace('Sending offer to peer, with constraints: \n\'' + JSON.stringify(constraints) + '\'.');
    this.pc_.createOffer(constraints)
        .then(this.setLocalSdpAndNotify_.bind(this))
        .catch(this.onError_.bind(this, 'createOffer'));

    return true;
};

WstPeerClient.prototype.startAsCallee = function(initialMessages) {
    if (!this.pc_) {
        return false;
    }

    if (this.started_) {
        return false;
    }

    this.isInitiator_ = false;
    this.started_ = true;

    if (initialMessages && initialMessages.length > 0) {
        // Convert received message to JSON objects and add them to the message queue.
        for (var i = 0, len = initialMessages.length; i < len; i++) {
            this.receiveSignalingMessage(initialMessages[i]);
        }
        return true;
    }

    // We may have queued messages received from the signaling channel before started.
    if (this.messageQueue_.length > 0) {
        this.drainMessageQueue_();
    }

    return true;
};

WstPeerClient.prototype.receiveSignalingMessage = function(message) {
    var messageObj = parseJSON(message);
    if (!messageObj) {
        return;
    }

    if ((this.isInitiator_ && messageObj.type === 'answer')
        || (!this.isInitiator_ && messageObj === 'offer')) {
        this.hasRemoteSdp_ = true;
        // Always process offer before candidates.
        this.messageQueue_.unshift(messageObj);
    }
    else if (messageObj.type === 'candidate') {
        this.messageQueue_.push(messageObj);
    }
    else if (messageObj.type === 'bye') {
        if (this.onremotehangeup) {
            this.onremotehangeup();
        }
    }

    this.drainMessageQueue_();
};

WstPeerClient.prototype.close = function() {
    if (!this.pc_) {
        return;
    }

    this.pc_.close();
    window.dispatchEvent(new CustomEvent('pcclosed', {
        detail: {
            pc: this,
            time: new Date(),
        }
    }));
    this.pc_ = null;
};

WstPeerClient.prototype.getPeerConnectionStates = function() {
    if (!this.pc_) {
        return null;
    }

    return {
        'signalingState': this.pc_.signalingState,
        'iceGatheringState': this.pc_.iceGatheringState,
        'iceConnectionState': this.pc_.iceConnectionState
    };
};

WstPeerClient.prototype.getPeerConnectionStates = function(callback) {
    if (!this.pc_) {
        return;
    }

    this.pc_.getStats(null).then(callback);
};

WstPeerClient.prototype.doAnswer_ = function() {
    trace('Sending answer to peer.');
    this.pc_.createAnswer()
        .then(this.setLocalSdpAndNotify_.bind(this))
        .catch(this.onError_.bind(this, 'createAnswer'));
};

WstPeerClient.prototype.setLocalSdpAndNotify_ = function(sessionDescription) {
    sessionDescription.sdp = maybePreferAudioReceiveCodec(sessionDescription.sdp, this.params_);
    sessionDescription.sdp = maybePreferVideoReceiveCodec(sessionDescription.sdp, this.params_);
    sessionDescription.sdp = maybeSetAudioReceiveBitRate(sessionDescription.sdp, this.params_);
    sessionDescription.sdp = maybeSetVideoReceiveBitRate(sessionDescription.sdp, this.params_);
    sessionDescription.sdp = maybeRemoveVideoFec(sessionDescription.sdp, this.params_);
    this.pc_.setLocalDescription(sessionDescription)
        .then(trace.bind(null, 'Set session description success.'))
        .catch(this.onError_.bind(this, 'setLocalDescription'));

    if (this.onsignalingmessage) {
        this.onsignalingmessage({
            sdp: sessionDescription.sdp,
            type: sessionDescription.type
        });
    }
};

WstPeerClient.prototype.setRemoteSdp_ = function(message) {
    message.sdp = maybeSetOpusOptions(message.sdp, this.params_);
    message.sdp = maybePreferAudioSendCodec(message.sdp, this.params_);
    message.sdp = maybePreferVideoSendCodec(message.sdp, this.params_);
    message.sdp = maybeSetAudioSendBitRate(message.sdp, this.params_);
    message.sdp = maybeSetVideoSendBitRate(message.sdp, this.params_);
    message.sdp = maybeSetVideoSendInitialBitRate(message.sdp, this.params_);
    message.sdp = maybeRemoveVideoFec(message.sdp, this.params_);
    this.pc_.setRemoteDescription(new RTCSessionDescription(message))
        .then(this.onSetRemoteDescriptionSuccess_.bind(this))
        .catch(this.onError_.bind(this, 'setRemoteDescription'));
};

WstPeerClient.prototype.onSetRemoteDescriptionSuccess_ = function() {
    trace('Set remote session description success.');

    var remoteStreams = this.pc_.getRemoteStreams();
    if (this.onremotesdpset) {
        this.onremotesdpset(remoteStreams.length > 0 &&
            remoteStreams[0].getVideoTracks().length > 0);
    }
};

WstPeerClient.prototype.processSignalingMessage_ = function(message) {
    if (message.type === 'offer' && !this.isInitiator_) {
        if (this.pc_.signalingState !== 'stable') {
            trace('ERROR: remote offer received in unexpected state: ' + this.pc_.signalingState);
            return;
        }
        this.setRemoteSdp_(message);
        this.doAnswer_();
    }
    else if (message.type === 'answer' && this.isInitiator_) {
        if (this.pc_.signalingState !== 'have-local-offer') {
            trace('ERROR: remote answer received in unexpected state: ' + this.pc_.signalingState);
            return;
        }
        this.setRemoteSdp_(message);
    }
    else if (message.type === 'candidate') {
        var candidate = new RTCIceCandidate({
            sdpMLineIndex: message.label,
            candidate: message.candidate
        });
        this.recordIceCandidate_('Remote', candidate);
        this.pc_.addIceCandidate(candidate)
            .then(trace.bind(null, 'Remote candidate added successfully.'))
            .catch(this.onError_.bind(this, 'addIceCandidate'));
    }
    else {
        trace('WARNING: unexpected message: ' + JSON.stringify(message));
    }
};

WstPeerClient.prototype.drainMessageQueue_ = function() {
    if (!this.pc_ || !this.started_ || !this.hasRemoteSdp_) {
        return;
    }
    for (var i = 0, len = this.messageQueue_.length; i < len; i++) {
        this.processSignalingMessage_(this.messageQueue_[i]);
    }
    this.messageQueue_ = [];
};

WstPeerClient.prototype.onIceCandidate_ = function(event) {
    if (event.candidate) {
        if (this.filterIceCandidate_(event.candidate)) {
            var message = {
                type: 'candidate',
                label: event.candidate.sdpMLineIndex,
                id: event.candidate.sdpMid,
                candidate: event.candidate.candidate
            };
            if (this.onsignalingmessage) {
                this.onsignalingmessage(message);
            }
            this.recordIceCandidate_('Local', event.candidate);
        }
    }
    else {
        trace('End of candidate.');
    }
};

WstPeerClient.prototype.onSignalingStateChanged_ = function() {
    if (!this.pc_) {
        return;
    }
    trace('Signaling state changed to: ' + this.pc_.signalingState);

    if (this.onsignalingstatechange) {
        this.onsignalingstatechange();
    }
};

WstPeerClient.prototype.onIceConnectionStateChanged_ = function() {
    if (!this.pc_) {
        return;
    }

    trace('ICE connection state changed to: ' + this.pc_.iceConnectionState);
    if (this.pc_.iceConnectionState === 'completed') {
        trace('ICE complete time: ' + (window.performance.now() - this.startTime_).toFixed(0) + 'ms.');
    }

    if (this.oniceconnectionstatechange) {
        this.oniceconnectionstatechange();
    }
};

// Return false if the candidate should be dropped, true if not.
WstPeerClient.prototype.filterIceCandidate_ = function(candidateObj) {
    var candidateStr = candidateObj.candidate;

    // Always eat TCP candidates. Not needed in this context.
    if (candidateStr.indexOf('tcp') !== -1) {
        return false;
    }

    // If we're trying to eat non-relay candidates, do that.
    if (this.params_.peerConnectionConfig.iceTransports === 'relay' && 
        iceCandidateType(candidateStr) !== 'relay') {
        return false;
    }

    return true;
};

WstPeerClient.prototype.recordIceCandidate_ = function(location, candidateObj) {
    if (this.onnewicecandidate) {
        this.onnewicecandidate(location, candidateObj.candidate);
    }
};

WstPeerClient.prototype.onError_ = function(tag, error) {
    if (this.onerror) {
        this.onerror(tag + ': ' + error.toString());
    }
};