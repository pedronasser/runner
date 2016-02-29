package main

import (
	// "github.com/iron-io/titan/api/client"
	// "github.com/iron-io/titan/api/models"
	log "github.com/Sirupsen/logrus"
	"github.com/iron-io/titan/runner/docker"
	"github.com/iron-io/titan/runner/swagger"
	"os"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	host := os.Getenv("API_URL")
	if host == "" {
		host = "http://localhost:8080"
	}

	jc := swagger.NewCoreApiWithBasePath(host)
	for {
		log.Infoln("Asking for job")
		jobs, err := jc.JobsGet()
		if err != nil {
			log.Errorln("We've got an error!", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(jobs) < 1 || len(jobs[0].Jobs) < 1 {
			time.Sleep(1 * time.Second)
			continue
		}
		job := jobs[0].Jobs[0]
		job.StartedAt = time.Now()
		log.Infoln("Got job:", job)
		s, err := docker.DockerRun(job)
		job.FinishedAt = time.Now()
		if err != nil {
			log.Errorln("We've got an error!", err)
			job.Status = "error"
			job.Error = err.Error()
			jc.JobIdPatch(job.Id, swagger.JobWrapper{job})
			continue
		}
		job.Status = "success"
		jc.JobIdPatch(job.Id, swagger.JobWrapper{job})
		log.Infoln("output:", s)
	}
}
