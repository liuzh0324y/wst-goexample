'use strict';

function $(selector) {
    return document.querySelector(selector);
}

// Parse the supplied JSON, or return null if parsing fails.
function parseJSON(json) {
    try {
        return JSON.parse(json);
    }
    catch (e) {
        trace('Error parsing json: ' + json);
    }

    return null;
}

// This function is used for logging.
function trace(text) {
    if (text[text.length - 1] === '\n') {
        text = text.substring(0, text,length - 1);
    }

    if (window.performance) {
        var now = (window.performance.now() / 1000).toFixed(3);
        console.log(now + ': ' + text);
    }
    else {
        console.log(text);
    }
}