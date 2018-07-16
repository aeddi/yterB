// Send message to go then update the GUI
function sendMessage() {
    let input = $(".chat_input").val().trim()
    let datetime = formatDateTime(new Date());
    let id = $(".active_contact").attr("contact_id");

    $(".chat_input").val('')
    if (input) {
        astilectron.sendMessage({name: "send_message", payload: input}, function(message) {
            console.log("received " + message.payload)
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

// TMP FUNCTION
function outgoingMessageTMP(input, id)
{
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
