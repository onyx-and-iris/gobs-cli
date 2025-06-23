# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

# [0.13.0] - 2025-06-23

### Added

-   record split and record chapter commands, see [RecordCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#recordcmd)

# [0.12.1] - 2025-06-21

### Added

-   Various colouring styles, see [Style](https://github.com/onyx-and-iris/gobs-cli/tree/main?tab=readme-ov-file#style)
    -   colouring is applied to list tables as well as highlighted information in stdout/stderr output.
    -   table border styling may be optionally disabled with the --no-border flag.

### Changed

-   if an itemName is passed to a sceneitem command that's in a group, without the --group flag, a friendlier error message is displayed.
    -   it will suggest using *gobs-cli si ls* to list sources in the scene.
-   if an invalid --monitor-index is passed to projector open a friendlier error message is displayed.
    -   it will suggest using *gobs-cli prj ls-m* to list available monitors.


# [0.11.0] - 2025-06-20

### Added

-   input list, scene list and sceneitem list now accept --uuid flag.
-   Active column added to scene list table.

### Changed

-   scene list no longer prints the UUIDs by default, enable it with the --uuid flag.

# [0.10.3] - 2025-06-07

### Added

-   filter list:
    -   --ffmpeg, --vlc flags
    -   Muted column to list table

# [0.10.2] - 2025-06-04

### Added

-   screenshot save command, see [ScreenshotCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#screenshotcmd)

### Fixed

-   filter list:
    -   sourceName arg now defaults to current scene. 
    -   defaults are printed for any unmodified values.
-   sceneitem list:
    -   prints enabled mark instead of true/false

# [0.9.0] - 2025-06-02

### Added

-   --version/-v option. See [Flags](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#flags)

### Changed

-   version command renamed to obs-version

# [0.8.2] - 2025-05-29

### Added

-   record start/stop and stream start/stop commands check outputActive states first. 
    -   Errors are returned if the command cannot be performed.

### Changed

-   The --parent flag for the sceneitem commands has been renamed to --group, see [SceneItemCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#sceneitemcmd)

# [0.8.0] - 2025-05-27

### Added

-   record directory command, see [directory under RecordCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#recordcmd)

### Changed

-   record stop now prints the output path of the recording.


# [0.7.0] - 2025-05-26

### Added

-   projector commands, see [ProjectorCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#projectorcmd)


# [0.6.1] - 2025-05-25

### Added

-   filter commands, see [FilterCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#filtercmd)

### Changed

-   list commands are now printed as tables.
    - This affects group, hotkey, input, profile, scene, scenecollection and sceneitem command groups.

# [0.5.0] - 2025-05-22

### Added

-   hotkey commands, see [HotkeyCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#hotkeycmd)

# [0.4.2] - 2025-05-08

### Added

-   replaybuffer toggle command
-   studiomode enable/disable now print output to console
-   stream start/stop now print output to console
-   Unit tests

# [0.3.1] - 2025-05-02

### Added

-   --man flag for generating/viewing a man page.
-   Ability to load env vars from env files, see the [README](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#environment-variables)

# [0.2.0] - 2025-04-27

### Added

-   sceneitem transform, see *transform* under [SceneItemCmd](https://github.com/onyx-and-iris/gobs-cli?tab=readme-ov-file#sceneitemcmd)

# [0.1.0] - 2025-04-24

### Added

-   Initial release.
