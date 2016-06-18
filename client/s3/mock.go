package s3

type MockClient struct {
}

func (c MockClient) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {

	return nil, nil
}
