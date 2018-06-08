window.ws = new WebSocket("ws://" + window.location.host + "/ws");

ws.onopen = function() {
  console.log('connection live');
}

ws.onmessage = function(event) {
  var message = JSON.parse(event.data);
  console.log("received message:", message);

  if(message.code === 'registered.success') {
    console.log('received register.success message from server');
    // update the WebsocketId locally
    window.userInfo.WebsocketId = message.websocketId;
    return;
  }

  // this condition means someone has initiated a DM. subscribe this user to
  // the channel with pubnub and update the channel with messages as they come in
  if(message.code === 'subscribe.direct') {
    console.log('received dm request from server. Subscribing to dm request');

    var channelsEl = PUBNUB.$('channels');
    var channelId = message.channel;

    // create the new channel box
    var el = document.createElement('div');
    el.className = 'channel-container';
    var elString = '<div class="channel-chat-box" id=chat-box-'+channelId+'><div class="channel-title">Direct message: '+channelId+'</div></div><div class="channel-input">Chat: <input id="'+channelId+'" size="50" placeholder="Say something"/></div>';
    el.innerHTML = elString;
    channelsEl.appendChild(el);

    window.pubnub.subscribe({
      channel: channelId,
      callback: function(message) {
        var chatBoxEl = PUBNUB.$('chat-box-'+channelId);
        var el = document.createElement('div');
        el.className = 'message';
        if(message.userId === window.userInfo.UserId) {
          el.className = 'message-mine';
        }
        var username = message.username;
        el.innerHTML = '<span>'+username+': '+message.value+'</span>';
        chatBoxEl.appendChild(el);
      }
    });

    registerInputs();
  }
}

ws.onerror = function(event) {
  console.log('error in ws client', event);
}

ws.onclose = function(event) {
  console.log('ws closing', event);
}
