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

// StyleConfig holds the configuration for styling the CLI output.
type StyleConfig struct {
	Style string `help:"Style used in output." flag:"style" default:"" env:"GOBS_STYLE" short:"s" enum:",red,magenta,purple,blue,cyan,green,yellow,orange,white,grey,navy,black"`
}

// CLI is the main command line interface structure.
// It embeds the ObsConfig struct to inherit its fields and flags.
type CLI struct {
	ObsConfig   `embed:"" help:"OBS WebSocket configuration."`
	StyleConfig `embed:"" help:"Style configuration."`

	Man     mangokong.ManFlag `help:"Print man page."`
	Version VersionFlag       `help:"Print gobs-cli version information and quit" name:"version" short:"v"`

	ObsVersion      ObsVersionCmd      `help:"Print OBS client and websocket version." cmd:"" aliases:"v"`
	Scene           SceneCmd           `help:"Manage scenes."                          cmd:"" aliases:"sc"  group:"Scene"`
	Sceneitem       SceneItemCmd       `help:"Manage scene items."                     cmd:"" aliases:"si"  group:"Scene Item"`
	Group           GroupCmd           `help:"Manage groups."                          cmd:"" aliases:"g"   group:"Group"`
	Input           InputCmd           `help:"Manage inputs."                          cmd:"" aliases:"i"   group:"Input"`
	Record          RecordCmd          `help:"Manage recording."                       cmd:"" aliases:"rec" group:"Recording"`
	Stream          StreamCmd          `help:"Manage streaming."                       cmd:"" aliases:"st"  group:"Streaming"`
	Scenecollection SceneCollectionCmd `help:"Manage scene collections."               cmd:"" aliases:"scn" group:"Scene Collection"`
	Profile         ProfileCmd         `help:"Manage profiles."                        cmd:"" aliases:"p"   group:"Profile"`
	Replaybuffer    ReplayBufferCmd    `help:"Manage replay buffer."                   cmd:"" aliases:"rb"  group:"Replay Buffer"`
	Studiomode      StudioModeCmd      `help:"Manage studio mode."                     cmd:"" aliases:"sm"  group:"Studio Mode"`
	Virtualcam      VirtualCamCmd      `help:"Manage virtual camera."                  cmd:"" aliases:"vc"  group:"Virtual Camera"`
	Hotkey          HotkeyCmd          `help:"Manage hotkeys."                         cmd:"" aliases:"hk"  group:"Hotkey"`
	Filter          FilterCmd          `help:"Manage filters."                         cmd:"" aliases:"f"   group:"Filter"`
	Projector       ProjectorCmd       `help:"Manage projectors."                      cmd:"" aliases:"prj" group:"Projector"`
	Screenshot      ScreenshotCmd      `help:"Take screenshots."                       cmd:"" aliases:"ss"  group:"Screenshot"`
}

type context struct {
	Client *goobs.Client
	Out    io.Writer
	Style  *Style
}

func newContext(client *goobs.Client, out io.Writer, styleName string) *context {
	return &context{
		Client: client,
		Out:    out,
		Style:  styleFromFlag(styleName),
	}
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

	ctx.Bind(newContext(client, os.Stdout, cli.StyleConfig.Style))

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
