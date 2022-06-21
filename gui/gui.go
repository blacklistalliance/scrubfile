package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	die(err)
	win.SetTitle("Blacklist Alliance Scrub")
	win.Connect(
		"destroy", func() {
			gtk.MainQuit()
		},
	)
	st, err := gtk.StackNew()
	die(err)
	img, err := gtk.ImageNewFromFile("logo-mini.png")
	die(err)
	st.Add(img)
	win.Add(st)

	// Create a new label widget to show in the window.
	l, err := gtk.LabelNew("Hello, gotk3!")
	die(err)

	// Add the label to the window.
	win.Add(l)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}