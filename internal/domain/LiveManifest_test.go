package domain

import (
	"bufio"
	"errors"
	"github.com/grafov/m3u8"
	"os"
	"testing"
	"time"
)

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func originalLiveManifest() *m3u8.MediaPlaylist {
	dir := getCurrentDir()
	f, err := os.Open(dir + "/video.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediaPL := p.(*m3u8.MediaPlaylist)
		if err := mediaPL.SetWinSize(mediaPL.Count()); err != nil {
			panic(err)
		}
		mediaPL.Closed = false
		return mediaPL
	case m3u8.MASTER:
		panic(errors.New("master playlist not supported"))
	}
	return nil
}

func TestCreateLiveManifest(t *testing.T) {
	dir := getCurrentDir()
	mediaPL := originalLiveManifest()

	type args struct {
		source string
		now    time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *m3u8.MediaPlaylist
		wantErr bool
	}{
		{
			"Check return.",
			args{dir + "/video.m3u8",
				time.Date(2019, 10, 24, 0, 0, 0, 0, time.UTC)},
			mediaPL,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateLiveManifest(tt.args.source, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLiveManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.String() != tt.want.String() {
					t.Errorf("CreateLiveManifest() got = %v, want %v", got, tt.want)
				}
			} else if got != tt.want {
				t.Errorf("CreateLiveManifest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shiftLiveManifest(t *testing.T) {
	mediaPL := originalLiveManifest()
	shiftedMediaPL := originalLiveManifest()
	shiftedMediaPL.Slide(shiftedMediaPL.Segments[0].URI, shiftedMediaPL.Segments[0].Duration, shiftedMediaPL.Segments[0].Title)
	shiftedMediaPLOneHour := originalLiveManifest()
	shiftedMediaPLOneHour.SeqNo = 700
	for i := 0; i < 19; i++ {
		shiftedMediaPLOneHour.Slide(shiftedMediaPL.Segments[i].URI, shiftedMediaPL.Segments[i].Duration, shiftedMediaPL.Segments[i].Title)
	}

	type args struct {
		mediaPL *m3u8.MediaPlaylist
		now     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *m3u8.MediaPlaylist
		wantErr bool
	}{
		{
			"Shift 5 seconds.",
			args{originalLiveManifest(),
				time.Date(2019, 10, 24, 0, 0, 5, 0, time.UTC)},
			mediaPL,
			false,
		},
		{
			"Shift 6 seconds.",
			args{originalLiveManifest(),
				time.Date(2019, 10, 24, 0, 0, 6, 0, time.UTC)},
			shiftedMediaPL,
			false,
		},
		{
			"Shift one hour.",
			args{originalLiveManifest(),
				time.Date(2019, 10, 24, 1, 0, 0, 0, time.UTC)},
			shiftedMediaPLOneHour,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shiftLiveManifest(tt.args.mediaPL, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("shiftLiveManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.String() != tt.want.String() {
					t.Errorf("CreateLiveManifest() got = %v, want %v", got, tt.want)
				}
			} else if got != tt.want {
				t.Errorf("CreateLiveManifest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_totalTime(t *testing.T) {
	mediaPL := originalLiveManifest()
	type args struct {
		mediaPL *m3u8.MediaPlaylist
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"Sum total time of segments.",
			args{mediaPL},
			139.99000199999998,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := totalTime(tt.args.mediaPL); got != tt.want {
				t.Errorf("totalTime() = %v, want %v", got, tt.want)
			}
		})
	}
}