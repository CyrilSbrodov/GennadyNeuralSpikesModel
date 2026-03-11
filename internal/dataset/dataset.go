package dataset

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Sample struct {
	Label     string
	ImagePath string
}

func Load(root string) ([]Sample, error) {
	entries, readDirErr := os.ReadDir(root)
	if readDirErr != nil {
		return nil, readDirErr
	}
	samples := make([]Sample, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		label := entry.Name()
		dir := filepath.Join(root, label)
		imgs, readImgDirErr := os.ReadDir(dir)
		if readImgDirErr != nil {
			return nil, readImgDirErr
		}
		for _, img := range imgs {
			if img.IsDir() {
				continue
			}
			fmt.Println(img)
			ext := strings.ToLower(filepath.Ext(img.Name()))
			if ext != ".jpg" && ext != ".png" && ext != ".jpeg" && ext != ".gif" {
				continue
			}
			samples = append(samples, Sample{Label: img.Name(), ImagePath: filepath.Join(dir, img.Name())})
		}
	}
	sort.Slice(samples, func(i, j int) bool {
		if samples[i].Label == samples[j].Label {
			return samples[i].ImagePath < samples[j].ImagePath
		}
		return samples[i].Label < samples[j].Label
	})
	if len(samples) == 0 {
		return nil, errors.New("no samples found")
	}

	return samples, nil
}
