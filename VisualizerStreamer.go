package main

import (
	"github.com/faiface/beep"
)

type VisualizerStreamer struct {
	Streamer       beep.Streamer
	lastGainValues []float64
}

const averageGainSampleAmount = 30

func (v *VisualizerStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.Streamer.Stream(samples)

	highestGain := 0.0

	for i := range samples[:n] {
		if samples[i][0] > highestGain {
			highestGain = samples[i][0]
		}
	}

	v.lastGainValues = append(v.lastGainValues, highestGain)
	if len(v.lastGainValues) > averageGainSampleAmount {
		v.lastGainValues = v.lastGainValues[len(v.lastGainValues)-averageGainSampleAmount:]
	}

	testDisplay.SetValue(getAverage(v.lastGainValues))

	return n, ok
}

func (v *VisualizerStreamer) Err() error {
	return v.Streamer.Err()
}

func getAverage(list []float64) float64 {
	tmp := 0.0
	for _, n := range list {
		tmp += n
	}

	return tmp / float64(len(list))
}
