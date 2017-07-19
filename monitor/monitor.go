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
	start bool
}

func parseMonitorArgs(opts []string) (*monitorOpts, error) {
	args := &monitorOpts{}
	fs := flag.NewFlagSet("monitor", flag.ContinueOnError)
	fs.BoolVar(&args.start, "start", false, "start the coin daemon if it is not running")
	if err := fs.Parse(opts); err != nil {
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
	cStateInit                        = iota
	cStateWaitStart                   = iota
	cStateMasternodePendingActivation = iota
)

func (m *Monitor) Start() error {
	// TODO: Check to see if coin is running, if it is we are ok. If not and
	//       we do not have the --start option set, then we kick off the
	// 		 coin's daemon and wait for it to get going.
	if m.Opts.start {
		fmt.Printf("Starting daemon!\n")
		if err := m.Coin.StartDaemon(); err != nil {
			// TODO: if err is already started, dont bail
			return err
		}
	}

	for {

		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("Monitoring coin %s\n", m.Coin.GetName())
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
