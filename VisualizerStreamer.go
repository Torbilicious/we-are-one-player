package main

import (
	"github.com/faiface/beep"
	"math"
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
		highestGain = math.Max(samples[i][0], highestGain)
		highestGain = math.Max(samples[i][1], highestGain)
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
	sum := 0.0
	for _, n := range list {
		sum += n
	}

	return sum / float64(len(list))
}
