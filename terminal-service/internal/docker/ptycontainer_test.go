package docker

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateBash(t *testing.T) {
	// github action 不执行该测试
	if testing.Short() {
		t.Skip("Skipping test in github action.")
	}

	bashContainer, err := NewPtyContainer("bash:alpine3.16", fmt.Sprint(time.Now().Unix()))
	if err != nil {
		t.Error(err)
	}

	if err = bashContainer.Start(); err != nil {
		t.Error(err)
	}

	if err = bashContainer.Remove(); err != nil {
		t.Error(err)
	}
}
