/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/12 11:19
 * 描述     ：管理镜像相关
 */

package manage

import (
	"Evo/service/docker"
	"Evo/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// PostImage 接受上传的tar文件，打包成镜像
// 注意！必须是一个文件夹（名字为题目名）打包为名字.tar的打包文件   //TODO
func PostImage(c *gin.Context) {
	// 解析表单
	file, err := c.FormFile("image")
	if err != nil {
		util.Fail(c, "上传失败", nil)
		return
	}
	name := c.PostForm("name")

	if name == "" {
		util.Fail(c, "缺少参数", nil)
		return
	}

	// 检查字段是否过长
	if len(name) > 200 {
		util.Fail(c, "字段过长", gin.H{
			"name": name,
		})
	}

	imagePath := viper.GetString("image.path")
	dst := imagePath + name + ".tar"
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		util.Error(c, "保存失败", nil)
		log.Println(err.Error())
		return
	}
	resp, err := docker.BuildImage(dst, "Dockerfile", name)
	if err != nil {
		log.Println(err.Error())
		if err == docker.ErrRead {
			util.Success(c, "读取异常", nil)
			return
		} else {
			log.Println(err)
			util.Fail(c, "镜像构建失败", gin.H{
				"error": err,
				"resp":  string(resp),
			})
			return
		}
	}
	util.Success(c, "", gin.H{
		"process": string(resp),
	})
}

// GetImage 列出所有镜像
func GetImage(c *gin.Context) {
	images, err := docker.ListImage()
	if err != nil {
		log.Println(err.Error())
		util.Error(c, "出错了", nil)
		return
	}
	util.Success(c, "成功", gin.H{
		"images": images,
	})
}

func DelImage(c *gin.Context) {
	imageId := c.Query("image")
	err := docker.RemoveImage(imageId)
	if err != nil {
		log.Println(err.Error())
		util.Error(c, "删除失败", gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("删除镜像", imageId)
	util.Success(c, "成功", nil)
}
