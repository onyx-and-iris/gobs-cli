package main

import (
	"os"
	"testing"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
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
	defer client.Disconnect()

	setup(client)

	// Run the tests
	exitCode := m.Run()

	teardown(client)

	// Exit with the appropriate code
	os.Exit(exitCode)
}

func setup(client *goobs.Client) {
	client.Config.SetStreamServiceSettings(config.NewSetStreamServiceSettingsParams().
		WithStreamServiceType("rtmp_common").
		WithStreamServiceSettings(&typedefs.StreamServiceSettings{
			Server: "auto",
			Key:    os.Getenv("OBS_STREAM_KEY"),
		}))

	client.Config.SetCurrentSceneCollection(config.NewSetCurrentSceneCollectionParams().
		WithSceneCollectionName("test-collection"))

	client.Scenes.CreateScene(scenes.NewCreateSceneParams().
		WithSceneName("gobs-test"))
	client.Inputs.CreateInput(inputs.NewCreateInputParams().
		WithSceneName("gobs-test").
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
		WithSceneName("gobs-test").
		WithInputName("gobs-test-input-2").
		WithInputKind("color_source_v3").
		WithInputSettings(map[string]any{
			"color":   1789347616,
			"width":   720,
			"height":  480,
			"visible": true,
		}).
		WithSceneItemEnabled(true))
}

func teardown(client *goobs.Client) {
	client.Scenes.RemoveScene(scenes.NewRemoveSceneParams().
		WithSceneName("gobs-test"))

	client.Config.SetCurrentSceneCollection(config.NewSetCurrentSceneCollectionParams().
		WithSceneCollectionName("default"))

	client.Stream.StopStream()
	client.Record.StopRecord()
}
