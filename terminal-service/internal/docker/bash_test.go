package docker

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func TestCreateBash(t *testing.T) {
	// github action 不执行该测试
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	bashContainer, err := NewBashContainer(fmt.Sprint(time.Now().Unix()))
	if err != nil {
		log.Fatal().Err(err)
		t.Fail()
	}

	if err = bashContainer.Remove(); err != nil {
		log.Fatal().Err(err)
		t.Fail()
	}

	log.Info().Msgf("container ID: %s, container name: %s", bashContainer.id, bashContainer.name)
}
