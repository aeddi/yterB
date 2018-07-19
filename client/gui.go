package main

import (
	"encoding/json"
	"flag"
	"strconv"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Constants
const html_about = `Welcome on <b>ytreB</b>!<br>
Developed by aeddi using libp2p, crypto and astilectron.`

var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	a       *astilectron.Astilectron
)

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icons/icon.icns",
			AppIconDefaultPath: "resources/icons/icon.png",
		},
		Debug: *debug,
		MenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astilectron.PtrStr("Menu"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label: astilectron.PtrStr("About"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "about", html_about, func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
									return
								}
								astilog.Infof("About modal has been displayed and payload is %s!", s)
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending about event failed"))
							}
							return
						},
					},
					{Type: astilectron.MenuItemTypeSeparator},
					{
						Label: astilectron.PtrStr("Debug"),
						Role:  astilectron.MenuItemRoleToggleDevTools,
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "detach", "", func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
									return
								}
								astilog.Infof("About modal has been displayed and payload is %s!", s)
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending detach event failed"))
							}
							return
						},
					},
					{Role: astilectron.MenuItemRoleClose},
					{Role: astilectron.MenuItemRoleQuit},
				},
			},
			{
				Label: astilectron.PtrStr("Edit"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Role: astilectron.MenuItemRoleCopy},
					{Role: astilectron.MenuItemRoleCut},
					{Role: astilectron.MenuItemRolePaste},
					{Role: astilectron.MenuItemRoleSelectAll},
				},
			},
			{Role: astilectron.MenuItemRoleAbout},
		},
		OnWait: func(aa *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			a = aa
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Center:    astilectron.PtrBool(true),
				Height:    astilectron.PtrInt(600),
				Width:     astilectron.PtrInt(1100),
				MinHeight: astilectron.PtrInt(400),
				MinWidth:  astilectron.PtrInt(800),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// Handle messages sent from JS
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {

	switch m.Name {
	case "send_message":
		var decode Message
		unquoted, _ := strconv.Unquote(string(m.Payload))
		err := json.Unmarshal([]byte(unquoted), &decode)
		if err != nil {
			payload = err.Error()
		} else {
			requestAuthcode(decode)
			payload = decode.Content
		}
	case "username":
		initClient(unmarshalPayload(m))
	}
	return
}

// Marshal payload then send command
func sendCommandToJS(header string, command interface{}) {

	encoded, _ := json.Marshal(command)
	bootstrap.SendMessage(w, header, string(encoded), func(m *bootstrap.MessageIn) { _ = m })
}

// Send command logs to javascript console
func consoleLog(log interface{}) {

	sendCommandToJS("debug", log)
}

func addContactToGUI(client Client) {

	sendCommandToJS("add_contact", client)
}

// Unmarshal payload
func unmarshalPayload(m bootstrap.MessageIn) (message string) {

	if len(m.Payload) > 0 {
		if err := json.Unmarshal(m.Payload, &message); err != nil {
			message = err.Error()
			return
		}
	}
	return
}
