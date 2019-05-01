package internal

import (
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
)

const (
	dockerPSQL = `andboson/postgres-multiple-db`
)

type DockerSuite struct {
	suite.Suite

	pool      *dockertest.Pool
	resources []*dockertest.Resource
	authConf  docker.AuthConfiguration
	networkID string
	rootID    string
}

func (s *DockerSuite) SetupPSQL(databases string) (addr string) {
	resource, err := s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: dockerPSQL,
		Env: []string{
			"POSTGRES_PASSWORD=admin",
			"POSTGRES_USER=admin",
			"POSTGRES_MULTIPLE_DATABASES=" + databases,
		},
	})

	s.Require().NoError(err)

	s.resources = append(s.resources, resource)

	ip, err := s.ConnectNetwork(s.networkID, resource.Container.ID, "db")
	s.Require().NoError(err)

	addr = ip.String() + ":5432"
	if runtime.GOOS == "darwin" { // hacks for mac local docker network
		addr = "localhost:" + resource.GetPort("5432/tcp")
	}

	log.Printf("docker started: %s - %s", "Postgresql", resource.Container.ID[:13])

	return addr
}

func (s *DockerSuite) ConnectNetwork(networkID, containerID, alias string) (net.IP, error) {
	err := s.pool.Client.ConnectNetwork(networkID, docker.NetworkConnectionOptions{
		Container: containerID,
		EndpointConfig: &docker.EndpointConfig{
			Aliases: []string{alias},
		},
	})
	if err != nil {
		return nil, err
	}

	nk, err := s.pool.Client.NetworkInfo(networkID)
	if err != nil {
		return nil, err
	}

	ip, _, err := net.ParseCIDR(nk.Containers[containerID].IPv4Address)
	if err != nil {
		return nil, err
	}

	return ip.To4(), nil
}

func (s *DockerSuite) Down() {
	for _, resource := range s.resources {
		err := s.pool.Purge(resource)
		s.Require().NoError(err, `Could not purge resource: %s`, err)
	}

	if s.rootID != "" {
		err := s.pool.Client.DisconnectNetwork(s.networkID, docker.NetworkConnectionOptions{
			Container: s.rootID,
		})
		s.NoError(err, `Could not disconnect network: %s`, err)
	}

	err := s.pool.Client.RemoveNetwork(s.networkID)
	s.Require().NoError(err, `Could not remove network: %s`, err)
}

func (s *DockerSuite) Setup(networkName string) {
	var err error

	s.pool, err = dockertest.NewPool(``)
	s.Require().NoError(err)

	nk, err := s.pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: networkName})
	s.Require().NoError(err)
	s.networkID = nk.ID

	c, err := s.pool.Client.ListContainers(docker.ListContainersOptions{})
	s.Require().NoError(err)

	hostName, _ := os.Hostname()
	for _, container := range c {
		if strings.Contains(container.Names[0], hostName) {
			s.rootID = container.ID
			_, err = s.ConnectNetwork(s.networkID, s.rootID, hostName)
			s.Require().NoError(err)
			break
		}
	}
}
