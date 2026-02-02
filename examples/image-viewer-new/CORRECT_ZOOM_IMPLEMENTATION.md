# Image Zoom Fix - Correct Implementation

## The Problem

You were absolutely right! The original implementation was fundamentally wrong. Zoom was incorrectly implemented as **changing display size** instead of **changing sampling resolution**.

## What Zoom Should Do

**Zoom changes the image resolution/detail level, NOT the character cell size:**

| Zoom Level | Sampled Area | Display Size | Effect |
|------------|--------------|--------------|--------|
| 1.0x | 40x40 pixels | 80x40 chars | Normal (1 pixel = 2 chars vertically) |
| 2.0x | 20x20 pixels | 80x40 chars | **Zoomed in** - each char shows 0.5x0.5 pixel area (higher detail) |
| 4.0x | 10x10 pixels | 80x40 chars | **Very zoomed in** - each char shows 0.25x0.25 pixel area |
| 0.5x | 80x80 pixels | 80x40 chars | **Zoomed out** - each char shows 2x2 pixel area (lower detail) |

## The Fix

### Before (Wrong)
```go
// ScaledSize() multiplied base size by zoom level
scaledW := int(float64(baseW) * img.zoomLevel)  // ❌ Makes display bigger
scaledH := int(float64(baseH) * img.zoomLevel)
```

This caused:
- 2x zoom → 160x80 character cells (exceeds terminal width)
- Image just got bigger, not clearer
- Content got cut off

### After (Correct)
```go
// ScaledSize() divides base size by zoom level  
sampledW := int(float64(baseW) / img.zoomLevel)  // ✅ Changes resolution
sampledH := int(float64(baseH) / img.zoomLevel)
```

This causes:
- 2x zoom → Sample 20x20 pixels from original
- Display still 80x40 chars
- **Each character cell shows higher detail** (fewer pixels per cell)

## How It Works

### Example: 100x100 pixel image in 80x40 terminal

**At 1.0x zoom:**
```
Base size: 40x40 (fits in 80x40 with aspect ratio)
Sampled: 40x40 pixels
Display: 80x40 characters
Result: 1 pixel = 2x1 character cells (normal)
```

**At 2.0x zoom:**
```
Base size: 40x40
Sampled: 20x20 pixels (40 / 2.0)
Display: 80x40 characters (same!)
Result: 1 pixel = 4x2 character cells (higher detail!)
```

**At 4.0x zoom:**
```
Base size: 40x40
Sampled: 10x10 pixels (40 / 4.0)
Display: 80x40 characters (same!)
Result: 1 pixel = 8x4 character cells (maximum detail!)
```

### Visual Representation

```
1.0x zoom (normal):
┌────────────────────────────────────┐
│ Each ▀ shows 2x1 pixels           │
│ Image looks normal                 │
└────────────────────────────────────┘

2.0x zoom (zoomed in 2x):
┌────────────────────────────────────┐
│ Each ▀ shows 1x0.5 pixels         │
│ Same size, MORE DETAIL visible     │
│ You can see finer details          │
└────────────────────────────────────┘

0.5x zoom (zoomed out 2x):
┌────────────────────────────────────┐
│ Each ▀ shows 4x2 pixels           │
│ Same size, LESS detail visible     │
│ You see more of the image         │
└────────────────────────────────────┘
```

## Key Insight

**Zoom is about sampling density**, not display size.

The renderer (`blocks.go`) does this:
```go
// For each character cell (x, y):
imgX := (x * scaledW) / width   // Map to original image
imgY := (y * scaledH) / height
pixel := GetPixelColor(imgX, imgY)
```

When `scaledW, scaledH` are smaller (higher zoom), each character cell samples a **smaller region** of the original image, showing **more detail**.

## Files Modified

1. **ui/components/image/image.go** - `ScaledSize()` method
   - Changed from `base * zoomLevel` to `base / zoomLevel`
   - Added bounds checking (min 1px, max original size)

2. **ui/components/image/decoder/decoder.go**
   - Kept `ScaleWithOptions()` for future use
   - Original `Scale()` still prevents upscaling (correct for base calculation)

3. **All renderers** - Reverted to use `Scale()` (not `ScaleWithOptions`)
   - `blocks.go` - Both `renderBlocks()` and `renderASCII()`
   - `sixel.go` - `scaleToFit()`
   - `kitty.go` - `scaleToFit()`
   - `iterm.go` - `scaleToFit()`

## Test Results

```
✅ TestResolutionBasedZoom - All scenarios pass
✅ TestZoomDoesNotChangeDisplaySize - Verified display size is constant

Sample output:
  1.0x zoom: Sample 40x40 pixels
  2.0x zoom: Sample 20x20 pixels (2x detail!)
  4.0x zoom: Sample 10x10 pixels (4x detail!)
  0.5x zoom: Sample 80x80 pixels (0.5x detail)
```

## Usage

```bash
cd examples/image-viewer-new
go run main.go <path-to-image>
```

Press `+` to zoom in (see more detail) or `-` to zoom out (see more of the image).

The image will remain the same size on screen but will show **different levels of detail** based on zoom level.
