package cli

import (
	"io"
	"os"

	"github.com/huanghao/dida365-cli/internal/config"
)

type App struct {
	In          io.Reader
	Out         io.Writer
	Err         io.Writer
	Debug       bool
	DryRun      bool
	ConfigPath  string
	ConfigStore *config.Store
}

func NewApp() (*App, error) {
	store, err := config.NewStore("")
	if err != nil {
		return nil, err
	}
	return &App{
		In:          os.Stdin,
		Out:         os.Stdout,
		Err:         os.Stderr,
		ConfigStore: store,
	}, nil
}

func (a *App) ReloadConfigStore() error {
	store, err := config.NewStore(a.ConfigPath)
	if err != nil {
		return err
	}
	a.ConfigStore = store
	return nil
}
