package controller

import (
	"github.com/klauspost/reedsolomon"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func Encode(filename string, shards []string, n, k int) error {
	log.Debugf("filename: %s, shards: %v, n: %d, k: %d", filename, shards, n, k)

	// open file
	file, err := os.Open(filename)
	if err != nil {
		log.WithError(err).Errorf("Open %s failed.", filename)
		return err
	}
	defer file.Close()

	// Create encoding matrix
	enc, err := reedsolomon.NewStream(n, k)
	if err != nil {
		log.WithError(err).Errorf("Create new encoder(n=%d,k=%d) failed.", n, k)
		return err
	}

	// Create the resulting files
	out := make([]*os.File, n+k)
	for i := range out {
		out[i], err = os.Create(shards[i])
		if err != nil {
			log.WithError(err).Errorf("os.Create(%s) failed.", shards[i])
			return err
		}
	}

	// Split into files.
	instat, err := file.Stat()
	if err != nil {
		log.WithError(err).Errorf("file.Stat() failed: %s.", file.Name())
		return err
	}
	data := make([]io.Writer, n)
	for i := range data {
		data[i] = out[i]
	}
	err = enc.Split(file, data, instat.Size())
	if err != nil {
		log.WithError(err).Errorf("Split file %v(%vB) failed.", file.Name(), instat.Size())
		return err
	}

	// Close and re-open the files.
	input := make([]io.Reader, n)
	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		if err != nil {
			log.WithError(err).Errorf("Open file %s failed.", out[i].Name())
			return err
		}
		input[i] = f
		out[i] = f
	}

	// Create parity output writers
	parity := make([]io.Writer, k)
	for i := range parity {
		parity[i] = out[n+i]
	}

	// Encode parity
	err = enc.Encode(input, parity)
	if err != nil {
		log.WithError(err).Errorf("Encode parity shards failed.")
	}

	// Close result files
	for _, f := range out {
		f.Close()
	}

	return nil
}

func decode(filename string, size int64, shards []string, n, k int) error {
	log.Debugf("filename: %s, shards: %v, n: %d, k: %d", filename, shards, n, k)

	// read shards
	inputs := make([]io.Reader, n+k)
	for i, s := range shards {
		f, err := os.Open(s)
		if err != nil {
			log.WithError(err).Errorf("Open file %s failed.", s)
			return err
		}
		inputs[i] = f
		defer f.Close()
	}

	// create file
	file, err := os.Create(filename)
	if err != nil {
		log.WithError(err).Errorf("Create file %s failed.", filename)
		return err
	}
	defer file.Close()

	enc, err := reedsolomon.NewStream(n, k)
	if err != nil {
		log.WithError(err).Errorf("Create new encoder failed.")
		return err
	}

	// ok, err := enc.Verify(inputs)
	// logrus.WithError(err).Info(ok)

	err = enc.Join(file, inputs, size)
	if err != nil {
		log.WithError(err).Error("reconsruct failed")
		return err
	}

	return nil
}
