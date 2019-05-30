package providers

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift-metalkube/kni-installer/pkg/types"
)

// Destroyer allows multiple implementations of destroy
// for different platforms.
type Destroyer interface {
	Run() error
}

// NewFunc is an interface for creating platform-specific destroyers.
type NewFunc func(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (Destroyer, error)