package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"path"
	"slices"
	"strings"
)

// An RPG Maker MV side-view battler is a 9×6 grid of cells: three columns of
// motion blocks, each block three frames wide, six motions per column (layout:
// assets/enemy/Free Sprites/imagePattern.png).
const (
	gridCols    = 9
	gridRows    = 6
	blockFrames = 3
)

// Mapping says how to bake one battler sheet: which cells make each clip.
type Mapping struct {
	Name       string                 // enemy name: output files are <Name>_<clip>.png
	Sheet      string                 // battler sheet, relative to the working directory
	FacesLeft  bool                   // the art faces left; strips are mirrored to face right
	DrawHeight int                    // output cell height in px (the size drawn in game)
	Clips      map[string]ClipMapping // standard clip name → its frames
}

// ClipMapping is one clip: battler cells in play order, plus the timing the
// printed blueprint snippet carries.
type ClipMapping struct {
	Frames   [][3]int // [motion block 0-2, row 0-5, frame 0-2] per frame
	FPS      float64
	Loop     bool
	Hold     []float64 // per-frame hold multipliers (EnemyClip.Hold)
	HitFrame int       // attack clips: the frame the strike lands on
}

type clipName struct{ name, constant string }

// clipNames are the game's standard clip names and their Go constants, in the
// game's clip order (internal/core/enemy_clips.go). The tool can't import the
// game, so keep this list in sync with clipNames there.
var clipNames = []clipName{
	{"idle", "ClipIdle"}, {"walk", "ClipWalk"}, {"alert", "ClipAlert"}, {"run", "ClipRun"},
	{"attack_1", "ClipAttack1"}, {"attack_2", "ClipAttack2"}, {"guard", "ClipGuard"},
	{"hurt", "ClipHurt"}, {"death", "ClipDeath"},
}

func (m Mapping) validate() error {
	var errs []error
	if m.Name == "" {
		errs = append(errs, errors.New("name is empty"))
	}
	if m.Sheet == "" {
		errs = append(errs, errors.New("sheet is empty"))
	}
	if m.DrawHeight <= 0 {
		errs = append(errs, errors.New("drawHeight must be > 0"))
	}
	for name, c := range m.Clips {
		if !slices.ContainsFunc(clipNames, func(n clipName) bool { return n.name == name }) {
			errs = append(errs, fmt.Errorf("clip %q: not a standard clip name", name))
		}
		if len(c.Frames) == 0 {
			errs = append(errs, fmt.Errorf("clip %q: no frames", name))
		}
		if c.FPS <= 0 {
			errs = append(errs, fmt.Errorf("clip %q: fps must be > 0", name))
		}
		if c.HitFrame < 0 || c.HitFrame >= max(len(c.Frames), 1) {
			errs = append(errs, fmt.Errorf("clip %q: hitFrame %d is past its %d frames", name, c.HitFrame, len(c.Frames)))
		}
		for i, f := range c.Frames {
			if f[0] < 0 || f[0] >= gridCols/blockFrames || f[1] < 0 || f[1] >= gridRows || f[2] < 0 || f[2] >= blockFrames {
				errs = append(errs, fmt.Errorf("clip %q frame %d: %v is outside the 3×6 blocks of 3 frames", name, i, f))
			}
		}
	}
	return errors.Join(errs...)
}

// cellRect is the pixel rect of grid cell (col, row). Cell edges are rounded
// from the exact fractional ones, so the cells tile the sheet even when its
// size doesn't divide by 9 or 6; a cell is then 1 px larger than its neighbour.
func cellRect(b image.Rectangle, col, row int) image.Rectangle {
	edge := func(lo, size, i, n int) int { return lo + int(math.Round(float64(size*i)/float64(n))) }
	return image.Rect(
		edge(b.Min.X, b.Dx(), col, gridCols), edge(b.Min.Y, b.Dy(), row, gridRows),
		edge(b.Min.X, b.Dx(), col+1, gridCols), edge(b.Min.Y, b.Dy(), row+1, gridRows),
	)
}

// frameSize is the output cell size for a sheet: DrawHeight tall, as wide as
// keeps the average cell's aspect ratio.
func frameSize(sheet image.Rectangle, drawHeight int) (w, h int) {
	cellW := float64(sheet.Dx()) / gridCols
	cellH := float64(sheet.Dy()) / gridRows
	return int(math.Round(cellW * float64(drawHeight) / cellH)), drawHeight
}

// bakeClip cuts the clip's cells out of the sheet, shrinks each to the frame
// size, mirrors them if the art faces left, and lays them out left to right.
func bakeClip(sheet image.Image, m Mapping, c ClipMapping) (*image.RGBA, error) {
	b := sheet.Bounds()
	w, h := frameSize(b, m.DrawHeight)
	if minCell := b.Dy() / gridRows; h > minCell {
		return nil, fmt.Errorf("drawHeight %d is above the %d px cell height: the tool only shrinks", h, minCell)
	}
	strip := image.NewRGBA(image.Rect(0, 0, w*len(c.Frames), h))
	for i, f := range c.Frames {
		r := cellRect(b, f[0]*blockFrames+f[2], f[1])
		cell := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
		draw.Draw(cell, cell.Bounds(), sheet, r.Min, draw.Src)
		frame := shrink(cell, w, h)
		if m.FacesLeft {
			mirror(frame)
		}
		draw.Draw(strip, image.Rect(i*w, 0, (i+1)*w, h), frame, image.Point{}, draw.Src)
	}
	return strip, nil
}

