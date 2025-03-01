# LunarFetch Image Display Tests

This directory contains test scripts for verifying the image display functionality of LunarFetch.

## Available Tests

### Basic Image Display Test

The `test-image.go` script tests the basic image display functionality:

```bash
go run tests/image/test-image.go <path-to-image>
```

This test verifies that:
- Images can be loaded from the specified path
- The terminal size is correctly detected
- The image is properly rendered using the best available protocol
- The image is displayed with the correct dimensions

## Image Display Features

The image display functionality in LunarFetch supports:

1. **Multiple Protocols**: Auto-detection of the best available protocol (Sixel, Kitty, iTerm2, Chafa, etc.)
2. **Optimal Sizing**: Automatic adjustment of image dimensions based on terminal size
3. **Aspect Ratio Preservation**: Maintains the original aspect ratio of the image
4. **Debug Mode**: Provides detailed information about the image display process

## Running the Tests

To run the image display test:

```bash
# With a specific image
go run tests/image/test-image.go /path/to/your/image.png

# With the default test image
go run tests/image/test-image.go ~/.config/lunarfetch/images/test.png
```

## Troubleshooting

If the image doesn't display correctly:

1. Ensure your terminal supports at least one of the image display protocols
2. Check that the image file exists and is a valid image format
3. Try running with debug mode to see detailed information:
   ```bash
   go run tests/image/test-image.go --debug /path/to/your/image.png
   ```
4. Verify that the terminal dimensions are correctly detected
5. Try specifying a different protocol manually:
   ```bash
   # Example with Sixel protocol
   PROTOCOL=sixel go run tests/image/test-image.go /path/to/your/image.png
   ``` 