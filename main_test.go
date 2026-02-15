package main

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/requests/ui"
	typedefs "github.com/andreykaipov/goobs/api/typedefs"
)

func getClient(t *testing.T) (*goobs.Client, func()) {
	t.Helper()
	client, err := connectObs(ObsConfig{
		Host:     os.Getenv("OBS_HOST"),
		Port:     4455,
		Password: os.Getenv("OBS_PASSWORD"),
		Timeout:  5,
	})
	if err != nil {
		t.Fatalf("Failed to connect to OBS: %v", err)
	}
	return client, func() {
		if err := client.Disconnect(); err != nil {
			t.Fatalf("Failed to disconnect from OBS: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	client, err := connectObs(ObsConfig{
		Host:     os.Getenv("OBS_HOST"),
		Port:     4455,
		Password: os.Getenv("OBS_PASSWORD"),
		Timeout:  5,
	})
	if err != nil {
		os.Exit(1)
	}

	setup(client)

	// Run the tests
	exitCode := m.Run()

	teardown(client)
	client.Disconnect()

	// Exit with the appropriate code
	os.Exit(exitCode)
}

// nolint: misspell
func setup(client *goobs.Client) {
	client.Config.SetStreamServiceSettings(config.NewSetStreamServiceSettingsParams().
		WithStreamServiceType("rtmp_common").
		WithStreamServiceSettings(&typedefs.StreamServiceSettings{
			Server: "auto",
			Key:    os.Getenv("OBS_STREAM_KEY"),
		}))

	client.Config.CreateProfile(config.NewCreateProfileParams().
		WithProfileName("gobs-test-profile"))
	time.Sleep(100 * time.Millisecond) // Wait for the profile to be created
	client.Config.SetProfileParameter(config.NewSetProfileParameterParams().
		WithParameterCategory("SimpleOutput").
		WithParameterName("RecRB").
		WithParameterValue("true"))
	// hack to ensure the Replay Buffer setting is applied
	client.Config.SetCurrentProfile(config.NewSetCurrentProfileParams().
		WithProfileName("Untitled"))
	client.Config.SetCurrentProfile(config.NewSetCurrentProfileParams().
		WithProfileName("gobs-test-profile"))

	client.Scenes.CreateScene(scenes.NewCreateSceneParams().
		WithSceneName("gobs-test-scene"))
	client.Inputs.CreateInput(inputs.NewCreateInputParams().
		WithSceneName("gobs-test-scene").
		WithInputName("gobs-test-input").
		WithInputKind("color_source_v3").
		WithInputSettings(map[string]any{
			"color":   3279460728,
			"width":   1920,
			"height":  1080,
			"visible": true,
		}).
		WithSceneItemEnabled(true))
	client.Inputs.CreateInput(inputs.NewCreateInputParams().
		WithSceneName("gobs-test-scene").
		WithInputName("gobs-test-input-2").
		WithInputKind("color_source_v3").
		WithInputSettings(map[string]any{
			"color":   1789347616,
			"width":   720,
			"height":  480,
			"visible": true,
		}).
		WithSceneItemEnabled(true))

	// ensure Desktop Audio input is created
	desktopAudioKinds := map[string]string{
		"windows": "wasapi_output_capture",
		"linux":   "pulse_output_capture",
		"darwin":  "coreaudio_output_capture",
	}
	platform := os.Getenv("GOBS_TEST_PLATFORM")
	if platform == "" {
		platform = runtime.GOOS
	}
	client.Inputs.CreateInput(inputs.NewCreateInputParams().
		WithSceneName("gobs-test-scene").
		WithInputName("Desktop Audio").
		WithInputKind(desktopAudioKinds[platform]).
		WithInputSettings(map[string]any{
			"device_id": "default",
		}))
	// ensure Mic/Aux input is created
	micKinds := map[string]string{
		"windows": "wasapi_input_capture",
		"linux":   "pulse_input_capture",
		"darwin":  "coreaudio_input_capture",
	}
	client.Inputs.CreateInput(inputs.NewCreateInputParams().
		WithSceneName("gobs-test-scene").
		WithInputName("Mic/Aux").
		WithInputKind(micKinds[platform]).
		WithInputSettings(map[string]any{
			"device_id": "default",
		}))

	// Create source filter on an audio input
	client.Filters.CreateSourceFilter(filters.NewCreateSourceFilterParams().
		WithSourceName("Mic/Aux").
		WithFilterName("test_filter").
		WithFilterKind("compressor_filter").
		WithFilterSettings(map[string]any{
			"threshold":        -20,
			"ratio":            4,
			"attack_time":      10,
			"release_time":     100,
			"output_gain":      -3.6,
			"sidechain_source": nil,
		}))

	// Create source filter on a scene
	client.Filters.CreateSourceFilter(filters.NewCreateSourceFilterParams().
		WithSourceName("gobs-test-scene").
		WithFilterName("test_filter").
		WithFilterKind("luma_key_filter_v2").
		WithFilterSettings(map[string]any{
			"luma": 0.5,
		}))
}

func teardown(client *goobs.Client) {
	client.Config.RemoveProfile(config.NewRemoveProfileParams().
		WithProfileName("gobs-test-profile"))

	client.Filters.RemoveSourceFilter(filters.NewRemoveSourceFilterParams().
		WithSourceName("Mic/Aux").
		WithFilterName("test_filter"))
	client.Filters.RemoveSourceFilter(filters.NewRemoveSourceFilterParams().
		WithSourceName("gobs-test-scene").
		WithFilterName("test_filter"))

	client.Scenes.RemoveScene(scenes.NewRemoveSceneParams().
		WithSceneName("gobs-test-scene"))

	client.Config.SetCurrentSceneCollection(config.NewSetCurrentSceneCollectionParams().
		WithSceneCollectionName("Untitled"))

	client.Stream.StopStream()
	client.Record.StopRecord()
	client.Outputs.StopReplayBuffer()
	client.Ui.SetStudioModeEnabled(ui.NewSetStudioModeEnabledParams().
		WithStudioModeEnabled(false))
}
