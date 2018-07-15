// Format date and time
function formatDateTime(datetime) {
  var formated = {};
  const days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

  formated.time = datetime.getHours() + ":" + (datetime.getMinutes() < 10 ? '0' + datetime.getMinutes() : datetime.getMinutes());
  formated.date = days[datetime.getDay()] + " " + datetime.getDate();

  return formated;
}

// Return matching contact_data index using the ID
function getContactIndexFromId(id) {
  for (i = 0; i < contact_data.length; i++) {
    if (contact_data[i].id == id) {
      return i;
    }
  }
  return -1;
}

// Append message to the chat box then scrolldown
function appendMessage(message) {
  let index = getContactIndexFromId(message.id);

  $(message.data).appendTo(".chat_history");
  scrollDownChat();

  contact_data[index].last_datetime = message.datetime;
  contact_data[index].data = contact_data[index].data + message.data;
  updateContact(message);
}

// Scroll chat div content down
function scrollDownChat() {
  let chat = $(".chat_history");

  chat.animate({scrollTop: chat[0].scrollHeight}, 500);
}

// Update contact card in contact list
function updateContact(data) {
  let updated_contact = $("div[contact_id='" + data.id + "']");
  updated_contact.find(".chat_date").text(data.datetime.time);
  updated_contact.find("p").text(data.text);
}
