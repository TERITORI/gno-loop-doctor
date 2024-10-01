package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	indexerBinary := flag.String("indexer-bin", "/gno-tx-indexer", "path to indexer binary")
	dbPath := flag.String("db-path", "/indexer-db", "path to db dir")

	remote := "https://rpc.gno.land"
	prevTime := ""

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	var prevCmd *exec.Cmd

	jobsCtx, cancelJobsCtx := context.WithCancel(context.Background())

	// FIXME: races

	_, err = s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				err := func() error {
					t, err := getGenesisTime(remote)
					if err != nil {
						return fmt.Errorf("failed to get genesis time: %w", err)
					}

					if t != prevTime {
						logger.Info("restarting", zap.String("genesis-time", t), zap.String("prev-time", prevTime))

						if prevCmd != nil {
							if err := prevCmd.Process.Kill(); err != nil {
								return fmt.Errorf("failed to stop previous instance: %w", err)
							}
							if err := prevCmd.Wait(); err != nil {
								return fmt.Errorf("failed to wait end of previous instance: %w", err)
							}
							logger.Info("stopped previous instance")
						}

						// clear db
						if err := os.RemoveAll(*dbPath); err != nil {
							return fmt.Errorf("failed to clear db: %w", err)
						}

						logger.Info("cleared db")

						// start indexer
						indexer := exec.CommandContext(jobsCtx, *indexerBinary, "start", "-remote", remote, "-db-path", *dbPath)
						indexer.WaitDelay = 10 * time.Second
						indexer.Stdout = os.Stdout
						indexer.Stderr = os.Stderr
						prevCmd = indexer
						prevTime = t
						if err := indexer.Run(); err != nil {
							return fmt.Errorf("failed to run indexer: %w", err)
						}
					}

					return nil
				}()
				if err != nil {
					logger.Error("job failed", zap.Error(err))
				}
			},
		),
	)
	if err != nil {
		panic(err)
	}

	s.Start()
	<-sigCh
	cancelJobsCtx()
	s.Shutdown()
}

func getGenesisTime(remote string) (string, error) {
	res, err := http.Get(remote + "/genesis")
	if err != nil {
		return "", fmt.Errorf("failed to get genesis: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var data genesisRes
	if err := decoder.Decode(&data); err != nil {
		return "", fmt.Errorf("failed to unmarshal genesis time: %w", err)
	}

	return data.Result.Genesis.GenesisTime, nil
}

type genesisRes struct {
	Result struct {
		Genesis struct {
			GenesisTime string `json:"genesis_time"`
		} `json:"genesis"`
	} `json:"result"`
}