// shrink scales src down to w×h with a box filter: each output pixel is the
// average of the source area it covers, partly covered pixels weighted by how
// much of them it covers. image.RGBA is premultiplied, so averaging the
// channels directly keeps transparent edges free of dark fringes.
func shrink(src *image.RGBA, w, h int) *image.RGBA {
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	// Pass 1: rows, sw → w. Pass 2: columns, sh → h. Four float channels per pixel.
	tmp := make([]float64, w*sh*4)
	for y := range sh {
		for x := range w {
			boxSum(x, w, sw, func(sx int, weight float64) {
				p := src.PixOffset(src.Bounds().Min.X+sx, src.Bounds().Min.Y+y)
				for ch := range 4 {
					tmp[(y*w+x)*4+ch] += float64(src.Pix[p+ch]) * weight
				}
			})
		}
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var sum [4]float64
			boxSum(y, h, sh, func(sy int, weight float64) {
				for ch := range 4 {
					sum[ch] += tmp[(sy*w+x)*4+ch] * weight
				}
			})
			p := dst.PixOffset(x, y)
			for ch := range 4 {
				dst.Pix[p+ch] = uint8(math.Round(min(sum[ch], 255)))
			}
		}
	}
	return dst
}

// boxSum calls add for each source pixel that output pixel i (of n) covers
// when size source pixels map onto n, with weights summing to 1.
func boxSum(i, n, size int, add func(src int, weight float64)) {
	scale := float64(size) / float64(n)
	lo, hi := float64(i)*scale, float64(i+1)*scale
	for s := int(lo); s < size && float64(s) < hi; s++ {
		cover := min(float64(s+1), hi) - max(float64(s), lo)
		add(s, cover/scale)
	}
}

// mirror flips img left to right in place.
func mirror(img *image.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for l, r := b.Min.X, b.Max.X-1; l < r; l, r = l+1, r-1 {
			pl, pr := img.PixOffset(l, y), img.PixOffset(r, y)
			for ch := range 4 {
				img.Pix[pl+ch], img.Pix[pr+ch] = img.Pix[pr+ch], img.Pix[pl+ch]
			}
		}
	}
}

// stripName is the output file name of a clip's strip.
func stripName(m Mapping, clip string) string { return m.Name + "_" + clip + ".png" }

// snippet is the Go to paste into the enemy's blueprint: frame-size constants
// and the Clips map, in the game's clip order. outDir is the strips' folder
// relative to the repo root.
func snippet(m Mapping, outDir string, w, h int) string {
	id := strings.ToLower(m.Name[:1]) + m.Name[1:]
	var sb strings.Builder
	fmt.Fprintf(&sb, "// %s: baked by tools/battler2strip from %s\nconst (\n", m.Name, m.Sheet)
	fmt.Fprintf(&sb, "\t%sFrameWidth  = %d\n\t%sFrameHeight = %d\n)\n\n", id, w, id, h)
	sb.WriteString("Clips: map[ClipName]EnemyClip{\n")
	for _, n := range clipNames {
		c, ok := m.Clips[n.name]
		if !ok {
			continue
		}
		sheet := path.Join("..", strings.ReplaceAll(outDir, `\`, "/"), stripName(m, n.name))
		fmt.Fprintf(&sb, "\t%s: {Sheet: %q, Frames: %d, FrameW: %sFrameWidth, FrameH: %sFrameHeight, FPS: %g",
			n.constant, sheet, len(c.Frames), id, id, c.FPS)
		if c.Loop {
			sb.WriteString(", Loop: true")
		}
		if len(c.Hold) > 0 {
			holds := make([]string, len(c.Hold))
			for i, v := range c.Hold {
				holds[i] = fmt.Sprintf("%g", v)
			}
			fmt.Fprintf(&sb, ", Hold: []float64{%s}", strings.Join(holds, ", "))
		}
		if c.HitFrame > 0 {
			fmt.Fprintf(&sb, ", HitFrame: %d", c.HitFrame)
		}
		sb.WriteString("},\n")
	}
	sb.WriteString("},\n")
	return sb.String()
}

// Contact sheet look: light grey under the art, so transparent pixels and
// stray marks show, a darker line at each frame's left edge, and a gap
// between clips.
var (
	contactBackground = color.RGBA{225, 225, 225, 255}
	contactLine       = color.RGBA{160, 160, 160, 255}
)

const contactGap = 4

// contactSheet stacks the clip strips top to bottom, so every baked frame can
// be checked by eye in one image.
func contactSheet(strips []*image.RGBA, frameW int) *image.RGBA {
	w, h := 0, 0
	for _, s := range strips {
		w = max(w, s.Bounds().Dx())
		h += s.Bounds().Dy() + contactGap
	}
	sheet := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(sheet, sheet.Bounds(), &image.Uniform{contactBackground}, image.Point{}, draw.Src)
	y := 0
	for _, s := range strips {
		r := s.Bounds().Sub(s.Bounds().Min).Add(image.Pt(0, y))
		draw.Draw(sheet, r, s, s.Bounds().Min, draw.Over)
		for x := frameW; x < r.Dx(); x += frameW {
			for ly := r.Min.Y; ly < r.Max.Y; ly++ {
				sheet.SetRGBA(x, ly, contactLine)
			}
		}
		y = r.Max.Y + contactGap
	}
	return sheet
}
