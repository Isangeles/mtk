/*
 * messagequeue.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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
	"github.com/gopxl/pixel"
)

// MessageQueue struct for list with messages to display.
type MessageQueue struct {
	queue []*MessageWindow
	focus *Focus
}

// NewMessageQueue creates new messages queue.
func NewMessageQueue(focus *Focus) *MessageQueue {
	mq := new(MessageQueue)
	mq.queue = make([]*MessageWindow, 0)
	mq.focus = focus
	return mq
}

// Draw draws all messages
func (mq *MessageQueue) Draw(t pixel.Target, matrix pixel.Matrix) {
	for _, m := range mq.queue {
		if m.Opened() {
			m.Draw(t, matrix)
		}
	}
}

// Update updates all messages in queue.
func (mq *MessageQueue) Update(win *Window) {
	for i, m := range mq.queue {
		if m.Opened() {
			if i == len(mq.queue)-1 {
				m.Active(true)
				mq.focus.Focus(m)
			} else {
				m.Active(false)
			}
			m.Update(win)
		}
		if m.Dismissed() {
			mq.Remove(i)
		}
	}
}

// Append adds specified message to the front of queue.
func (mq *MessageQueue) Append(m *MessageWindow) {
	mq.queue = append(mq.queue, m)
}

// Remove removes message with specified index from queue.
func (mq *MessageQueue) Remove(i int) {
	mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
}

// ContainsPosition checks whether specified position is
// contained by any message window in the queue.
func (mq *MessageQueue) ContainsPosition(pos pixel.Vec) bool {
	for _, msg := range mq.queue {
		if msg.DrawArea().Contains(pos) {
			return true
		}
	}
	return false
}
