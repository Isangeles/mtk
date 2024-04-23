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
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

// Struct for audio player.
type AudioPlayer struct {
	playlist []*beep.Buffer
	playID   int
	control  *beep.Ctrl
	volume   *effects.Volume
}

// NewAudioPlayer creates new audio player.
func NewAudioPlayer() *AudioPlayer {
	ap := new(AudioPlayer)
	ap.playlist = make([]*beep.Buffer, 0)
	ap.control = new(beep.Ctrl)
	ap.volume = &effects.Volume{
		Streamer: ap.control,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	return ap
}

// AddAudio adds specified audio stream to the current playlist.
func (ap *AudioPlayer) AddAudio(ab *beep.Buffer) error {
	ap.playlist = append(ap.playlist, ab)
	return nil
}

// Playlist returns the audio player playlist.
func (ap *AudioPlayer) Playlist() []*beep.Buffer {
	return ap.playlist
}

// SetPlaylist sets specified slice with audio streams
// as player playlist.
func (ap *AudioPlayer) SetPlaylist(playlist []*beep.Buffer) {
	ap.playlist = playlist
}

// ResumePlaylist starts playing audio from the playlist
// for current playlist ID.
func (ap *AudioPlayer) ResumePlaylist() {
	if ap.playID < 0 || ap.playID > len(ap.playlist)-1 {
		return
	}
	buffer := ap.playlist[ap.playID]
	ap.Play(buffer)
	return
}

// Play starts playing specified audio stream.
func (ap *AudioPlayer) Play(buffer *beep.Buffer) {
	streamer := buffer.Streamer(0, buffer.Len())
	ap.control.Streamer = streamer
	speaker.Play(ap.volume)
}

// Stop stops player.
func (ap *AudioPlayer) Stop() {
	if ap.control.Streamer == nil {
		return
	}
	speaker.Lock()
	ap.control.Streamer = nil
	speaker.Unlock()
}

// Reset stops player and moves play index to
// first music playlist index.
func (ap *AudioPlayer) Reset() {
	ap.Stop()
	ap.SetPlayIndex(0)
}

// SetVolume sets specified value as current
// audio volume value.
// 0 - unmodified(system volume), > 0 - lauder, < 0 quieter.
func (ap *AudioPlayer) SetVolume(v float64) {
	ap.volume.Volume = v
}

// Volume returns current volume value.
// 0 - unmodified(system volume), > 0 - lauder, < 0 quieter.
func (ap *AudioPlayer) Volume() float64 {
	return ap.volume.Volume
}

// SetMute mutes/unmutes audio player.
func (ap *AudioPlayer) SetMute(m bool) {
	ap.volume.Silent = m
}

// Muted checks if audio player is muted.
func (ap *AudioPlayer) Muted() bool {
	return ap.volume.Silent
}

// Playing checks if audio player is playing
// any audio buffer.
func (ap *AudioPlayer) Playing() bool {
	return ap.control.Streamer != nil
}

// Clear clears music playlist.
func (ap *AudioPlayer) Clear() {
	ap.playlist = make([]*beep.Buffer, 0)
}

// PlayIndex returns index of currently playing audio
// buffer from the playlist.
func (ap *AudioPlayer) PlayIndex() int {
	return ap.playID
}

// SetPlayIndex sets specified index as current index
// on music playlist.
// If specified value is bigger than playlist lenght
// then first index is set, if is lower than 0 then
// last index is set.
func (ap *AudioPlayer) SetPlayIndex(id int) {
	switch {
	case id > len(ap.playlist)-1:
		ap.playID = 0
	case id < 0:
		ap.playID = len(ap.playlist) - 1
	default:
		ap.playID = id
	}
}
