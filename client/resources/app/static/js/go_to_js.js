// Append contact to contact list
function addContact(contact) {
	console.log(contact)
	contact = JSON.parse(contact)
	console.log(contact)
    let new_contact = `
    <div class="contact" contact_id="` + contact.Peer_id + `" onclick="openChat(this)">
        <div class="contact_box">
            <div class="contact_avatar"><img src="https://api.adorable.io/avatars/285/` + contact.Peer_id + `.png"></div>
            <div class="contact_infos">
                <h5>` + contact.Name + `<span class="chat_date"><span style="font-weight: 700">New</span></span></h5>
                <p>Click to start a conversation</p>
				<div class="unread_dot"></div>
            </div>
        </div>
    </div>`;
    $(new_contact).appendTo(".contact_list");
    setTimeout(function() {
        $("div[contact_id='" + contact.Peer_id + "']").last().addClass("appended");
    }, 300);

    // If contact does'nt already exist, add it to the list
    if (getContactIndexFromId(contact.Peer_id) === -1) {
        contact_data.push({
            id: contact.Peer_id,
            name: contact.Name,
            data: ""
        });
    }
}

// Remove contact from contact list
function removeContact(contact) {
    let del_contact = $("div[contact_id='" + contact.id + "']");
	let is_active = del_contact.hasClass("active_contact");

    del_contact.css("z-index", 1);
    del_contact.css("margin-top", "-" + (del_contact.height() + 35) + "px");
    setTimeout(function() {
        $("div[contact_id='" + contact.id + "']").remove();
    }, 800);

	if ($(".contact").length == 0 || is_active) {
		noContactPlaceHolder();
	}
}

// Create incoming message element
function incomingMessage(message)
{
    //
    let datetime = formatDateTime(new Date());
    //
    let data = `
    <div class="received_avatar"><img src="https://api.adorable.io/avatars/285/` + message.sender + `.png"></div>
    <div class="received_content">
        <div class="received_message">
            <p>` + message.content + `</p>
            <span class="time_date">` + datetime.date + `&ensp;∣&ensp;` + datetime.time + `</span>
        </div>
    </div>
    `;

    appendMessage({
        id: message.sender,
        data: data,
        datetime: datetime,
        text: message.content
    });
	recv_sound.play();
}
