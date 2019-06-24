package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"log"
	"net/http"
	"time"

	"github.com/faiface/beep/mp3"

	. "fyne.io/fyne/widget"
)

var (
	w             fyne.Window
	playing       = false
	volume        *effects.Volume
	initialVolume = -2.0

	testDisplay *ProgressBar

	radioStations = map[string]string{
		"HouseTime.FM":  "http://lw1.mp3.tb-group.fm/ht.mp3",
		"TranceBase.FM": "http://lw2.mp3.tb-group.fm/trb.mp3",
		"TechnoBase.FM": "http://lw1.mp3.tb-group.fm/tb.mp3",
		"HardBase.FM":   "http://lw3.mp3.tb-group.fm/hb.mp3",
		"CoreTime.FM":   "http://lw3.mp3.tb-group.fm/ct.mp3",
		"ClubTime.FM":   "http://lw3.mp3.tb-group.fm/clt.mp3",
		"TeaTime.FM":    "http://lw3.mp3.tb-group.fm/tt.mp3",
	}
)

func main() {
	application := app.New()
	w = application.NewWindow("Player")

	stationPicker := NewRadio(getKeys(radioStations), func(selected string) {
		log.Print(selected)

		go initAudio(radioStations[selected])
	})

	stationPicker.SetSelected("HouseTime.FM")

	var playPauseButton *Button
	playPauseButton = NewButton("Play", func() {
		playing = !playing

		if playing {
			playPauseButton.SetText("Pause")
		} else {
			playPauseButton.SetText("Play")
		}
	})

	volumeArea := initVolumeArea()

	testDisplay = NewProgressBar()
	testDisplay.Min = 0.0
	testDisplay.Max = 0.2

	w.SetContent(
		NewVBox(
			NewHBox(
				NewLabel("Currently Playing: "),
				stationPicker,
			),
			playPauseButton,
			volumeArea,
			//testDisplay,
		),
	)

	go initAudio(radioStations[stationPicker.Selected])

	w.ShowAndRun()
}

func initVolumeArea() *Box {
	volumeDisplay := NewProgressBar()
	volumeDisplay.Min = -6.0
	volumeDisplay.Max = 0.0
	volumeDisplay.Value = initialVolume

	stepSize := 0.5

	decreaseVolumeButton := NewButton("-", func() {
		volume.Volume -= stepSize
		volumeDisplay.SetValue(volume.Volume)

		log.Printf("Current Volume: %v", volume.Volume)
	})
	increaseVolumeButton := NewButton("+", func() {
		volume.Volume += stepSize
		volumeDisplay.SetValue(volume.Volume)

		log.Printf("Current Volume: %v", volume.Volume)
	})
	volumeArea := NewHBox(
		decreaseVolumeButton,
		volumeDisplay,
		increaseVolumeButton,
	)
	return volumeArea
}

func initAudio(url string) {
	resp, err := http.Get(url)
	check(err)
	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	check(err)

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume = &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   initialVolume,
		Silent:   false,
	}

	viz := &VisualizerStreamer{Streamer: volume}

	done := make(chan bool)
	speaker.Play(beep.Seq(viz, beep.Callback(func() {
		done <- true
	})))

	for {
		speaker.Lock()
		ctrl.Paused = !playing
		speaker.Unlock()
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func getKeys(in map[string]string) (keys []string) {
	keys = make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}

	return
}
