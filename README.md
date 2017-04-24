# face_utils #


Face项目[face](https://github.com/xpzouying/face)相关的一些工具。

比如，

[X] 对一个目录批量抽取人脸并且保存



# 接口说明 #
## 批量detect人脸 ##
批量detect人脸。将`src`目录下的图片文件批量检测人脸，并且将人脸框保存到`des`目录下。

```bash
go run *.go -src=face_images/ -des=/tmp/faces/
```

## 类型说明 ##

### FaceInfo ###
人脸信息。  
- Image：人脸框图片，jpeg, base64格式。
- Rect：人脸框坐标信息，[left, top, right, bottom]

```go
type FaceInfo struct {
	Image []byte `json:"image"`
	Rect  [4]int `json:"rect"`
}
```

### FaceDetectResult ###
人脸检测响应信息。  
- TimeUsed: 检测人脸的耗时。
- FacesInfo：人脸信息。如图片存在多人脸，则返回多个人脸信息。

```go
type FaceDetectResult struct {
	TimeUsed  float64    `json:"time_used"`
	FacesInfo []FaceInfo `json:"faces"`
}


```