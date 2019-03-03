package leifengtrend

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// Config OSS 配置
type Config struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// Object       string
type Object struct {
	ObjectName    string
	LocalFileName string
}

// PutOSS 数据上传至OSS
func PutOSS(cf *Config, obj *Object) error {
	// Endpoint以杭州为例，其它Region请按实际情况填写。 endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// accessKeyID 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	// 创建OSSClient实例。
	client, err := oss.New(cf.Endpoint, cf.AccessKeyID, cf.AccessKeySecret)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(cf.BucketName)
	if err != nil {
		return err
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(obj.ObjectName, obj.LocalFileName)
	if err != nil {
		return err
	}
	return nil
}
