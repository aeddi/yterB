var contact_data = [];
var sent_sound = new Audio('static/sound/blop.wav');
var recv_sound = new Audio('static/sound/pling.wav');

let init = {
	about: function(html) {
		let c = document.createElement("div");
		c.innerHTML = html;
		asticode.modaler.setContent(c);
		asticode.modaler.show();
	},
	detach: function() {
		// Detach dev tools
		try {
			mainWindow.devToolsWebContents.executeJavaScript('document.getElementsByClassName("long-click-glyph")[0].click()');
		} catch(e) {}
	},
	init: function() {
		// Init libs
		asticode.loader.init();
		asticode.modaler.init();
		asticode.notifier.init();

		// Wait for astilectron to be ready
		document.addEventListener('astilectron-ready', function() {
			init.listen();
			bind();
			$(function() {
				noContactPlaceHolder();
				init.fakedata();
			});
		})
	},
	fakedata: function() {
		setTimeout(function() {
			addContact({id: 1, name: "Jessica Sarti"});
			addContact({id: 2, name: "Jeremie Gazaniol"});
			addContact({id: 3, name: "Enguerrand Lepretre"});
			addContact({id: 4, name: "Carole Barreau"});
			addContact({id: 5, name: "Jeremy Salmeron"});
			addContact({id: 6, name: "John Doe"});
			setTimeout(function() {removeContact({id: 5})}, 5000);

			outgoingMessageTMP("Test, which is a new approach to have", 3);
			incomingMessage({sender: 3, content: "Apollo University, Delhi, India Test"});
			outgoingMessageTMP("We work directly with our designers and suppliers, and sell direct to you, which means quality, exclusive products, at a price anyone can afford.", 3);
			incomingMessage({sender: 3, content: "Test which is a new approach to have all solutions"});
			outgoingMessageTMP("Test, which is a new approach to have", 3);
			incomingMessage({sender: 3, content: "Test, which is a new approach to have"});
			outgoingMessageTMP("Apollo University, Delhi, India Test", 3);
			incomingMessage({sender: 3, content: "We work directly with our designers and suppliers, and sell direct to you, which means quality, exclusive products, at a price anyone can afford."});
			outgoingMessageTMP("Test which is a new approach to have all solutions", 3);
			incomingMessage({sender: 3, content: "Test, which is a new approach to have"});
		}, 1000);

	},
	listen: function() {
		astilectron.onMessage(function(message) {
			switch (message.name) {
				case "about":
				init.about(message.payload);
				return {payload: "payload"};
				break;
				case "detach":
				init.detach();
				return {payload: "payload"};
				break;
			}
		});
	}
};
