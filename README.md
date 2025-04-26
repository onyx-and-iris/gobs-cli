# gobs-cli

A command line interface for OBS Websocket v5

For an outline of past/future changes refer to: [CHANGELOG](CHANGELOG.md)

## Configuration

#### Flags

Pass `--host`, `--port` and `--password` as flags to the root command, for example:

```console
gobs-cli --host=localhost --port=4455 --password=<websocket password> --help
```

#### Environment Variables

Load connection details from your environment:

```bash
#!/usr/bin/env bash

export OBS_HOST=localhost
export OBS_PORT=4455
export OBS_PASSWORD=<websocket password>
export OBS_TIMEOUT=5
```

## Commands

### VersionCmd

```console
gobs-cli version
```

### SceneCmd

-   list: List all scenes.

```console
gobs-cli scene list
```

-   current: Get the current scene.
    -   flags:

        *optional*
        -   --preview:  Preview scene.

```console
gobs-cli scene current

gobs-cli scene current --preview
```

-   switch: Switch to a scene.
    -   flags:

        *optional*
        -   --preview:  Preview scene.
    -   args: SceneName

```console
gobs-cli scene switch LIVE

gobs-cli scene switch --preview LIVE
```

### SceneItemCmd

-   list: List all scene items.
    -   args: SceneName

```console
gobs-cli sceneitem list LIVE
```

-   show: Show scene item.
    -   flags:

        *optional*
        -   --parent: Parent group name.
    -   args: SceneName ItemName

```console
gobs-cli sceneitem show START "Colour Source"
```

-   hide: Hide scene item.
    -   flags:

        *optional*
        -   --parent: Parent group name.
    -   args: SceneName ItemName

```console
gobs-cli sceneitem hide START "Colour Source"
```

-   toggle: Toggle scene item.
    -   flags:

        *optional*
        -   --parent: Parent group name.
    -   args: SceneName ItemName

```console
gobs-cli sceneitem toggle --parent=test_group START "Colour Source 3"
```

-   visible: Get scene item visibility.
    -   flags:

        *optional*
        -   --parent: Parent group name.
    -   args: SceneName ItemName

```console
gobs-cli sceneitem visible --parent=test_group START "Colour Source 4"
```

### GroupCmd

-   list: List all groups.
    -   args: SceneName

```console
gobs-cli group list START
```

-   show: Show group details.
    -   args: SceneName GroupName

```console
gobs-cli group show START "test_group"
```

-   hide: Hide group.
    -   args: SceneName GroupName

```console
gobs-cli group hide START "test_group"
```

-   toggle: Toggle group.
    -   args: SceneName GroupName

```console
gobs-cli group toggle START "test_group"
```

-   status: Get group status.
    -   args: SceneName GroupName

```console
gobs-cli group status START "test_group"
```

### InputCmd

-   list: List all inputs.
    -   flags:

        *optional*
        -   --input: List all inputs.
        -   --output: List all outputs.
        -   --colour: List all colour sources.

```console
gobs-cli input list

gobs-cli input list --input --colour
```

-   mute: Mute input.
    -   args: InputName

```console
gobs-cli input mute "Mic/Aux"
```

-   unmute: Unmute input.
    -   args: InputName

```console
gobs-cli input unmute "Mic/Aux"
```

-   toggle: Toggle input.
    -   args: InputName

```console
gobs-cli input toggle "Mic/Aux"
```

### RecordCmd

-   start: Start recording.

```console
gobs-cli record start
```

-   stop: Stop recording.

```console
gobs-cli record stop
```

-   status: Get recording status.

```console
gobs-cli record status
```

-   toggle: Toggle recording.

```console
gobs-cli record toggle
```

-   pause: Pause recording.

```console
gobs-cli record pause
```

-   resume: Resume recording.

```console
gobs-cli record resume
```

### StreamCmd

-   start: Start streaming.

```console
gobs-cli stream start
```

-   stop: Stop streaming.

```console
gobs-cli stream stop
```

-   status: Get streaming status.

```console
gobs-cli stream status
```

-   toggle: Toggle streaming.

```console
gobs-cli stream toggle
```

### SceneCollectionCmd

-   list: List scene collections.

```console
gobs-cli scenecollection list
```

-   current: Get current scene collection.

```console
gobs-cli scenecollection current
```

-   switch: Switch scene collection.
    -   args: Name

```console
gobs-cli scenecollection switch test-collection
```

-   create: Create scene collection.
    -   args: Name

```console
gobs-cli scenecollection create test-collection
```

### ProfileCmd

-   list: List profiles.

```console
gobs-cli profile list
```

-   current: Get current profile.

```console
gobs-cli profile current
```

-   switch: Switch profile.
    -   args: Name

```console
gobs-cli profile switch test-profile
```

-   create: Create profile.
    -   args: Name

```console
gobs-cli profile create test-profile
```

-   remove: Remove profile.
    -   args: Name

```console
gobs-cli profile remove test-profile
```

### ReplayBufferCmd

-   start: Start replay buffer.

```console
gobs-cli replaybuffer start
```

-   stop: Stop replay buffer.

```console
gobs-cli replaybuffer stop
```

-   status: Get replay buffer status.

```console
gobs-cli replaybuffer status
```

-   save: Save replay buffer.

```console
gobs-cli replaybuffer save
```

### StudioModeCmd

-   enable: Enable studio mode.

```console
gobs-cli studiomode enable
```

-   disable: Disable studio mode.

```console
gobs-cli studiomode disable
```

-   toggle: Toggle studio mode.

```console
gobs-cli studiomode toggle
```

-   status: Get studio mode status.

```console
gobs-cli studiomode status
```

### VirtualCamCmd

-   start: Start virtual camera.

```console
gobs-cli virtualcam start
```

-   stop: Stop virtual camera.

```console
gobs-cli virtualcam stop
```

-   toggle: Toggle virtual camera.

```console
gobs-cli virtualcam toggle
```

-   status: Get virtual camera status.

```console
gobs-cli virtualcam status
```