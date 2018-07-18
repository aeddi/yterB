// Send message to go then update the GUI
function sendMessage() {
    let input = $(".chat_input").val().trim()
    let datetime = formatDateTime(new Date());
    let id = $(".active_contact").attr("contact_id");
	let payload = {
		Content: input,
		Datetime:  datetime.date + ` ∣ ` + datetime.time,
		Recipient: id,
	}
	let test = JSON.stringify(payload)
	console.log(test)
    $(".chat_input").val('')
    if (input) {
        astilectron.sendMessage({name: "send_message", payload: test}, function(data) {
            console.log("received " + data.payload)
        });
        let data = `
        <div class="sent_content">
            <div class="sent_message">
                <p>` + input + `</p>
                <span class="time_date">` + datetime.time + `&ensp;∣&ensp;` + datetime.date + `</span>
            </div>
        </div>
        `;

        appendMessage({
            id: id,
            data: data,
            datetime: datetime,
            text: input
        });
    	sent_sound.play();
    }
}

// Send username to go
function sendUsernameToGo() {
	var username = $(".username_input").val().trim();
	if (username === "") {
		return;
	} else {
		$(".name_prompt_panel").fadeOut(500);
		astilectron.sendMessage({name: "username", payload: username}, function(data){});
	}
}

// Test function to generate fake data
function outgoingMessageTMP(input, id) {
    let datetime = formatDateTime(new Date());
    let data = `
    <div class="sent_content">
        <div class="sent_message">
            <p>` + input + `</p>
            <span class="time_date">` + datetime.date + `&ensp;∣&ensp;` + datetime.time + `</span>
        </div>
    </div>
    `;

    appendMessage({
        id: id,
        data: data,
        datetime: datetime,
        text: input
    });
	sent_sound.play();
}
