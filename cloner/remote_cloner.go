package cloner

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"
	specsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type RemoteCloner struct {
	fetcher      Fetcher
	unpacker     Unpacker
	volumeDriver groot.VolumeDriver
}

func NewRemoteCloner(fetcher Fetcher, unpacker Unpacker, volumizer groot.VolumeDriver) *RemoteCloner {
	return &RemoteCloner{
		fetcher:      fetcher,
		unpacker:     unpacker,
		volumeDriver: volumizer,
	}
}

func (c *RemoteCloner) Clone(logger lager.Logger, spec groot.CloneSpec) error {
	logger = logger.Session("remote-cloning", lager.Data{"spec": spec})
	logger.Info("start")
	defer logger.Info("end")

	imageURL, err := url.Parse(spec.Image)
	if err != nil {
		return fmt.Errorf("parsing URL: %s", err)
	}

	digests, config, err := c.fetcher.LayersDigest(logger, imageURL)
	if err != nil {
		return fmt.Errorf("fetching list of digests: %s", err)
	}
	logger.Debug("fetched-layers-digests", lager.Data{"digests": digests})

	if err := c.writeImageJSON(logger, spec.Bundle, config); err != nil {
		return fmt.Errorf("creating image.json: %s", err)
	}

	streamer, err := c.fetcher.Streamer(logger, imageURL)
	if err != nil {
		return fmt.Errorf("initializing streamer: %s", err)
	}

	for _, digest := range digests {
		volumePath, err := c.volumeDriver.Path(logger, wrapVolumeID(spec, digest.ChainID))
		if err == nil {
			logger.Debug("volume-exists", lager.Data{
				"volumePath":    volumePath,
				"blobID":        digest.BlobID,
				"diffID":        digest.DiffID,
				"chainID":       digest.ChainID,
				"parentChainID": digest.ParentChainID,
			})
			continue
		}

		volumePath, err = c.volumeDriver.Create(logger,
			wrapVolumeID(spec, digest.ParentChainID),
			wrapVolumeID(spec, digest.ChainID),
		)
		if err != nil {
			return fmt.Errorf("creating volume for layer `%s`: %s", digest.DiffID, err)
		}
		logger.Debug("volume-created", lager.Data{
			"volumePath":    volumePath,
			"blobID":        digest.BlobID,
			"diffID":        digest.DiffID,
			"chainID":       digest.ChainID,
			"parentChainID": digest.ParentChainID,
		})

		stream, size, err := streamer.Stream(logger, digest.BlobID)
		if err != nil {
			return fmt.Errorf("streaming blob `%s`: %s", digest.BlobID, err)
		}
		logger.Debug("got-stream-for-blob", lager.Data{
			"size":          size,
			"blobID":        digest.BlobID,
			"diffID":        digest.DiffID,
			"chainID":       digest.ChainID,
			"parentChainID": digest.ParentChainID,
		})

		unpackSpec := UnpackSpec{
			TargetPath:  volumePath,
			Stream:      stream,
			UIDMappings: spec.UIDMappings,
			GIDMappings: spec.GIDMappings,
		}
		if err := c.unpacker.Unpack(logger, unpackSpec); err != nil {
			return fmt.Errorf("unpacking layer `%s`: %s", digest.DiffID, err)
		}
		logger.Debug("layer-unpacked", lager.Data{
			"blobID":        digest.BlobID,
			"diffID":        digest.DiffID,
			"chainID":       digest.ChainID,
			"parentChainID": digest.ParentChainID,
		})
	}

	lastVolumeID := wrapVolumeID(spec, digests[len(digests)-1].ChainID)
	if err := c.volumeDriver.Snapshot(logger, lastVolumeID, spec.Bundle.RootFSPath()); err != nil {
		return fmt.Errorf("snapshotting the image to path `%s`: %s", spec.Bundle.RootFSPath(), err)
	}
	logger.Debug("last-volume-got-snapshotted", lager.Data{
		"lastVolumeID": lastVolumeID,
		"rootFSPath":   spec.Bundle.RootFSPath(),
	})

	return nil
}

func wrapVolumeID(spec groot.CloneSpec, volumeID string) string {
	if volumeID == "" {
		return ""
	}

	if len(spec.UIDMappings) > 0 || len(spec.GIDMappings) > 0 {
		return fmt.Sprintf("%s-namespaced", volumeID)
	}

	return volumeID
}

func (c *RemoteCloner) writeImageJSON(logger lager.Logger, bundle groot.Bundle, config specsv1.Image) error {
	logger = logger.Session("writing-image-json")
	logger.Info("start")
	defer logger.Info("end")

	configWriter, err := os.OpenFile(filepath.Join(bundle.Path(), "image.json"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	if err = json.NewEncoder(configWriter).Encode(config); err != nil {
		return err
	}
	return nil
}
