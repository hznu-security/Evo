/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 17:44
 * 描述     ：管理镜像的相关函数
 */

package docker

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
)

var (
	ErrRead = errors.New("读取buildResponse错误")
)

// BuildImage 构建镜像
func BuildImage(path string, dockerfile string, name string) ([]byte, error) {
	ctx := context.Background()
	cli, err := getClient()
	if err != nil {
		return nil, err
	}
	buildContext, err := os.Open(path) //这里必须接受一个打包的文件
	if err != nil {
		return nil, err
	}
	buildResponse, err := cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Tags:        []string{name},
		Remove:      true,
		ForceRemove: true,
		Dockerfile:  dockerfile,
	})
	if err != nil {
		return nil, fmt.Errorf("镜像构建错误 %v", err)
	}
	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("读取buildResponse错误 %v", err)
	}

	err = buildResponse.Body.Close()
	err = buildContext.Close()
	return response, nil
}

// Image 包括id和tags
type Image struct {
	ID   string `json:"Id"`
	Tags []string
}

func ListImage() ([]Image, error) {
	ctx := context.Background()
	cli, err := getClient()
	if err != nil {
		return nil, err
	}
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	res := make([]Image, len(images))
	for k, img := range images {
		image := Image{
			ID:   img.ID,
			Tags: img.RepoTags,
		}
		res[k] = image
	}
	return res, nil
}

// RemoveImage 删除镜像
func RemoveImage(id string) error {
	ctx := context.Background()
	cli, err := getClient()
	if err != nil {
		return err
	}
	_, err = cli.ImageRemove(ctx, id, types.ImageRemoveOptions{})
	if err != nil {
		return fmt.Errorf("删除镜像失败 %v", err)
	}

	return nil
}

func CheckImage(imgId string) error {
	ctx := context.Background()
	cli, err := getClient()
	if err != nil {
		return err
	}
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, img := range images {
		if img.ID == imgId {
			return nil
		}
	}
	return errors.New("镜像不存在")
}
