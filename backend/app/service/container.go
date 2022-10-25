package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/dto"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/1Panel-dev/1Panel/backend/utils/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerService struct{}

type IContainerService interface {
	Page(req dto.PageContainer) (int64, interface{}, error)
	PageNetwork(req dto.PageInfo) (int64, interface{}, error)
	PageVolume(req dto.PageInfo) (int64, interface{}, error)
	ListVolume() ([]dto.Options, error)
	PageCompose(req dto.PageInfo) (int64, interface{}, error)
	CreateCompose(req dto.ComposeCreate) error
	ComposeOperation(req dto.ComposeOperation) error
	ContainerCreate(req dto.ContainerCreate) error
	ContainerOperation(req dto.ContainerOperation) error
	ContainerLogs(param dto.ContainerLog) (string, error)
	ContainerStats(id string) (*dto.ContainterStats, error)
	Inspect(req dto.InspectReq) (string, error)
	DeleteNetwork(req dto.BatchDelete) error
	CreateNetwork(req dto.NetworkCreat) error
	DeleteVolume(req dto.BatchDelete) error
	CreateVolume(req dto.VolumeCreat) error
}

func NewIContainerService() IContainerService {
	return &ContainerService{}
}

func (u *ContainerService) Page(req dto.PageContainer) (int64, interface{}, error) {
	var (
		records   []types.Container
		list      []types.Container
		backDatas []dto.ContainerInfo
	)
	client, err := docker.NewDockerClient()
	if err != nil {
		return 0, nil, err
	}
	options := types.ContainerListOptions{All: true}
	if len(req.Filters) != 0 {
		options.Filters = filters.NewArgs()
		options.Filters.Add("label", req.Filters)
	}
	list, err = client.ContainerList(context.Background(), options)
	if err != nil {
		return 0, nil, err
	}
	total, start, end := len(list), (req.Page-1)*req.PageSize, req.Page*req.PageSize
	if start > total {
		records = make([]types.Container, 0)
	} else {
		if end >= total {
			end = total
		}
		records = list[start:end]
	}

	for _, container := range records {
		backDatas = append(backDatas, dto.ContainerInfo{
			ContainerID: container.ID,
			CreateTime:  time.Unix(container.Created, 0).Format("2006-01-02 15:04:05"),
			Name:        container.Names[0][1:],
			ImageId:     strings.Split(container.ImageID, ":")[1],
			ImageName:   container.Image,
			State:       container.State,
			RunTime:     container.Status,
		})
	}

	return int64(total), backDatas, nil
}

func (u *ContainerService) Inspect(req dto.InspectReq) (string, error) {
	client, err := docker.NewDockerClient()
	if err != nil {
		return "", err
	}
	var inspectInfo interface{}
	switch req.Type {
	case "container":
		inspectInfo, err = client.ContainerInspect(context.Background(), req.ID)
	case "network":
		inspectInfo, err = client.NetworkInspect(context.TODO(), req.ID, types.NetworkInspectOptions{})
	case "volume":
		inspectInfo, err = client.VolumeInspect(context.TODO(), req.ID)
	}
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(inspectInfo)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *ContainerService) ContainerCreate(req dto.ContainerCreate) error {
	client, err := docker.NewDockerClient()
	if err != nil {
		return err
	}
	config := &container.Config{
		Image:  req.Image,
		Cmd:    req.Cmd,
		Env:    req.Env,
		Labels: stringsToMap(req.Labels),
	}
	hostConf := &container.HostConfig{
		AutoRemove:      req.AutoRemove,
		PublishAllPorts: req.PublishAllPorts,
		RestartPolicy:   container.RestartPolicy{Name: req.RestartPolicy},
	}
	if req.RestartPolicy == "on-failure" {
		hostConf.RestartPolicy.MaximumRetryCount = 5
	}
	if req.NanoCPUs != 0 {
		hostConf.NanoCPUs = req.NanoCPUs * 1000000000
	}
	if req.Memory != 0 {
		hostConf.Memory = req.Memory
	}
	if len(req.ExposedPorts) != 0 {
		hostConf.PortBindings = make(nat.PortMap)
		for _, port := range req.ExposedPorts {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(port.HostPort)}
			hostConf.PortBindings[nat.Port(fmt.Sprintf("%d/tcp", port.ContainerPort))] = []nat.PortBinding{bindItem}
		}
	}
	if len(req.Volumes) != 0 {
		config.Volumes = make(map[string]struct{})
		for _, volume := range req.Volumes {
			config.Volumes[volume.ContainerDir] = struct{}{}
			hostConf.Binds = append(hostConf.Binds, fmt.Sprintf("%s:%s:%s", volume.SourceDir, volume.ContainerDir, volume.Mode))
		}
	}
	container, err := client.ContainerCreate(context.TODO(), config, hostConf, &network.NetworkingConfig{}, &v1.Platform{}, req.Name)
	if err != nil {
		return err
	}
	if err := client.ContainerStart(context.TODO(), container.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("create successful but start failed, err: %v", err)
	}
	return nil
}

