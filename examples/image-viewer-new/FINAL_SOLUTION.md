# Image Zoom Fix - Complete Solution

## Root Cause Analysis

The zoom feature was completely broken due to **two critical bugs**:

### Bug 1: `View()` Ignored Zoom Level
**File**: `ui/components/image/image.go:267`

```go
// BEFORE (BROKEN):
displayW, displayH := img.calculateDisplaySize()  // ❌ Ignores zoomLevel!

// AFTER (FIXED):
displayW, displayH := img.ScaledSize()  // ✅ Uses zoomLevel
```

### Bug 2: `Scale()` Method Prevented Upscaling
**File**: `ui/components/image/decoder/decoder.go:147-149`

```go
// BEFORE (BROKEN):
// Don't upscale
if ratio > 1.0 {
    ratio = 1.0  // ❌ Prevents all zooming in!
}
```

This meant that even when `View()` passed larger dimensions (e.g., 80x80 for 2x zoom), the `Scale()` method would force the ratio back to 1.0, returning the original size.

## The Fix

### 1. Added `ScaleWithOptions()` Method
**File**: `ui/components/image/decoder/decoder.go`

Added new method that allows controlled upscaling:

```go
func (d *ImageData) ScaleWithOptions(maxWidth, maxHeight int, allowUpscale bool) (int, int) {
    // Calculate scaling ratio
    ratio := ...
    
    // Only upscale if explicitly allowed
    if !allowUpscale && ratio > 1.0 {
        ratio = 1.0
    }
    
    return int(float64(d.Width) * ratio), int(float64(d.Height) * ratio)
}
```

The original `Scale()` method now calls this with `allowUpscale=false` to preserve backward compatibility.

### 2. Modified All Renderers
Updated all renderers to use `ScaleWithOptions(xxx, xxx, true)`:

| Renderer | File | Method Modified |
|----------|------|-----------------|
| Blocks | `blocks.go:58` | `renderBlocks()` |
| Blocks | `blocks.go:99` | `renderASCII()` |
| Sixel | `sixel.go:63` | `scaleToFit()` |
| Kitty | `kitty.go:83` | `scaleToFit()` |
| iTerm2 | `iterm.go:78` | `scaleToFit()` |

### 3. Rewrote `ScaledSize()` Algorithm
**File**: `ui/components/image/image.go:626-695`

New algorithm properly calculates zoomed dimensions:

```go
func (img *Image) ScaledSize() (int, int) {
    // Step 1: Get display bounds
    displayW, displayH := img.calculateDisplaySize()
    
    // Step 2: Calculate base size based on zoom mode
    var baseW, baseH int
    switch img.zoomMode {
    case ZoomFit:
        // Fit within display (maintain aspect ratio)
        baseW = displayW
        baseH = int(float64(baseW) / aspectRatio)
        if baseH > displayH {
            baseH = displayH
            baseW = int(float64(baseH) * aspectRatio)
        }
    // ... other modes
    }
    
    // Step 3: Apply zoom level
    scaledW := int(float64(baseW) * img.zoomLevel)
    scaledH := int(float64(baseH) * img.zoomLevel)
    
    return scaledW, scaledH
}
```

## Complete Data Flow

### Before (Broken)
```
User presses +
→ ZoomIn() sets zoomLevel = 1.1
→ View() calls calculateDisplaySize() → returns 80x40
→ blocks.Render(80, 40)
→ blocks.Scale(80, 80) → ratio forced to 1.0 → returns 40x40
→ Image appears same size ❌
```

### After (Fixed)
```
User presses +
→ ZoomIn() sets zoomLevel = 1.1
→ View() calls ScaledSize() → calculates base 40x40, applies zoom → returns 44x44
→ blocks.Render(44, 44)
→ blocks.ScaleWithOptions(44, 88, true) → ratio = 1.1 → returns 44x88
→ Image appears 10% larger ✅
```

## Test Results

All new tests pass:

```
✅ TestScaleWithOptions - All upscaling scenarios
✅ TestScaleBackwardCompatibility - Old behavior preserved
✅ TestScaledSizeAlgorithm - All zoom modes work correctly
✅ TestWideImageScaledSize - Aspect ratio maintained
```

Example output:
```
Fit at 1.0x zoom: 40x40 (Square 100x100 image in 80x40)
Fit at 2.0x zoom: 80x80 (2x zoom works!)
Fit at 0.5x zoom: 20x20 (0.5x zoom works!)
Stretch at 2.0x zoom: 160x80 (Ignores aspect ratio as expected)
```

## Files Modified

1. **ui/components/image/image.go**
   - Line 267: Use `ScaledSize()` instead of `calculateDisplaySize()`
   - Lines 626-695: Rewrote `ScaledSize()` algorithm

2. **ui/components/image/decoder/decoder.go**
   - Added `ScaleWithOptions()` method with `allowUpscale` parameter
   - Original `Scale()` preserved for backward compatibility

3. **ui/components/image/blocks.go**
   - Line 58: Use `ScaleWithOptions(..., true)` in `renderBlocks()`
   - Line 99: Use `ScaleWithOptions(..., true)` in `renderASCII()`

4. **ui/components/image/sixel.go**
   - Line 63: Use `ScaleWithOptions(..., true)` in `scaleToFit()`

5. **ui/components/image/kitty.go**
   - Line 83: Use `ScaleWithOptions(..., true)` in `scaleToFit()`

6. **ui/components/image/iterm.go**
   - Line 78: Use `ScaleWithOptions(..., true)` in `scaleToFit()`

7. **ui/components/image/scaledsize_test.go** (NEW)
   - Comprehensive tests for `ScaledSize()` algorithm

8. **ui/components/image/decoder/scale_test.go** (NEW)
   - Tests for `ScaleWithOptions()` method

## Usage

```bash
cd examples/image-viewer-new
go run main.go <path-to-image>
```

### Key Controls
- `+` / `=` : Zoom in
- `-` / `_` : Zoom out
- `0` : Reset to 100%
- `*` : Zoom to 200%
- `%` : Zoom to 50%
- `[` / `]` : Fine zoom (1% steps)
- `m` : Cycle zoom modes
- `f` / `F` / `s` / `o` : Select specific zoom mode

## Verification

The zoom feature now works correctly:
- ✅ Zoom in increases image size
- ✅ Zoom out decreases image size
- ✅ All zoom levels work (0.1x to 5.0x)
- ✅ Aspect ratio maintained in Fit mode
- ✅ Image can exceed terminal bounds for inspection
- ✅ Backward compatibility preserved
