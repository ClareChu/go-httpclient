package go_httpclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPost(t *testing.T) {
	s := New()
	baseUrl := "http://dev.apps.cloud2go.cn/cloudos/dev/manager"
	_, body, err := s.Get(fmt.Sprintf("%s/%s", baseUrl, "pipeline/id/8a963a19-a2b6-11e9-b657-0242ac121217"), nil)
	assert.Equal(t, nil, err)
	b := string(body)
	fmt.Print(b)
}