func (u *ContainerService) ContainerOperation(req dto.ContainerOperation) error {
	var err error
	ctx := context.Background()
	client, err := docker.NewDockerClient()
	if err != nil {
		return err
	}
	switch req.Operation {
	case constant.ContainerOpStart:
		err = client.ContainerStart(ctx, req.ContainerID, types.ContainerStartOptions{})
	case constant.ContainerOpStop:
		err = client.ContainerStop(ctx, req.ContainerID, nil)
	case constant.ContainerOpRestart:
		err = client.ContainerRestart(ctx, req.ContainerID, nil)
	case constant.ContainerOpKill:
		err = client.ContainerKill(ctx, req.ContainerID, "SIGKILL")
	case constant.ContainerOpPause:
		err = client.ContainerPause(ctx, req.ContainerID)
	case constant.ContainerOpUnpause:
		err = client.ContainerUnpause(ctx, req.ContainerID)
	case constant.ContainerOpRename:
		err = client.ContainerRename(ctx, req.ContainerID, req.NewName)
	case constant.ContainerOpRemove:
		err = client.ContainerRemove(ctx, req.ContainerID, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true})
	}
	return err
}

func (u *ContainerService) ContainerLogs(req dto.ContainerLog) (string, error) {
	var (
		options types.ContainerLogsOptions
		logs    io.ReadCloser
		buf     *bytes.Buffer
		err     error
	)
	client, err := docker.NewDockerClient()
	if err != nil {
		return "", err
	}
	options = types.ContainerLogsOptions{
		ShowStdout: true,
		Timestamps: true,
	}
	if req.Mode != "all" {
		options.Since = req.Mode
	}
	if logs, err = client.ContainerLogs(context.Background(), req.ContainerID, options); err != nil {
		return "", err
	}
	defer logs.Close()
	buf = new(bytes.Buffer)
	if _, err = stdcopy.StdCopy(buf, nil, logs); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (u *ContainerService) ContainerStats(id string) (*dto.ContainterStats, error) {
	client, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	res, err := client.ContainerStats(context.TODO(), id, false)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var stats *types.StatsJSON
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}
	var data dto.ContainterStats
	previousCPU := stats.PreCPUStats.CPUUsage.TotalUsage
	previousSystem := stats.PreCPUStats.SystemUsage
	data.CPUPercent = calculateCPUPercentUnix(previousCPU, previousSystem, stats)
	data.IORead, data.IOWrite = calculateBlockIO(stats.BlkioStats)
	data.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		data.Cache = float64(cache) / 1024 / 1024
	}
	data.Memory = data.Memory - data.Cache
	data.NetworkRX, data.NetworkTX = calculateNetwork(stats.Networks)
	data.ShotTime = stats.Read
	return &data, nil
}

func stringsToMap(list []string) map[string]string {
	var lableMap = make(map[string]string)
	for _, label := range list {
		sps := strings.Split(label, "=")
		if len(sps) > 1 {
			lableMap[sps[0]] = sps[1]
		}
	}
	return lableMap
}
func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}
func calculateBlockIO(blkio types.BlkioStats) (blkRead float64, blkWrite float64) {
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = (blkRead + float64(bioEntry.Value)) / 1024 / 1024
		case "write":
			blkWrite = (blkWrite + float64(bioEntry.Value)) / 1024 / 1024
		}
	}
	return
}
func calculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes) / 1024
		tx += float64(v.TxBytes) / 1024
	}
	return rx, tx
}