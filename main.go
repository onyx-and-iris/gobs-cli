// Package main provides a command-line interface (CLI) tool for interacting with OBS WebSocket.
// It allows users to manage various aspects of OBS, such as scenes, inputs, recording, streaming,
// and more, by leveraging the goobs library for communication with the OBS WebSocket server.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
	"github.com/andreykaipov/goobs"
	kongcompletion "github.com/jotaen/kong-completion"
	kongdotenv "github.com/titusjaka/kong-dotenv-go"
)

var version string // Version of the CLI, set at build time.

// VersionFlag is a custom flag type that prints the version and exits.
type VersionFlag string

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }  // nolint: revive
func (v VersionFlag) IsBool() bool                       { return true } // nolint: revive
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error { // nolint: revive, unparam
	fmt.Printf("gobs-cli version: %s\n", vars["version"])
	app.Exit(0)
	return nil
}

// ObsConfig holds the configuration for connecting to the OBS WebSocket server.
type ObsConfig struct {
	Host     string `flag:"host"     help:"Host to connect to."          default:"localhost" env:"OBS_HOST"     short:"H"`
	Port     int    `flag:"port"     help:"Port to connect to."          default:"4455"      env:"OBS_PORT"     short:"P"`
	Password string `flag:"password" help:"Password for authentication." default:""          env:"OBS_PASSWORD" short:"p"`
	Timeout  int    `flag:"timeout"  help:"Timeout in seconds."          default:"5"         env:"OBS_TIMEOUT"  short:"T"`
}

// StyleConfig holds the configuration for styling the CLI output.
type StyleConfig struct {
	Style    string `help:"Style used in output."                   flag:"style"     default:""      env:"GOBS_STYLE"           short:"s" enum:",red,magenta,purple,blue,cyan,green,yellow,orange,white,grey,navy,black"`
	NoBorder bool   `help:"Disable table border styling in output." flag:"no-border" default:"false" env:"GOBS_STYLE_NO_BORDER" short:"b"`
}

// CLI is the main command line interface structure.
// It embeds ObsConfig and StyleConfig to provide configuration options.
type CLI struct {
	ObsConfig   `embed:"" help:"OBS WebSocket configuration."`
	StyleConfig `embed:"" help:"Style configuration."`

	Man     mangokong.ManFlag `help:"Print man page."`
	Version VersionFlag       `help:"Print gobs-cli version information and quit" name:"version" short:"v"`

	Completion kongcompletion.Completion `help:"Generate shell completion scripts." cmd:"" aliases:"c"`

	ObsVersion      ObsVersionCmd      `help:"Print OBS client and websocket version." cmd:"" aliases:"v"`
	Scene           SceneCmd           `help:"Manage scenes."                          cmd:"" aliases:"sc"  group:"Scene"`
	Sceneitem       SceneItemCmd       `help:"Manage scene items."                     cmd:"" aliases:"si"  group:"Scene Item"`
	Group           GroupCmd           `help:"Manage groups."                          cmd:"" aliases:"g"   group:"Group"`
	Input           InputCmd           `help:"Manage inputs."                          cmd:"" aliases:"i"   group:"Input"`
	Text            TextCmd            `help:"Manage text inputs."                     cmd:"" aliases:"t"   group:"Text Input"`
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
	Settings        SettingsCmd        `help:"Manage video and profile settings."      cmd:"" aliases:"set" group:"Settings"`
	Media           MediaCmd           `help:"Manage media inputs."                    cmd:"" aliases:"mi"  group:"Media Input"`
}

type context struct {
	Client *goobs.Client
	Out    io.Writer
	Style  *Style
}

func newContext(client *goobs.Client, out io.Writer, styleCfg StyleConfig) *context {
	return &context{
		Client: client,
		Out:    out,
		Style:  styleFromFlag(styleCfg),
	}
}

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user config directory: %v\n", err)
		os.Exit(1)
	}

	var cli CLI
	kongcompletion.Register(kong.Must(&cli))
	ctx := kong.Parse(
		&cli,
		kong.Name("gobs-cli"),
		kong.Description("A command line tool to interact with OBS Websocket."),
		kong.Configuration(
			kongdotenv.ENVFileReader,
			".env",
			filepath.Join(userConfigDir, "gobs-cli", "config.env"),
		),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": func() string {
				if version == "" {
					info, ok := debug.ReadBuildInfo()
					if !ok {
						return "(unable to read build info)"
					}
					version = strings.Split(info.Main.Version, "-")[0]
				}
				return version
			}(),
		},
	)

	ctx.FatalIfErrorf(run(ctx, cli.ObsConfig, cli.StyleConfig))
}

// run executes the command line interface.
// It connects to the OBS WebSocket server and binds the context to the selected command.
// It also handles the "completion" command separately to avoid unnecessary connections.
func run(ctx *kong.Context, obsCfg ObsConfig, styleCfg StyleConfig) error {
	if ctx.Selected().Name == "completion" {
		return ctx.Run()
	}

	client, err := connectObs(obsCfg)
	if err != nil {
		return err
	}

	defer func() error {
		if err := client.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect from OBS: %w", err)
		}
		return nil
	}() // nolint: errcheck

	ctx.Bind(newContext(client, os.Stdout, styleCfg))

	return ctx.Run()
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
