/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package p9

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControllerServer struct {
	Driver *Driver
	// Working directory for the provisioner to temporarily mount p9 shares at
	workingMountDir string
}

// p9Volume is an internal representation of a volume
// created by the provisioner.
type p9Volume struct {
	// Volume id
	id string
	// Address of the P9 server.
	// Matches paramServer.
	server string
	// Base directory of the P9 server to create volumes under
	// Matches paramShare.
	baseDir string
	// Subdirectory of the P9 server to create volumes under
	subDir string
	// size of volume
	size int64
}

// Ordering of elements in the CSI volume id.
// ID is of the form {server}/{baseDir}/{subDir}.
// TODO: This volume id format limits baseDir and
// subDir to only be one directory deep.
// Adding a new element should always go at the end
// before totalIDElements
const (
	idServer = iota
	idBaseDir
	idSubDir
	totalIDElements // Always last
)

func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// Validate arguments
	name := req.GetName()
	if len(name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume name must be provided")
	}

	if err := cs.validateVolumeCapabilities(req.GetVolumeCapabilities()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	reqCapacity := req.GetCapacityRange().GetRequiredBytes()
	p9Vol, err := cs.newP9Volume(name, reqCapacity, req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Mount p9 base share so we can create a subdirectory
	if err = cs.internalMount(ctx, p9Vol); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mount p9 server: %v", err.Error())
	}
	defer func() {
		if err = cs.internalUnmount(ctx, p9Vol); err != nil {
			glog.Warningf("failed to unmount p9 server: %v", err.Error())
		}
	}()

	// Create subdirectory under base-dir
	// TODO: revisit permissions
	internalVolumePath := cs.getInternalVolumePath(p9Vol)
	if err = os.Mkdir(internalVolumePath, 0777); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make subdirectory: %v", err.Error())
	}
	// Remove capacity setting when provisioner 1.4.0 is available with fix for
	// https://github.com/kubernetes-csi/external-provisioner/pull/271
	return &csi.CreateVolumeResponse{Volume: cs.p9VolToCSI(p9Vol, reqCapacity)}, nil
}

func (cs *ControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	if volumeID == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	p9Vol, err := cs.getNfsVolFromID(volumeID)
	if err != nil {
		// An invalid ID should be treated as doesn't exist
		glog.V(5).Infof("failed to get p9 volume for volume id %v deletion: %v", volumeID, err)
		return &csi.DeleteVolumeResponse{}, nil
	}

	// Mount p9 base share so we can delete the subdirectory
	if err = cs.internalMount(ctx, p9Vol); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mount p9 server: %v", err.Error())
	}
	defer func() {
		if err = cs.internalUnmount(ctx, p9Vol); err != nil {
			glog.Warningf("failed to unmount p9 server: %v", err.Error())
		}
	}()

	// Delete subdirectory under base-dir
	internalVolumePath := cs.getInternalVolumePath(p9Vol)

	glog.V(4).Infof("Removing subdirectory at %v", internalVolumePath)
	if err = os.RemoveAll(internalVolumePath); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete subdirectory: %v", err.Error())
	}

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *ControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if req.GetVolumeCapabilities() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	}

	// supports all AccessModes, no need to check capabilities here
	return &csi.ValidateVolumeCapabilitiesResponse{Message: ""}, nil
}

func (cs *ControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetCapabilities implements the default GRPC callout.
// Default supports all capabilities
func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	glog.V(5).Infof("Using default ControllerGetCapabilities")

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.Driver.cscap,
	}, nil
}

func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) validateVolumeCapabilities(caps []*csi.VolumeCapability) error {
	if len(caps) == 0 {
		return fmt.Errorf("volume capabilities must be provided")
	}

	for _, c := range caps {
		if err := cs.validateVolumeCapability(c); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ControllerServer) validateVolumeCapability(c *csi.VolumeCapability) error {
	if c == nil {
		return fmt.Errorf("volume capability must be provided")
	}

	// Validate access mode
	accessMode := c.GetAccessMode()
	if accessMode == nil {
		return fmt.Errorf("volume capability access mode not set")
	}
	if !cs.Driver.cap[accessMode.Mode] {
		return fmt.Errorf("driver does not support access mode: %v", accessMode.Mode.String())
	}

	// Validate access type
	accessType := c.GetAccessType()
	if accessType == nil {
		return fmt.Errorf("volume capability access type not set")
	}
	return nil
}

// Mount p9 server at base-dir
func (cs *ControllerServer) internalMount(ctx context.Context, vol *p9Volume) error {
	sharePath := filepath.Join(string(filepath.Separator) + vol.baseDir)
	targetPath := cs.getInternalMountPath(vol)
	stdVolCap := csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{},
		},
	}

	glog.V(4).Infof("internally mounting %v:%v at %v", vol.server, sharePath, targetPath)
	_, err := cs.Driver.ns.NodePublishVolume(ctx, &csi.NodePublishVolumeRequest{
		TargetPath: targetPath,
		VolumeContext: map[string]string{
			paramServer: vol.server,
			paramShare:  sharePath,
		},
		VolumeCapability: &stdVolCap,
		VolumeId:         vol.id,
	})
	return err
}

