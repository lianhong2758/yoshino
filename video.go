package yoshino

import (
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"sync"

	"github.com/gen2brain/mpeg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type VideoPlayer struct {
	Mpg *mpeg.MPEG
	// yCbCrImage is the current frame image in YCbCr format.
	// An MPEG frame is stored in this image first.
	// Then, this image data is converted to RGB to frameImage.
	yCbCrImage *ebiten.Image

	// yCbCrBytes is the byte slice to store YCbCr data.
	// This includes Y, Cb, Cr, and alpha (always 0xff) data for each pixel.
	yCbCrBytes []byte
	// yCbCrShader is the shader to convert YCbCr to RGB.
	yCbCrShader *ebiten.Shader
	// frameImage is the current frame image in RGB format.
	frameImage *ebiten.Image

	closeOnce sync.Once

	src io.ReadCloser
}

// updateFrame upadtes the current video frame.
func (v *VideoPlayer) updateFrame(musicplayer *audio.Player) error {
	pos := musicplayer.Position().Seconds()

	video := v.Mpg.Video()
	if video.HasEnded() {
		v.frameImage.Clear()
		var err error
		v.closeOnce.Do(func() {
			fmt.Println("The video has ended.")
			if err1 := v.src.Close(); err1 != nil {
				err = err1
			}
		})
		return err
	}

	d := 1 / v.Mpg.Framerate()
	var mpegFrame *mpeg.Frame
	for video.Time()+d <= pos && !video.HasEnded() {
		mpegFrame = video.Decode()
	}

	if mpegFrame == nil {
		return nil
	}

	img := mpegFrame.YCbCr()
	if img.SubsampleRatio != image.YCbCrSubsampleRatio420 {
		return fmt.Errorf("video: subsample ratio must be 4:2:0")
	}
	w, h := v.Mpg.Width(), v.Mpg.Height()
	for j := 0; j < h; j++ {
		yi := j * img.YStride
		ci := (j / 2) * img.CStride
		// Create temporary slices to encourage BCE (boundary-checking elimination).
		ys := img.Y[yi : yi+w]
		cbs := img.Cb[ci : ci+w/2]
		crs := img.Cr[ci : ci+w/2]
		for i := 0; i < w; i++ {
			idx := 4 * (j*w + i)
			buf := v.yCbCrBytes[idx : idx+3]
			buf[0] = ys[i]
			buf[1] = cbs[i/2]
			buf[2] = crs[i/2]
			// p.yCbCrBytes[3] = 0xff is not needed as the shader ignores this part.
		}
	}

	v.yCbCrImage.WritePixels(v.yCbCrBytes)

	// Converting YCbCr to RGB on CPU is slow. Use a shader instead.
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = v.yCbCrImage
	op.Blend = ebiten.BlendCopy
	v.frameImage.DrawRectShader(w, h, v.yCbCrShader, op)

	return nil
}
func (v *VideoPlayer) Draw(screen *ebiten.Image, musicplayer *audio.Player) {
	if err := v.updateFrame(musicplayer); err != nil {
		log.Println("Error:", err)
		return
	}

	frame := v.frameImage
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	fw, fh := frame.Bounds().Dx(), frame.Bounds().Dy()

	op := ebiten.DrawImageOptions{}
	wf, hf := float64(sw)/float64(fw), float64(sh)/float64(fh)
	s := wf
	if hf < wf {
		s = hf
	}
	op.GeoM.Scale(s, s)

	offsetX, offsetY := float64(screen.Bounds().Min.X), float64(screen.Bounds().Min.Y)
	op.GeoM.Translate(offsetX+(float64(sw)-float64(fw)*s)/2, offsetY+(float64(sh)-float64(fh)*s)/2)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(frame, &op)
}

type mpegAudio struct {
	audio *mpeg.Audio

	// leftovers is the remaining audio samples of the previous Read call.
	leftovers []byte

	// m is the mutex shared with the mpegPlayer.
	// As *mpeg.MPEG is not concurrent safe, this mutex is necessary.
	//m *sync.Mutex
}

func (a *mpegAudio) Read(buf []byte) (int, error) {
	// a.m.Lock()
	// defer a.m.Unlock()

	var readBytes int
	if len(a.leftovers) > 0 {
		n := copy(buf, a.leftovers)
		readBytes += n
		buf = buf[n:]

		copy(a.leftovers, a.leftovers[n:])
		a.leftovers = a.leftovers[:len(a.leftovers)-n]
	}

	for len(buf) > 0 && !a.audio.HasEnded() {
		mpegSamples := a.audio.Decode()
		if mpegSamples == nil {
			break
		}

		bs := make([]byte, len(mpegSamples.Interleaved)*4)
		for i, s := range mpegSamples.Interleaved {
			v := math.Float32bits(s)
			bs[4*i] = byte(v)
			bs[4*i+1] = byte(v >> 8)
			bs[4*i+2] = byte(v >> 16)
			bs[4*i+3] = byte(v >> 24)
		}

		n := copy(buf, bs)
		readBytes += n
		buf = buf[n:]

		if n < len(bs) {
			a.leftovers = append(a.leftovers, bs[n:]...)
			break
		}
	}

	if a.audio.HasEnded() {
		return readBytes, io.EOF
	}
	return readBytes, nil
}
