// Bind html elements and javascript
function bind() {
  // Bind send message function
  $( ".send_button" ).click(function() {
    sendMessage();
  });
  $(".chat_input").keyup(function(e) {
    var code = e.keyCode || e.which;
    if(code == 13) { // Enter keycode
      sendMessage();
    }
  });

  // Search bar filter
  $(".search_input").keyup(function(e) {
    var value = $(this).val().toLowerCase();
    $(".contact_list").filter(function() {
      $(this).toggle($(this).find("h5").text().toLowerCase().indexOf(value) > -1)
    });
  });
}

// Bind to switch between chats onclick
function openChat(clicked) {
  if ($(clicked).hasClass("active_contact") === false) {
    let active = $(".active_contact");
    let chat = $(".chat_history");
    let data = contact_data[getContactIndexFromId($(clicked).attr("contact_id"))].data;
    let chat_placeholder = `
      <div class="chat_placeholder">
        <img src="static/img/icon.png"></img>
        <p>Type a message to start a conversation</p>
      </div>
    `;

    $(clicked).addClass("active_contact");
    if(active.length) {
      active.removeClass("active_contact");
    }
    chat.empty();
    if (data === "") {
      chat.append(chat_placeholder);
    }
    else {
      chat.append(data);
    }
    scrollDownChat();
  }
  $(".chat_input").focus();
}
