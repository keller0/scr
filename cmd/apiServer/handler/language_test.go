package handler

import "testing"

func TestAllRunners(t *testing.T) {
	availableRunners = nil
	t.Log(availableRunners)
	t.Log(expectRunners)
	t.Log(getAllDockerImages())
	allLocalImages := getAllDockerImages()
	for _, r := range expectRunners {
		tmpRunner := &runner{Language: r.Language}

		for _, v := range r.Versions {
			tmpImg := V2Images(r.Language, v.Version)
			if containsString(allLocalImages, tmpImg) {
				tmpRunner.Versions = append(tmpRunner.Versions, v)
			}
		}
		if len(tmpRunner.Versions) > 0 {
			availableRunners = append(availableRunners, *tmpRunner)
		}
	}
	t.Log(availableRunners)

}
