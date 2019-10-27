package domain

import (
	"bufio"
	"errors"
	"github.com/grafov/m3u8"
	"os"
	"time"
)

const LIVE_START string = "2019-10-24 00:00:00 +0000"
const LIVE_START_FMT string = "2006-01-02 15:04:05 -0700"

func CurrentManifest(source string) (*m3u8.MediaPlaylist, error) {
	return CreateLiveManifest(source, time.Now())
}

func CreateLiveManifest(source string, now time.Time) (*m3u8.MediaPlaylist, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		return nil, err
	}
	switch listType {
	case m3u8.MEDIA:
		mediaPL := p.(*m3u8.MediaPlaylist)
		// Convert VOD playlist to Live playlist.
		if err := mediaPL.SetWinSize(mediaPL.Count()); err != nil {
			return nil, err
		}
		mediaPL.Closed = false
		shiftedPL, err := shiftLiveManifest(mediaPL, now)
		if err != nil {
			return nil, err
		}
		return shiftedPL, nil
	case m3u8.MASTER:
		return nil, errors.New("master playlist not supported")
	}
	return nil, nil
}

func shiftLiveManifest(mediaPL *m3u8.MediaPlaylist, now time.Time) (*m3u8.MediaPlaylist, error) {
	startTime, err := time.Parse(LIVE_START_FMT, LIVE_START)
	if err != nil {
		return nil, err
	}

	// Difference between live start time and current time.
	// It is cast to milli sec.
	timeDiff := uint64(now.Unix() - startTime.Unix()) * 1000

	// Sum total time and cast to milli sec.
	totalTime := uint64(totalTime(mediaPL) * 1000)

	// Total count of manifest creation.
	createCount := timeDiff / totalTime * uint64(mediaPL.Count())

	// Shift Time for shifted media playlist.
	shiftTime := float64(timeDiff % totalTime) / 1000

	mediaPL.SeqNo = createCount

	var sumDuration float64
	for i := 0; i < int(mediaPL.Count()); i++ {
		sumDuration += mediaPL.Segments[i].Duration
		if sumDuration > shiftTime {
			break
		}
		mediaPL.Slide(mediaPL.Segments[i].URI, mediaPL.Segments[i].Duration, mediaPL.Segments[i].Title)
	}

	return mediaPL, nil
}

func totalTime(mediaPL *m3u8.MediaPlaylist) float64 {
	var sumDuration float64
	for i := 0; i < int(mediaPL.Count()); i++ {
		sumDuration += mediaPL.Segments[i].Duration
	}
	return sumDuration
}
