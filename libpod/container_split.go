//go:build !remote && (linux || freebsd)

package libpod

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.podman.io/podman/v6/libpod/define"
)

// SplitContainer creates a COW child from a running parent container
// via the OCI runtime's split command. The child shares the parent's
// memory (copy-on-write) and rootfs (overlayfs).
//
// The child bundle is created under the parent's GraphRoot so storage
// cleanup follows Podman's normal lifecycle.
func (r *Runtime) SplitContainer(parent *Container, childName string, noCleanup bool, shareNetwork bool, shareIPC bool, shareUTS bool, sharePID bool) error {
	if parent == nil {
		return fmt.Errorf("parent container must not be nil: %w", define.ErrInvalidArg)
	}

	if !parent.ensureState(define.ContainerStateRunning, define.ContainerStatePaused) {
		return fmt.Errorf("parent container %s must be running or paused to split: %w", parent.ID(), define.ErrCtrStateInvalid)
	}

	if !parent.ociRuntime.SupportsSplit() {
		return fmt.Errorf("configured OCI runtime %s does not support container split: %w", parent.ociRuntime.Name(), define.ErrInvalidArg)
	}

	// Child bundle lives under the runtime graph root so it is cleaned
	// up alongside other container data by Podman.
	childBundle := filepath.Join(r.GraphRoot(), "split", childName)

	logrus.Debugf("Live-clone/split container %s -> %s (bundle=%s, noCleanup=%t)",
		parent.ID(), childName, childBundle, noCleanup)

	return parent.ociRuntime.SplitContainer(parent, childName, childBundle, noCleanup, shareNetwork, shareIPC, shareUTS, sharePID)
}
