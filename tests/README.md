# LunarFetch Tests

This directory contains test scripts for verifying various functionalities of LunarFetch.

## Available Tests

### Image Display Test

The `test-image.go` script tests the image display functionality:

```bash
go run tests/test-image.go <path-to-image>
```

This test verifies that:
- Images can be loaded from the specified path
- The terminal size is correctly detected
- The image is properly rendered using the best available protocol
- The image is displayed with the correct dimensions

### Configuration Tests

The `config/` directory contains tests for verifying configuration options:

```bash
go run tests/config/main.go
```

See the [Configuration Tests README](config/README.md) for more details.

## Running All Tests

To run all tests, you can use:

```bash
# Run image test with default image
go run tests/test-image.go ~/.config/lunarfetch/images/test.png

# Run configuration tests
go run tests/config/main.go

# Test configurations
go run main.go -c ~/.config/lunarfetch/test-above.json
go run main.go -c ~/.config/lunarfetch/test-side.json
```

## Debugging

Add the `--debug` flag to see detailed information about the test execution:

```bash
go run main.go --debug -c ~/.config/lunarfetch/test-above.json
``` 