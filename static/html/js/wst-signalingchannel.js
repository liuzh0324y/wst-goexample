'use strict';

var WstSignalingChannel = function(wssUrl, wssPostUrl) {
    this.wssUrl_ = wssUrl;
    this.wssPostUrl_ = wssPostUrl;
    this.roomId_ = null;
    this.clientId_ = null;
    this.registered_ = false;

    // Public callbacks. Keep it sorted.
    this.onerror = null;
    this.onmessage = null;
};

WstSignalingChannel.prototype.open = function() {
    if (this.websocket_) {
        
    }
};

WstSignalingChannel.prototype.register = function(roomId, clientId) {

};

WstSignalingChannel.prototype.close = function(async) {

};

WstSignalingChannel.prototype.send = function(message) {

};

WstSignalingChannel.prototype.getWssPostUrl = function() {
    return this.wssPostUrl_ + '/' + this.roomId_ + '/' + this.clientId_;
};