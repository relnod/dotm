package e2e

import (
	"archive/tar"
	"bytes"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func BuildImage() error {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return err
	}

	t := time.Now()

	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte("From base\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
		Name:         "test",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := client.BuildImage(opts); err != nil {
		return err
	}

	return nil
}

func Test1(t *testing.T) {
	if err := BuildImage(); err != nil {
		t.Fatal(err)
	}
}

/*
import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.35"))
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}

func (c *Client) buildImage(ctx context.Context) error {

	dockerBuildContext, err := os.Create("/tmp/tarfile.tar")
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "test/e2e/Dockerfile",
	}

	resp, err := c.cli.ImageBuild(ctx, dockerBuildContext, buildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) startContainer(ctx context.Context) error {
	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: "dotm:latest",
		User:  "testuser",
	}, nil, nil, "")
	if err != nil {
		return err
	}

	err = c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := c.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)

	return nil
}

func TestSmoke(t *testing.T) {
	ctx := context.Background()

	c, err := NewClient()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.buildImage(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.startContainer(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	// subCmds := []runner.DotmCmd{
	// 	runner.DotmCmd{
	// 		SubCommand: "install",
	// 		Params:     make(map[string]string),
	// 	},
	// 	runner.DotmCmd{
	// 		SubCommand: "uninstall",
	// 		Params:     make(map[string]string),
	// 	},
	// }

	// for _, subCmd := range subCmds {
	// 	t.Run(subCmd.SubCommand, func(tt *testing.T) {
	// 		fs := testutil.NewFileSystem()
	// 		defer fs.Cleanup()
	// 		fs.MkdirAll("dotiles")
	// 		fs.MkdirAll("home")
	// 		fs.Create("dotiles/a.txt")

	// 		subCmd.Params["source"] = fs.BasePath()
	// 		subCmd.Params["destination"] = fs.Path("home")

	// 		r := runner.Run(subCmd)
	// 		assert.ErrorIsNil(tt, r.Error())
	// 	})
	// }
}
*/
