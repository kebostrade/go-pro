package structural

import "fmt"

/*
ADAPTER PATTERN

Purpose: Convert the interface of a class into another interface clients expect.

Use Cases:
- Integrating third-party libraries
- Legacy code integration
- API versioning
- Data format conversion

Go-Specific Implementation:
- Wrapper struct implementing target interface
- Composition over inheritance
*/

// MediaPlayer is the target interface
type MediaPlayer interface {
	Play(filename string) error
}

// MP3Player is our existing implementation
type MP3Player struct{}

func (m *MP3Player) PlayMP3(filename string) error {
	fmt.Printf("🎵 Playing MP3 file: %s\n", filename)
	return nil
}

// VLCPlayer is a third-party player (adaptee)
type VLCPlayer struct{}

func (v *VLCPlayer) PlayVLC(filename string) error {
	fmt.Printf("🎬 Playing VLC file: %s\n", filename)
	return nil
}

// VLCAdapter adapts VLCPlayer to MediaPlayer interface
type VLCAdapter struct {
	vlc *VLCPlayer
}

func NewVLCAdapter() *VLCAdapter {
	return &VLCAdapter{
		vlc: &VLCPlayer{},
	}
}

func (a *VLCAdapter) Play(filename string) error {
	return a.vlc.PlayVLC(filename)
}

// MP4Player is another third-party player
type MP4Player struct{}

func (m *MP4Player) PlayMP4(filename string) error {
	fmt.Printf("📹 Playing MP4 file: %s\n", filename)
	return nil
}

// MP4Adapter adapts MP4Player to MediaPlayer interface
type MP4Adapter struct {
	mp4 *MP4Player
}

func NewMP4Adapter() *MP4Adapter {
	return &MP4Adapter{
		mp4: &MP4Player{},
	}
}

func (a *MP4Adapter) Play(filename string) error {
	return a.mp4.PlayMP4(filename)
}

