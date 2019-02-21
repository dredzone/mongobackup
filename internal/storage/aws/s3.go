package aws

type S3 struct {
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	API       string `yaml:"api"`
	SecretKey string `yaml:"secretKey"`
	URL       string `yaml:"url"`
}

