/*
 * audioplayer.go
 *
 * Copyright 2019-2024 Dariusz Sikora <ds@isangeles.dev>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301, USA.
 *
 *
 */

package mtk

import (
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// Struct for audio player.
type AudioPlayer struct {
	playlist []beep.Streamer
	playID   int
	mixer    *beep.Mixer
	control  *beep.Ctrl
	volume   *effects.Volume
}

// NewAudioPlayer creates new audio player for specified
// stream format.
func NewAudioPlayer(format beep.Format) (*AudioPlayer, error) {
	p := new(AudioPlayer)
	p.playlist = make([]beep.Streamer, 0)
	p.mixer = new(beep.Mixer)
	p.control = new(beep.Ctrl)
	p.volume = &effects.Volume{
		Streamer: p.control,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize speaker: %v", err)
	}
	speaker.Play(p.mixer)
	return p, nil
}

// AddAudio adds specified audio stream to playlist.
func (p *AudioPlayer) AddAudio(ab *beep.Buffer) error {
	s := ab.Streamer(0, ab.Len())
	p.playlist = append(p.playlist, s)
	return nil
}

// SetPlaylist sets specified slice with audio streams
// as player playlist.
func (p *AudioPlayer) SetPlaylist(playlist []beep.Streamer) {
	p.playlist = playlist
}

// Play starts player.
func (p *AudioPlayer) ResumePlaylist() error {
	if p.playID < 0 || p.playID > len(p.playlist)-1 {
		return fmt.Errorf("audio_player: current playlist position nil")
	}
	m := p.playlist[p.playID]
	p.volume.Streamer = m
	p.mixer.Add(p.volume)
	return nil
}

// Play starts playing specified audio stream.
func (p *AudioPlayer) Play(ab *beep.Buffer) error {
	s := ab.Streamer(0, ab.Len())
	p.mixer.Add(s)
	return nil
}

// Stop stops player.
func (p *AudioPlayer) StopPlaylist() {
	if p.control.Streamer == nil {
		return
	}
	speaker.Lock()
	p.control.Streamer = nil
	speaker.Unlock()
}

// Reset stops player and moves play index to
// first music playlist index.
func (p *AudioPlayer) Reset() {
	p.StopPlaylist()
	p.SetPlayIndex(0)
}

// Next moves play index to next position
// on music playlist.
func (p *AudioPlayer) Next() {
	p.StopPlaylist()
	p.SetPlayIndex(p.playID + 1)
	p.ResumePlaylist()
}

// Prev moves play index to previous position
// on music playlist.
func (p *AudioPlayer) Prev() {
	p.StopPlaylist()
	p.SetPlayIndex(p.playID - 1)
	p.ResumePlaylist()
}

// SetVolume sets specified value as current
// value.
// 0 - unmodified, > 0 - lauder, < 0 quieter.
func (ap *AudioPlayer) SetVolume(v float64) {
	ap.volume.Volume = v
}

// Volume returns current volume value.
// 0 - unmodified, > 0 - lauder, < 0 quieter.
func (ap *AudioPlayer) Volume() float64 {
	return ap.volume.Volume
}

// SetMute mutes/unmutes audio player.
func (ap *AudioPlayer) SetMute(m bool) {
	ap.volume.Silent = m
}

// Muted check if audio player is muted.
func (ap *AudioPlayer) Muted() bool {
	return ap.volume.Silent
}

// Clear clears music playlist.
func (p *AudioPlayer) Clear() {
	p.playlist = make([]beep.Streamer, 0)
}

// SetPlayIndex sets specified index as current index
// on music playlist.
// If specified value is bigger than playlist lenght
// then first index is set, if is lower than 0 then
// last index is set.
func (p *AudioPlayer) SetPlayIndex(id int) {
	switch {
	case id > len(p.playlist)-1:
		p.playID = 0
	case id < 0:
		p.playID = len(p.playlist) - 1
	default:
		p.playID = id
	}
}
