# Image Viewer Zoom and Display Fix - Complete Solution

## Problems Fixed

### 1. Zoom Not Working (Root Cause)
**File**: `ui/components/image/image.go:267`

**Problem**: `View()` method called `calculateDisplaySize()` instead of `ScaledSize()`, completely ignoring the zoom level.

**Fix**:
```diff
- displayW, displayH := img.calculateDisplaySize()
+ displayW, displayH := img.ScaledSize()
```

### 2. ScaledSize() Algorithm Issue
**File**: `ui/components/image/image.go:626-662`

**Problem**: The old `ScaledSize()` implementation had issues:
- Used `img.imgData.Scale()` which has a "don't enlarge" constraint (ratio ≤ 1.0)
- Didn't properly calculate base size before applying zoom
- Zoom was applied after constraint, making it ineffective

**Solution**: Rewrote `ScaledSize()` to:
1. Calculate base dimensions that fit in display area (respecting zoom mode)
2. Apply zoom level as multiplier on base size
3. Maintain aspect ratio correctly
4. Allow zooming beyond display bounds (for inspection)

## How It Works Now

### ScaledSize() Algorithm

```go
func (img *Image) ScaledSize() (int, int) {
    // Step 1: Get display bounds (terminal size minus margins)
    displayW, displayH := img.calculateDisplaySize()
    
    // Step 2: Calculate base size based on zoom mode
    switch img.zoomMode {
    case ZoomFit:
        // Fit within display (maintain aspect ratio)
        baseW = displayW
        baseH = int(float64(baseW) / aspectRatio)
        if baseH > displayH {
            baseH = displayH  // Limited by height
            baseW = int(float64(baseH) * aspectRatio)
        }
        
    case ZoomFill:
        // Fill display (may crop)
        // Use larger scale to ensure coverage
        
    case ZoomStretch:
        // Stretch to fill (ignore aspect ratio)
        baseW = displayW
        baseH = displayH
        
    case ZoomOriginal:
        // Use original image size
        baseW = origW
        baseH = origH
    }
    
    // Step 3: Apply zoom level
    scaledW := int(float64(baseW) * img.zoomLevel)
    scaledH := int(float64(baseH) * img.zoomLevel)
    
    return scaledW, scaledH
}
```

### Example Behavior

For a **100x100 square image** in **80x40 display**:

| Mode | Zoom | Base Size | Scaled Size | Notes |
|------|------|-----------|-------------|-------|
| Fit | 1.0x | 40x40 | 40x40 | Limited by height |
| Fit | 2.0x | 40x40 | 80x80 | 2x base (may exceed display) |
| Fit | 0.5x | 40x40 | 20x20 | Half base |
| Stretch | 2.0x | 80x40 | 160x80 | Ignores aspect ratio |
| Original | 1.0x | 100x100 | 100x100 | Original size |

For a **200x100 wide image** in **80x40 display**:

| Mode | Zoom | Base Size | Scaled Size | Notes |
|------|------|-----------|-------------|-------|
| Fit | 1.0x | 80x40 | 80x40 | Fits perfectly (2:1 ratio) |
| Fit | 2.0x | 80x40 | 160x80 | 2x base |
| Original | 1.0x | 200x100 | 200x100 | Exceeds display |

## Key Behaviors

### 1. Fit Mode (ZoomFit)
- Maintains aspect ratio
- Fits within display area at 1.0x zoom
- At higher zoom, may exceed display bounds
- Best for general viewing

### 2. Fill Mode (ZoomFill)
- May crop image to fill display
- Ensures no empty space
- Good for wallpapers/backgrounds

### 3. Stretch Mode (ZoomStretch)
- Ignores aspect ratio
- Stretches to fill display exactly
- May distort image

### 4. Original Mode (ZoomOriginal)
- Uses actual image dimensions
- Zoom directly multiplies original size
- Good for pixel-perfect inspection

## Test Results

All tests pass:
```
✅ TestScaledSizeAlgorithm - All zoom modes and levels
✅ TestWideImageScaledSize - Wide images maintain aspect ratio
✅ ZoomIn/ZoomOut - Correctly adjust zoomLevel
✅ SetScale - Direct zoom level setting
✅ Aspect ratio preservation in Fit mode
```

## Usage

```bash
cd examples/image-viewer-new
go run main.go <path-to-image>
```

### Key Controls
- `+` / `=` : Zoom in (0.1x increment)
- `-` / `_` : Zoom out (0.1x decrement)
- `0` : Reset to 1.0x zoom
- `*` : Set zoom to 2.0x
- `%` : Set zoom to 0.5x
- `[` / `]` : Fine zoom (0.01 steps)
- `m` : Cycle zoom modes
- `f` : Fit mode
- `F` : Fill mode
- `s` : Stretch mode
- `o` : Original size mode

## Files Modified

1. **ui/components/image/image.go**
   - Line 267: Use `ScaledSize()` instead of `calculateDisplaySize()`
   - Lines 626-695: Rewrote `ScaledSize()` with proper algorithm

2. **ui/components/image/scaledsize_test.go** (NEW)
   - Comprehensive tests for ScaledSize() algorithm
   - Tests for all zoom modes and levels
   - Aspect ratio verification

## Why This Fix Works

### Before
```
User presses +
→ ZoomIn() sets zoomLevel = 1.1
→ View() calls calculateDisplaySize()
→ Returns 80x40 (ignores zoomLevel!)
→ Image appears same size ❌
```

### After
```
User presses +
→ ZoomIn() sets zoomLevel = 1.1
→ View() calls ScaledSize()
→ Calculates base size: 40x40
→ Applies zoom: 40x40 * 1.1 = 44x44
→ Image appears larger ✅
```

## Handling Large Images

When zoomed image exceeds terminal size:
- **Blocks renderer**: Will render full size (may scroll off screen)
- **Sixel/Kitty**: Will render at requested pixel size
- **User can**: Zoom out or switch to Fit mode

This is intentional - users can zoom in to inspect details, even if it exceeds the viewport.

## Future Improvements

Potential enhancements:
1. Add viewport/scroll support for large zoomed images
2. Add pan controls (arrow keys) when zoomed
3. Show scroll indicators when image exceeds bounds
4. Optimize rendering for very large zoom levels
