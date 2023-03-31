package docker

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func TestCreateBash(t *testing.T) {
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
