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
	Host     string `flag:"host"     help:"Host to connect to."          default:"localhost" env:"OBS_HOST"     short:"H" completion-short-enabled:"false"`
	Port     int    `flag:"port"     help:"Port to connect to."          default:"4455"      env:"OBS_PORT"     short:"P" completion-short-enabled:"false"`
	Password string `flag:"password" help:"Password for authentication." default:""          env:"OBS_PASSWORD" short:"p" completion-short-enabled:"false"`
	Timeout  int    `flag:"timeout"  help:"Timeout in seconds."          default:"5"         env:"OBS_TIMEOUT"  short:"T" completion-short-enabled:"false"`
}

// StyleConfig holds the configuration for styling the CLI output.
type StyleConfig struct {
	Style    string `flag:"style"     help:"Style used in output."                   default:""      env:"GOBS_STYLE"           short:"s" enum:",red,magenta,purple,blue,cyan,green,yellow,orange,white,grey,navy,black" completion-short-enabled:"false"`
	NoBorder bool   `flag:"no-border" help:"Disable table border styling in output." default:"false" env:"GOBS_STYLE_NO_BORDER" short:"b"                                                                                completion-short-enabled:"false"`
}

// CLI is the main command line interface structure.
// It embeds ObsConfig and StyleConfig to provide configuration options.
type CLI struct {
	ObsConfig   `embed:"" help:"OBS WebSocket configuration."`
	StyleConfig `embed:"" help:"Style configuration."`

	Man     mangokong.ManFlag `help:"Print man page."`
	Version VersionFlag       `help:"Print gobs-cli version information and quit" name:"version" short:"v"`

	Completion kongcompletion.Completion `help:"Generate shell completion scripts." cmd:""`

	ObsVersion      ObsVersionCmd      `cmd:"" help:"Print OBS client and websocket version." aliases:"v"   completion-command-alias-enabled:"false"`
	Scene           SceneCmd           `cmd:"" help:"Manage scenes."                          aliases:"sc"  completion-command-alias-enabled:"false" group:"Scene"`
	Sceneitem       SceneItemCmd       `cmd:"" help:"Manage scene items."                     aliases:"si"  completion-command-alias-enabled:"false" group:"Scene Item"`
	Group           GroupCmd           `cmd:"" help:"Manage groups."                          aliases:"g"   completion-command-alias-enabled:"false" group:"Group"`
	Input           InputCmd           `cmd:"" help:"Manage inputs."                          aliases:"i"   completion-command-alias-enabled:"false" group:"Input"`
	Text            TextCmd            `cmd:"" help:"Manage text inputs."                     aliases:"t"   completion-command-alias-enabled:"false" group:"Text Input"`
	Record          RecordCmd          `cmd:"" help:"Manage recording."                       aliases:"rec" completion-command-alias-enabled:"false" group:"Recording"`
	Stream          StreamCmd          `cmd:"" help:"Manage streaming."                       aliases:"st"  completion-command-alias-enabled:"false" group:"Streaming"`
	Scenecollection SceneCollectionCmd `cmd:"" help:"Manage scene collections."               aliases:"scn" completion-command-alias-enabled:"false" group:"Scene Collection"`
	Profile         ProfileCmd         `cmd:"" help:"Manage profiles."                        aliases:"p"   completion-command-alias-enabled:"false" group:"Profile"`
	Replaybuffer    ReplayBufferCmd    `cmd:"" help:"Manage replay buffer."                   aliases:"rb"  completion-command-alias-enabled:"false" group:"Replay Buffer"`
	Studiomode      StudioModeCmd      `cmd:"" help:"Manage studio mode."                     aliases:"sm"  completion-command-alias-enabled:"false" group:"Studio Mode"`
	Virtualcam      VirtualCamCmd      `cmd:"" help:"Manage virtual camera."                  aliases:"vc"  completion-command-alias-enabled:"false" group:"Virtual Camera"`
	Hotkey          HotkeyCmd          `cmd:"" help:"Manage hotkeys."                         aliases:"hk"  completion-command-alias-enabled:"false" group:"Hotkey"`
	Filter          FilterCmd          `cmd:"" help:"Manage filters."                         aliases:"f"   completion-command-alias-enabled:"false" group:"Filter"`
	Projector       ProjectorCmd       `cmd:"" help:"Manage projectors."                      aliases:"prj" completion-command-alias-enabled:"false" group:"Projector"`
	Screenshot      ScreenshotCmd      `cmd:"" help:"Take screenshots."                       aliases:"ss"  completion-command-alias-enabled:"false" group:"Screenshot"`
	Settings        SettingsCmd        `cmd:"" help:"Manage video and profile settings."      aliases:"set" completion-command-alias-enabled:"false" group:"Settings"`
	Media           MediaCmd           `cmd:"" help:"Manage media inputs."                    aliases:"mi"  completion-command-alias-enabled:"false" group:"Media Input"`
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
				if version != "" {
					return version
				}

				info, ok := debug.ReadBuildInfo()
				if !ok {
					return "(unable to read version)"
				}
				return strings.Split(info.Main.Version, "-")[0]
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
