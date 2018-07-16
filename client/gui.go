package main

import (
	"flag"
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Constants
const htmlAbout = `Welcome on <b>Berty</b> POC!<br>
Developed by aeddi using libp2p, crypto and astilectron.`

var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

func initGui() {
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
							if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
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
						Role: astilectron.MenuItemRoleToggleDevTools,
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
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				// var n = a.NewNotification(&astilectron.NotificationOptions{
				// 	Body: "My Body",
				// 	// HasReply: astilectron.PtrBool(true), // Only MacOSX
				// 	// Icon: "../icon.png",
				// 	// ReplyPlaceholder: "type your reply here", // Only MacOSX
				// 	Title: "My title",
				// })
				// n.Show()
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(600),
				Width:           astilectron.PtrInt(1100),
				MinHeight:       astilectron.PtrInt(400),
				MinWidth:        astilectron.PtrInt(800),
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
		// Unmarshal payload
		var message string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &message); err != nil {
				payload = err.Error()
				return
			}
		}

		// if payload, err = explore(path); err != nil {
		// 	payload = err.Error()
		// 	return
		// }
		payload = message
	}
	return
}
