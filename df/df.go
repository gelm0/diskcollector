package diskcollector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type UnixDiskStat struct {
	PrometheusDiskInfo
}

func (p UnixDiskStat) Collect() PrometheusDiskInfo {
	statfs, err := StatDisk(p.Path)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  p.Path,
			"error": err,
		}).Fatal("Not able to stat path")
	}
	p.Bsize = float64(statfs.Blocks * uint64(statfs.Bsize))
	p.Bused = float64(statfs.Blocks*uint64(statfs.Bsize) - statfs.Bfree*uint64(statfs.Bsize))
	p.Bfree = float64(statfs.Bavail * uint64(statfs.Bsize))
	p.Bavailable = float64(statfs.Bfree * uint64(statfs.Bsize))
	return p.PrometheusDiskInfo
}

// Init the description for the Prometheus metrics for disks
func InitPd(path string) (p *UnixDiskStat) {
	p = &UnixDiskStat{}
	p.InitDp(path)
	dl := prometheus.Labels{"mount": p.Path}
	p.DescSize = prometheus.NewDesc("disk_size_bytes",
		"Total disk space in bytes",
		nil, dl)
	p.DescUsed = prometheus.NewDesc("disk_used_bytes",
		"Total usage of the disk in bytes",
		nil,
		dl)
	p.DescFree = prometheus.NewDesc("disk_free_bytes",
		"Total free space disk in bytes",
		nil,
		dl)
	p.DescAvailable = prometheus.NewDesc("disk_available_bytes",
		"Total available disk space left in bytes",
		nil,
		dl)
	return
}

func (d *UnixDiskStat) InitDp(path string) {
	if path == "" {
		log.Info("No path given, defaulting to root /")
		path = "/"
	}
	log.Debug("Collecing metrics from path: ", path)
	d.Path = path
}

func StatDisk(path string) (*unix.Statfs_t, error) {
	statfs := unix.Statfs_t{}
	err := unix.Statfs(path, &statfs)
	if err != nil {
		return nil, fmt.Errorf("unable to call statfs on %s, %v", path, err)
	}
	return &statfs, nil
}
