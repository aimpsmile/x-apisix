package testconf

import (
	"log"
	"testing"
)

func TestInitEnv(t *testing.T) {
	serverName := "aimgo.test.v1.testconf"
	_, err := InitEnv(MicroEnv("MICRO_SERVER_NAME", serverName),
		ConfPath([]string{"srv.yml", "db.yml", "broker.yml", "cache.yml"}),
		EtcdPrefix(""),
		AimConfig(GetTestConfPath()),
	)
	if err != nil {
		t.Fatalf("init env error %v", err)
	}
}
func TestMock(t *testing.T) {
	if err := MockConf(); err != nil {
		log.Fatal(err)
	}
}
