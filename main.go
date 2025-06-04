// Package main provides a command-line interface (CLI) tool for interacting with OBS WebSocket.
// It allows users to manage various aspects of OBS, such as scenes, inputs, recording, streaming,
// and more, by leveraging the goobs library for communication with the OBS WebSocket server.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
	"github.com/andreykaipov/goobs"
	kongdotenv "github.com/titusjaka/kong-dotenv-go"
)

// ObsConfig holds the configuration for connecting to the OBS WebSocket server.
type ObsConfig struct {
	Host     string `flag:"host"     help:"Host to connect to."          default:"localhost" env:"OBS_HOST"     short:"H"`
	Port     int    `flag:"port"     help:"Port to connect to."          default:"4455"      env:"OBS_PORT"     short:"P"`
	Password string `flag:"password" help:"Password for authentication." default:""          env:"OBS_PASSWORD" short:"p"`
	Timeout  int    `flag:"timeout"  help:"Timeout in seconds."          default:"5"         env:"OBS_TIMEOUT"  short:"T"`
}

// CLI is the main command line interface structure.
// It embeds the ObsConfig struct to inherit its fields and flags.
type CLI struct {
	ObsConfig `embed:"" help:"OBS WebSocket configuration."`

	Man     mangokong.ManFlag `help:"Print man page."`
	Version VersionFlag       `help:"Print gobs-cli version information and quit" name:"version" short:"v"`

	ObsVersion      ObsVersionCmd      `help:"Print OBS client and websocket version." cmd:"" aliases:"v"`
	Scene           SceneCmd           `help:"Manage scenes."                          cmd:"" aliases:"sc"`
	Sceneitem       SceneItemCmd       `help:"Manage scene items."                     cmd:"" aliases:"si"`
	Group           GroupCmd           `help:"Manage groups."                          cmd:"" aliases:"g"`
	Input           InputCmd           `help:"Manage inputs."                          cmd:"" aliases:"i"`
	Record          RecordCmd          `help:"Manage recording."                       cmd:"" aliases:"rec"`
	Stream          StreamCmd          `help:"Manage streaming."                       cmd:"" aliases:"st"`
	Scenecollection SceneCollectionCmd `help:"Manage scene collections."               cmd:"" aliases:"scn"`
	Profile         ProfileCmd         `help:"Manage profiles."                        cmd:"" aliases:"p"`
	Replaybuffer    ReplayBufferCmd    `help:"Manage replay buffer."                   cmd:"" aliases:"rb"`
	Studiomode      StudioModeCmd      `help:"Manage studio mode."                     cmd:"" aliases:"sm"`
	Virtualcam      VirtualCamCmd      `help:"Manage virtual camera."                  cmd:"" aliases:"vc"`
	Hotkey          HotkeyCmd          `help:"Manage hotkeys."                         cmd:"" aliases:"hk"`
	Filter          FilterCmd          `help:"Manage filters."                         cmd:"" aliases:"f"`
	Projector       ProjectorCmd       `help:"Manage projectors."                      cmd:"" aliases:"prj"`
	Screenshot      ScreenshotCmd      `help:"Take screenshots."                       cmd:"" aliases:"ss"`
}

type context struct {
	Client *goobs.Client
	Out    io.Writer
}

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user config directory: %v\n", err)
		os.Exit(1)
	}

	var cli CLI
	ctx := kong.Parse(
		&cli,
		kong.Name("GOBS-CLI"),
		kong.Description("A command line tool to interact with OBS Websocket."),
		kong.Configuration(kongdotenv.ENVFileReader, ".env", filepath.Join(userConfigDir, "gobs-cli", "config.env")),
	)

	client, err := connectObs(cli.ObsConfig)
	ctx.FatalIfErrorf(err)

	ctx.Bind(&context{
		Client: client,
		Out:    os.Stdout,
	})

	ctx.FatalIfErrorf(run(ctx, client))
}

// connectObs creates a new OBS client and connects to the OBS WebSocket server.
func connectObs(cfg ObsConfig) (*goobs.Client, error) {
	client, err := goobs.New(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		goobs.WithPassword(cfg.Password),
		goobs.WithResponseTimeout(time.Duration(cfg.Timeout)*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// run executes the command line interface.
// It disconnects the OBS client after the command is executed.
func run(ctx *kong.Context, client *goobs.Client) error {
	defer func() error {
		if err := client.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect from OBS: %w", err)
		}
		return nil
	}()

	return ctx.Run()
}
