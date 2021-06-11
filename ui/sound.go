package ui

import (
	"github.com/gordonklaus/portaudio"
)

func PortAudioInitialize() {
	portaudio.Initialize()
}

func PortAudioTerminate() {
	portaudio.Terminate()
}

func OpenAudioStream(sampleBuffer chan float32) (*portaudio.Stream, float64, error) {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, -1, err
	}

	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	stream, err := portaudio.OpenStream(parameters, func(out []float32) {
		var output float32
		for i := range out {
			if i%parameters.Output.Channels == 0 {
				select {
				case sample := <-sampleBuffer:
					output = sample
				default:
					output = 0
				}
			}
			out[i] = output
		}
	})
	if err != nil {
		return nil, -1, err
	}

	if err := stream.Start(); err != nil {
		return nil, -1, err
	}
	return stream, parameters.SampleRate, nil
}
