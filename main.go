package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultThreadCount        = 2
	defaultHeight             = 480
	defaultWidth              = 640
	defaultCompressionQuality = 95
	defaultBlur               = 1
	defaultLoggingLevel       = "info"

	retOk int = 0
	retKo int = 1
)

var (
	threadCount        uint
	height             uint
	width              uint
	compressionQuality uint
	blur               float64
	input              string
	output             string
	loggingLevel       string
)

func resize(threadId int, taskChan chan string, wg *sync.WaitGroup) {
	subLogger := log.With().Int("threadId", threadId).Logger()

	var err error

	i := -1
	for task := range taskChan {
		i++
		subLogger.Info().Msgf("Converting %v...", task)

		var outputFile string
		err = func() error {
			f, err := os.Open(task)
			if err != nil {
				return fmt.Errorf("Error while opening file %v: %w", task, err)
			}
			defer f.Close()
			h := md5.New()
			if _, err := io.Copy(h, f); err != nil {
				return fmt.Errorf("Error while calculating md5 for %v: %w", task, err)
			}
			outputFile = filepath.Join(output, hex.EncodeToString(h.Sum(nil))+"_"+filepath.Base(task))
			subLogger.Debug().Msgf("Output for %v: %v", task, outputFile)
			return nil
		}()
		if err != nil {
			subLogger.Error().Msgf("Error while getting output path for %v: %v", task, err)
			continue
		}

		_, err := imagick.ConvertImageCommand([]string{
			"convert", task,"-resize", "640x480", outputFile,
		})
		if err != nil {
			subLogger.Error().Msgf("Error while resizing image %v: %v", task, err)
			continue
		}
	}
	wg.Done()
}

func execute(cmd *cobra.Command, args []string) error {
	var err error

	lvl, err := zerolog.ParseLevel(loggingLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)

	taskChan := make(chan string, threadCount*3)
	var wg sync.WaitGroup
	wg.Add(int(threadCount))

	for i := 0; i < int(threadCount); i++ {
		cur := i
		go resize(cur, taskChan, &wg)
	}

	err = filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			taskChan <- path
		}
		return nil
	})
	close(taskChan)
	if err != nil {
		log.Error().Msgf("Error while browsing source directory: %v", err)
	}

	wg.Wait()
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:          "picture-resizer",
		Short:        "Picture resizer",
		RunE:         execute,
		SilenceUsage: false,
	}
	cmd.Flags().UintVarP(&threadCount, "threadCount", "t", defaultThreadCount, "thread count")
	cmd.Flags().UintVarP(&height, "height", "H", defaultHeight, "height (pixels)")
	cmd.Flags().UintVarP(&width, "width", "w", defaultWidth, "width (pixels)")
	cmd.Flags().UintVarP(&compressionQuality, "compression", "c", defaultCompressionQuality, "compression quality (1-100)")
	cmd.Flags().Float64VarP(&blur, "blur", "b", defaultBlur, "blur (> 1 is blurry, < 1 is sharp)")
	cmd.Flags().StringVarP(&input, "input", "i", "", "input folder")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output folder")
	cmd.Flags().StringVarP(&loggingLevel, "loggingLevel", "l", defaultLoggingLevel, "logging level")
	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")

	if err := cmd.Execute(); err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(retKo)
	}
	os.Exit(retOk)

}
