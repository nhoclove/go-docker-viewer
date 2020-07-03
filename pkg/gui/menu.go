package gui

import (
	"context"
	"fmt"
	"go-docker-viewer/pkg/docker"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	viewArr          = []string{"top-left", "bottom-left", "right"}
	active           = 0
	dockerFeatureMap = map[int]string{
		0: "Containers",
		1: "Images",
		2: "Volumes",
	}
)

// ShowMenu shows the main menu
func ShowMenu() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	err = keybindings(g)
	if err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("top-left", 0, 0, maxX/3, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = " Docker Menu "
		for _, name := range dockerFeatureMap {
			fmt.Fprintf(v, " %s \n", name)
		}
		v.Editable = true
		v.Wrap = true
	}

	if v, err := g.SetView("bottom-left", 0, 5, maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.Title = " Info "
		v.Wrap = true
		v.Overwrite = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("right", maxX/3+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.Title = " Details "
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	// Keybinding for all views
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	// Keybinding for top-left menu
	if err := g.SetKeybinding("top-left", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		var line string
		_, err := g.SetCurrentView(v.Name())
		if err != nil {
			return err
		}

		_, cy := v.Cursor()
		line, err = v.Line(cy)
		if err != nil {
			return err
		}

		// Access bottom-left view
		bottomLeftView, err := g.View("bottom-left")
		if err != nil {
			return err
		}
		if err = onMenuSelect(line, bottomLeftView); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	// Keybinding for right menu

	// Keybinding for bottom-left menu
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	g.Cursor = true

	active = nextIndex
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func onMenuSelect(feature string, out *gocui.View) error {
	out.Clear()
	out.SetCursor(0, 0)
	out.SetOrigin(0, 0)
	feature = strings.Trim(feature, " ")
	switch feature {
	case "Containers":
		out.Title = " Containers "
		containers, err := docker.ListContainer(context.Background())
		if err != nil {
			return err
		}
		for _, container := range containers {
			fmt.Fprintf(out, " <%s> %s\n", container.State, container.Names[0][1:])
		}
	case "Images":
		out.Title = " Images "
		images, err := docker.ListImages(context.Background())
		if err != nil {
			return err
		}
		for _, image := range images {
			size := float64(image.Size / 1000000)
			fmt.Fprintf(out, " %s %.1fMB\n", image.RepoTags[0], size)
		}
	case "Volumes":
		out.Title = " Volumes "
		volumes, err := docker.ListVolumes(context.Background())
		if err != nil {
			return err
		}
		for _, volume := range volumes {
			name := volume.Name
			if len(volume.Name) > 32 {
				name = volume.Name[:32]
			}
			fmt.Fprintf(out, " <%s> %s\n", volume.Driver, name)
		}
	default:
		return fmt.Errorf("unsupported feature: %s", feature)
	}

	return nil
}
