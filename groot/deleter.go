package groot

import "code.cloudfoundry.org/lager"

type Deleter struct {
	bundler           Bundler
	dependencyManager DependencyManager
}

func IamDeleter(bundler Bundler, dependencyManager DependencyManager) *Deleter {
	return &Deleter{
		bundler:           bundler,
		dependencyManager: dependencyManager,
	}
}

func (d *Deleter) Delete(logger lager.Logger, id string) error {
	logger = logger.Session("groot-deleting", lager.Data{"bundleID": id})
	logger.Info("start")
	defer logger.Info("end")

	err := d.bundler.Destroy(logger, id)
	if derErr := d.dependencyManager.Deregister(id); derErr != nil {
		logger.Error("failed-to-deregister-dependencies", derErr)
	}

	return err
}
