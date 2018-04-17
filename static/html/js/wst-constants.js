// exported Constants

'use strict';

var WstConstants = {
    // Action type for remote web socket communication.
    WS_ACTION: 'ws',
    // Action type for remote xhr communication.
    XHR_ACTION: 'xhr',
    // Action type for adding a command to the remote clean up queue.
    QUEUEADD_ACTION: 'addtoqueue',
    // Action type for clearing the remote clean up queue.
    QUEUECLEAR_ACTION: 'clearqueue',
    // Web socket action type specifying that an event occured.
    EVENT_ACTION: 'event',

    // Web socket action type to create to a remote web socket.
    WS_CREATE_ACTION: 'create',
    // Web socket event type onerror.
    WS_EVENT_ONERROR: 'onerror',
    // Web socket event type onmessage.
    WS_EVENT_ONMESSAGE: 'onmessage',
    // Web socket event type onopen.
    WS_EVENT_ONOPEN: 'onopen',
    // Web socket event type onclose.
    WS_EVENT_ONCLOSE: 'onclose',
    // Web socket event sent when an error occurs while calling send.
    WS_EVENT_SENDERROR: 'onsenderror',
    // Web socket action type to send a message on the remote web socket.
    WS_SEND_ACTION: 'send',
    // Web socket action type to close the remote web socket.
    WS_CLOSE_ACTION: 'close',
};