// Unmount p9 server at base-dir
func (cs *ControllerServer) internalUnmount(ctx context.Context, vol *p9Volume) error {
	targetPath := cs.getInternalMountPath(vol)

	// Unmount p9 server at base-dir
	glog.V(4).Infof("internally unmounting %v", targetPath)
	_, err := cs.Driver.ns.NodeUnpublishVolume(ctx, &csi.NodeUnpublishVolumeRequest{
		VolumeId:   vol.id,
		TargetPath: cs.getInternalMountPath(vol),
	})
	return err
}

// Convert VolumeCreate parameters to an p9Volume
func (cs *ControllerServer) newP9Volume(name string, size int64, params map[string]string) (*p9Volume, error) {
	var (
		server  string
		baseDir string
	)

	// Validate parameters (case-insensitive).
	// TODO do more strict validation.
	for k, v := range params {
		switch strings.ToLower(k) {
		case paramServer:
			server = v
		case paramShare:
			baseDir = v
		default:
			return nil, fmt.Errorf("invalid parameter %q", k)
		}
	}

	// Validate required parameters
	if server == "" {
		return nil, fmt.Errorf("%v is a required parameter", paramServer)
	}
	if baseDir == "" {
		return nil, fmt.Errorf("%v is a required parameter", paramShare)
	}

	vol := &p9Volume{
		server:  server,
		baseDir: baseDir,
		subDir:  name,
		size:    size,
	}
	vol.id = cs.getVolumeIDFromNfsVol(vol)

	return vol, nil
}

// Get working directory for CreateVolume and DeleteVolume
func (cs *ControllerServer) getInternalMountPath(vol *p9Volume) string {
	// use default if empty
	if cs.workingMountDir == "" {
		cs.workingMountDir = "/tmp"
	}
	return filepath.Join(cs.workingMountDir, vol.subDir)
}

// Get internal path where the volume is created
// The reason why the internal path is "workingDir/subDir/subDir" is because:
//   * the semantic is actually "workingDir/volId/subDir" and volId == subDir.
//   * we need a mount directory per volId because you can have multiple
//     CreateVolume calls in parallel and they may use the same underlying share.
//     Instead of refcounting how many CreateVolume calls are using the same
//     share, it's simpler to just do a mount per request.
func (cs *ControllerServer) getInternalVolumePath(vol *p9Volume) string {
	return filepath.Join(cs.getInternalMountPath(vol), vol.subDir)
}

// Get user-visible share path for the volume
func (cs *ControllerServer) getVolumeSharePath(vol *p9Volume) string {
	return filepath.Join(string(filepath.Separator), vol.baseDir, vol.subDir)
}

// Convert into p9Volume into a csi.Volume
func (cs *ControllerServer) p9VolToCSI(vol *p9Volume, reqCapacity int64) *csi.Volume {
	return &csi.Volume{
		CapacityBytes: reqCapacity,
		VolumeId:      vol.id,
		VolumeContext: map[string]string{
			paramServer: vol.server,
			paramShare:  cs.getVolumeSharePath(vol),
		},
	}
}

// Given a p9Volume, return a CSI volume id
func (cs *ControllerServer) getVolumeIDFromNfsVol(vol *p9Volume) string {
	idElements := make([]string, totalIDElements)
	idElements[idServer] = strings.Trim(vol.server, "/")
	idElements[idBaseDir] = strings.Trim(vol.baseDir, "/")
	idElements[idSubDir] = strings.Trim(vol.subDir, "/")
	return strings.Join(idElements, "/")
}

// Given a CSI volume id, return a p9Volume
func (cs *ControllerServer) getNfsVolFromID(id string) (*p9Volume, error) {
	tokens := strings.Split(id, "/")
	if len(tokens) != totalIDElements {
		return nil, fmt.Errorf("volume id %q unexpected format: got %v token(s) instead of %v", id, len(tokens), totalIDElements)
	}

	return &p9Volume{
		id:      id,
		server:  tokens[idServer],
		baseDir: tokens[idBaseDir],
		subDir:  tokens[idSubDir],
	}, nil
}
