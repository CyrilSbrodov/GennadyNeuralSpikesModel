package encoder

import (
	"image"
	"image/color"

	"gennady-neural-spikes-model/internal/engine"
)

func EncodeImageToSpikes(img image.Image, t0 uint64, cfg EncoderConfig) []engine.TimedSpikeEvent {
	if img == nil {
		return nil
	}
	if err := cfg.Validate(); err != nil {
		return nil
	}

	resized := resizeNearest(img, cfg.ImageWidth, cfg.ImageHeight)
	channelsPerPixel := 1
	if cfg.UseRGB {
		channelsPerPixel = 3
	}

	events := make([]engine.TimedSpikeEvent, 0, cfg.ImageWidth*cfg.ImageHeight)
	unitIdx := 0
	for y := 0; y < cfg.ImageHeight; y++ {
		for x := 0; x < cfg.ImageWidth; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			vals := []float32{float32(r) / 65535.0}
			if cfg.UseRGB {
				vals = []float32{
					float32(r) / 65535.0,
					float32(g) / 65535.0,
					float32(b) / 65535.0,
				}
			}

			windowStart := t0 + uint64(unitIdx)*cfg.UnitWindowTicks
			for ci := 0; ci < channelsPerPixel; ci++ {
				channel := (unitIdx*channelsPerPixel + ci) % cfg.InputChannels
				events = append(events, buildTemporalSpikes(windowStart, channel, vals[ci], cfg)...)
			}
			unitIdx++
		}
	}

	return events
}

func resizeNearest(src image.Image, width, height int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	sb := src.Bounds()
	sw := sb.Dx()
	sh := sb.Dy()

	for y := 0; y < height; y++ {
		sy := sb.Min.Y + y*sh/height
		for x := 0; x < width; x++ {
			sx := sb.Min.X + x*sw/width
			c := color.NRGBAModel.Convert(src.At(sx, sy)).(color.NRGBA)
			dst.SetNRGBA(x, y, c)
		}
	}

	return dst
}
