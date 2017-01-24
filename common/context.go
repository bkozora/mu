package common

import (
	"bufio"
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path"
	"time"
)

var version string

// GetVersion returns the current version of the app
func GetVersion() string {
	if version == "" {
		return "0.0.0-local"
	}
	return version
}

// SetVersion returns the current version of the app
func SetVersion(v string) {
	version = v
}

// NewContext create a new context object
func NewContext() *Context {
	ctx := new(Context)
	return ctx
}

// InitializeFromFile loads config from file
func (ctx *Context) InitializeFromFile(muFile string) error {
	// load yaml config
	yamlFile, err := os.Open(muFile)
	if err != nil {
		return err
	}
	defer func() {
		yamlFile.Close()
	}()

	// set the basedir
	ctx.Config.Basedir = path.Dir(muFile)
	ctx.Repo.Name = path.Base(ctx.Config.Basedir)
	ctx.Repo.Revision = time.Now().Format("20060102150405")
	gitRevision, err := findGitRevision(muFile)
	if err == nil {
		ctx.Repo.Revision = gitRevision
	}

	return ctx.Initialize(bufio.NewReader(yamlFile))
}

// Initialize loads config object
func (ctx *Context) Initialize(configReader io.Reader) error {

	// load the configuration
	err := loadYamlConfig(&ctx.Config, configReader)
	if err != nil {
		return err
	}

	// initialize StackManager
	ctx.StackManager, err = newStackManager(ctx.Config.Region)
	if err != nil {
		return err
	}

	// initialize ClusterManager
	ctx.ClusterManager, err = newClusterManager(ctx.Config.Region)
	if err != nil {
		return err
	}

	// initialize DockerManager
	ctx.DockerManager, err = newClientDockerManager()
	if err != nil {
		return err
	}

	return nil
}

func loadYamlConfig(config *Config, yamlReader io.Reader) error {
	yamlBuffer := new(bytes.Buffer)
	yamlBuffer.ReadFrom(yamlReader)
	return yaml.Unmarshal(yamlBuffer.Bytes(), config)
}