package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ory/dockertest/v3"
)

type TestDocker struct {
	ContainerPort string
	cleanup       []func()
}

// NewTestDocker creates a new TestDocker instance.
func NewTestDocker() *TestDocker {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	var cleanup []func()
	cleanup = append(cleanup, cancel)

	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		panic(fmt.Errorf("%w: docker engine not runnig", err))
	}

	// pulls an image, creates a container based on it and runs it
	dbContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=listmonk",
			"POSTGRES_USER=listmonk",
			"POSTGRES_DB=listmonk",
		},
	})
	if err != nil {
		panic(err)
	}
	// check if container is running
	err = pool.Retry(func() error {
		_, err := pool.Client.InspectContainer(dbContainer.Container.ID)
		if err != nil {
			return err
		}
		exitCode, err := dbContainer.Exec([]string{"pg_isready", "-U", "listmonk"}, dockertest.ExecOptions{Env: []string{"PGPASSWORD=listmonk"}})
		if err != nil {
			return err
		}
		if exitCode != 0 {
			return fmt.Errorf("exit code %d", exitCode)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// // defer cleanup
	// defer pool.Purge(dbContainer)
	cleanup = append(cleanup, func() {
		err := pool.Purge(dbContainer)
		if err != nil {
			fmt.Println("cleanup: container purge err", err)
		}
	})
	// run listmonk container
	listmonkContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "listmonk/listmonk",
		Tag:        "latest",
		Env: []string{
			"LISTMONK_app__address=0.0.0.0:9000",
			"LISTMONK_app__admin_username=listmonk",
			"LISTMONK_app__admin_password=listmonk",
			fmt.Sprintf("LISTMONK_db__host=%s", dbContainer.Container.NetworkSettings.IPAddress),
			"LISTMONK_db__port=5432",
			"LISTMONK_db__user=listmonk",
			"LISTMONK_db__password=listmonk",
			"LISTMONK_db__database=listmonk",
		},
		Cmd: []string{"sh", "-c", "yes | ./listmonk --install && ./listmonk"},
	})
	if err != nil {
		panic(err)
	}
	// check if container is running
	err = pool.Retry(func() error {
		_, err := pool.Client.InspectContainer(listmonkContainer.Container.ID)
		if err != nil {
			return err
		}
		// checking health endpoint with http client to see if listmonk is ready
		httpClient := http.Client{}
		resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/health", listmonkContainer.GetPort("9000/tcp")))
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// defer cleanup
	// defer pool.Purge(listmonkContainer)
	cleanup = append(cleanup, func() {
		err := pool.Purge(listmonkContainer)
		if err != nil {
			fmt.Println("cleanup: container purge err", err)
		}
	})

	return &TestDocker{
		ContainerPort: listmonkContainer.GetPort("9000/tcp"),
		cleanup:       cleanup,
	}
}

// Cleanup cleans up the docker containers.
func (t *TestDocker) Cleanup() {
	for _, f := range t.cleanup {
		f()
	}
}
