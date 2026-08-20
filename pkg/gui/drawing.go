package gui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var (
	ColorBgDark        = color.RGBA{18, 20, 24, 255}
	ColorPanelBg       = color.RGBA{28, 32, 38, 255}
	ColorPanelBorder   = color.RGBA{45, 52, 62, 255}
	ColorLightSquare   = color.RGBA{235, 236, 212, 255}
	ColorDarkSquare    = color.RGBA{118, 150, 86, 255}
	ColorSelected      = color.RGBA{246, 246, 105, 200}
	ColorValidMove     = color.RGBA{80, 200, 120, 200}
	ColorMandatoryCap  = color.RGBA{230, 80, 80, 220}
	ColorLastMove      = color.RGBA{205, 210, 106, 180}
	ColorWhitePiece    = color.RGBA{245, 245, 240, 255}
	ColorWhitePieceRim = color.RGBA{190, 190, 185, 255}
	ColorBlackPiece    = color.RGBA{40, 40, 44, 255}
	ColorBlackPieceRim = color.RGBA{70, 70, 78, 255}
	ColorGoldCrown     = color.RGBA{255, 215, 0, 255}
	ColorSilverCrown   = color.RGBA{220, 220, 230, 255}
	ColorTextLight     = color.RGBA{240, 240, 245, 255}
	ColorTextDim       = color.RGBA{160, 165, 175, 255}
	ColorBtnNormal     = color.RGBA{40, 90, 160, 255}
	ColorBtnHover      = color.RGBA{55, 115, 200, 255}
	ColorBtnActive     = color.RGBA{30, 140, 90, 255}
	ColorAccentYellow  = color.RGBA{245, 195, 65, 255}
	ColorAccentGreen   = color.RGBA{80, 210, 130, 255}
	ColorAccentRed     = color.RGBA{235, 85, 85, 255}
)

var emptyImage = func() *ebiten.Image {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	return img
}()

func DrawRect(dst *ebiten.Image, x, y, w, h float32, clr color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w)/3.0, float64(h)/3.0)
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	dst.DrawImage(emptyImage, op)
}

func DrawRoundedRect(dst *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	if r <= 0 {
		DrawRect(dst, x, y, w, h, clr)
		return
	}
	DrawRect(dst, x+r, y, w-2*r, h, clr)
	DrawRect(dst, x, y+r, w, h-2*r, clr)
	DrawCircle(dst, x+r, y+r, r, clr)
	DrawCircle(dst, x+w-r, y+r, r, clr)
	DrawCircle(dst, x+r, y+h-r, r, clr)
	DrawCircle(dst, x+w-r, y+h-r, r, clr)
}

func DrawCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color) {
	if r <= 0 {
		return
	}
	segments := int(math.Max(16, float64(r*2.5)))
	var vertices []ebiten.Vertex
	var indices []uint16

	r32, g32, b32, a32 := clr.RGBA()
	rf := float32(r32) / 65535.0
	gf := float32(g32) / 65535.0
	bf := float32(b32) / 65535.0
	af := float32(a32) / 65535.0

	vertices = append(vertices, ebiten.Vertex{
		DstX:   cx,
		DstY:   cy,
		SrcX:   1.5,
		SrcY:   1.5,
		ColorR: rf,
		ColorG: gf,
		ColorB: bf,
		ColorA: af,
	})

	for i := 0; i <= segments; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(segments)
		vx := cx + float32(math.Cos(theta))*r
		vy := cy + float32(math.Sin(theta))*r

		vertices = append(vertices, ebiten.Vertex{
			DstX:   vx,
			DstY:   vy,
			SrcX:   1.5,
			SrcY:   1.5,
			ColorR: rf,
			ColorG: gf,
			ColorB: bf,
			ColorA: af,
		})

		if i > 0 {
			indices = append(indices, 0, uint16(i), uint16(i+1))
		}
	}

	dst.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
}

func DrawRing(dst *ebiten.Image, cx, cy, r, thickness float32, clr color.Color) {
	if r <= 0 || thickness <= 0 {
		return
	}
	outerR := r + thickness/2.0
	innerR := math.Max(0.0, float64(r-thickness/2.0))

	segments := int(math.Max(20, float64(r*3.0)))
	var vertices []ebiten.Vertex
	var indices []uint16

	r32, g32, b32, a32 := clr.RGBA()
	rf := float32(r32) / 65535.0
	gf := float32(g32) / 65535.0
	bf := float32(b32) / 65535.0
	af := float32(a32) / 65535.0

	for i := 0; i <= segments; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(segments)
		cosT := float32(math.Cos(theta))
		sinT := float32(math.Sin(theta))

		ox := cx + cosT*outerR
		oy := cy + sinT*outerR

		ix := cx + cosT*float32(innerR)
		iy := cy + sinT*float32(innerR)

		vertices = append(vertices,
			ebiten.Vertex{DstX: ox, DstY: oy, SrcX: 1.5, SrcY: 1.5, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
			ebiten.Vertex{DstX: ix, DstY: iy, SrcX: 1.5, SrcY: 1.5, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af},
		)

		if i > 0 {
			base := uint16((i - 1) * 2)
			indices = append(indices, base, base+1, base+2, base+1, base+3, base+2)
		}
	}

	dst.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
}

func DrawStar(dst *ebiten.Image, cx, cy, r float32, clr color.Color) {
	points := 5
	innerR := r * 0.45
	var vertices []ebiten.Vertex
	var indices []uint16

	r32, g32, b32, a32 := clr.RGBA()
	rf := float32(r32) / 65535.0
	gf := float32(g32) / 65535.0
	bf := float32(b32) / 65535.0
	af := float32(a32) / 65535.0

	vertices = append(vertices, ebiten.Vertex{DstX: cx, DstY: cy, SrcX: 1.5, SrcY: 1.5, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af})

	totalPoints := points * 2
	for i := 0; i <= totalPoints; i++ {
		theta := float64(i)*math.Pi/float64(points) - math.Pi/2.0
		radius := r
		if i%2 != 0 {
			radius = innerR
		}

		vx := cx + float32(math.Cos(theta))*radius
		vy := cy + float32(math.Sin(theta))*radius

		vertices = append(vertices, ebiten.Vertex{DstX: vx, DstY: vy, SrcX: 1.5, SrcY: 1.5, ColorR: rf, ColorG: gf, ColorB: bf, ColorA: af})
		if i > 0 {
			indices = append(indices, 0, uint16(i), uint16(i+1))
		}
	}

	dst.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
}

func DrawText(dst *ebiten.Image, str string, x, y int, clr color.Color) {
	text.Draw(dst, str, basicfont.Face7x13, x, y, clr)
}
