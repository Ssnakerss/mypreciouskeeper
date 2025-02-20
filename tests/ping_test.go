package tests

import (
	"testing"

	"github.com/Ssnakerss/mypreciouskeeper/tests/suite"
)

func Test_Ping(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	resp, err := st.PingClient.Ping(ctx, nil)
	t.Log(resp)
	t.Log(err)
}
