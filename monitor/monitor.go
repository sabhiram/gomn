package monitor

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"fmt"
	"time"

	"github.com/sabhiram/gomn/coin"
	"github.com/sabhiram/gomn/types"
)

////////////////////////////////////////////////////////////////////////////////

type monitorOpts struct {
	start              bool
	refreshIntervalStr string
	refreshInterval    time.Duration
}

func parseMonitorArgs(opts []string) (*monitorOpts, error) {
	args := &monitorOpts{}
	fs := flag.NewFlagSet("monitor", flag.ContinueOnError)
	fs.BoolVar(&args.start, "start", false, "start the coin daemon if it is not running")
	fs.StringVar(&args.refreshIntervalStr, "refresh", "30s", "refresh interval, default 30s")
	if err := fs.Parse(opts); err != nil {
		return nil, err
	}

	var err error
	args.refreshInterval, err = time.ParseDuration(args.refreshIntervalStr)
	if err != nil {
		return nil, err
	}
	return args, nil
}

////////////////////////////////////////////////////////////////////////////////

type Monitor struct {
	Coin *coin.Coin
	CLI  *types.CLI
	Opts *monitorOpts
}

func New(cli *types.CLI, opts []string) (*Monitor, error) {
	// Parse monitor specific arguments.
	mopts, err := parseMonitorArgs(opts)
	if err != nil {
		return nil, err
	}

	// What coin are we monitoring?
	c, err := coin.GetCoinByName(cli.Coin)
	if err != nil {
		return nil, err
	}

	c.UpdateDynamic(cli.Wallet, cli.BinPath, cli.DataPath)

	// TODO: Verify that coin is correctly setup

	return &Monitor{
		Coin: c,
		CLI:  cli,
		Opts: mopts,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

const (
	cStateInit                        = iota // Initial state
	cStateWaitStart                   = iota // Waiting for the daemon to startup
	cStateWaitMasternode              = iota // Waiting for masternode to come up
	cStateMasternodePendingActivation = iota // Masternode is pending client activation
	cStateMasternodeRunning           = iota // Masternode running fine and dandy
	cStateNotCapbaleMasternode        = iota // Error state: node not capable of being a masternode
)

func (m *Monitor) Start() error {
	daemonRunning := false
	if err := m.Coin.FnMap.GetInfoFn(m.Coin, nil); err == nil {
		daemonRunning = true
	}

	if !daemonRunning && !m.Opts.start {
		return fmt.Errorf("%s daemon is not running, try adding the --start option to monitor", m.Coin.GetName())
	} else if !daemonRunning {
		// Attempt to start the coin's daemon.
		fmt.Printf("Trying to start %s daemon\n", m.Coin.GetName())
		if err := m.Coin.StartDaemon(); err != nil {
			return err
		}
		fmt.Printf("... started at %s\n", time.Now().String())
	}

	for {
		select {
		case <-time.After(m.Opts.refreshInterval):
			fmt.Printf("Monitoring coin %s\n", m.Coin.GetName())
			err := m.Coin.FnMap.GetInfoFn(m.Coin, nil)
			if err != nil {
				fmt.Printf("Warning: Coin daemon down? : %s\n", err.Error())
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
