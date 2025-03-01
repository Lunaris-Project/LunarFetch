# LunarFetch Configuration Tests

This directory contains test scripts for verifying the configuration options of LunarFetch, particularly focusing on the image and logo display options.

## Available Tests

### Image Display Position Tests

These tests verify that images can be displayed either above or beside the system information.

- **Image Above Test**: Creates a configuration with the image displayed above the system information.
- **Image Beside Test**: Creates a configuration with the image displayed beside the system information.

## Running the Tests

To run all tests:

```bash
go run tests/config/main.go
```

To run a specific test:

```bash
# For image above test
go run tests/config/main.go above

# For image beside test
go run tests/config/main.go side
```

## Testing the Configurations

After running the tests, you can test the generated configurations with:

```bash
# Test image above configuration
go run main.go -c ~/.config/lunarfetch/test-above.json

# Test image beside configuration
go run main.go -c ~/.config/lunarfetch/test-side.json
```

To see debug information, add the `--debug` flag:

```bash
go run main.go --debug -c ~/.config/lunarfetch/test-above.json
```

## Configuration Details

### Image Above Configuration

- Logo is disabled
- Image is displayed above system information
- Uses auto protocol detection

### Image Beside Configuration

- Logo is enabled and positioned above system information
- Image is displayed beside system information
- Image is shown before system information
- Uses auto protocol detection

## Requirements

- Test image should be located at `~/.config/lunarfetch/images/test.png`
- The terminal should support at least one of the image display protocols (Sixel, Kitty, Chafa, etc.) 