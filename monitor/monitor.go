package monitor

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"time"

	"github.com/sabhiram/gomn/coin"
	"github.com/sabhiram/gomn/types"
)

////////////////////////////////////////////////////////////////////////////////

type Monitor struct {
	Coin *coin.Coin
	CLI  *types.CLI
}

func New(cli *types.CLI, opts []string) (*Monitor, error) {
	// What coin are we monitoring?
	c, err := coin.GetCoinByName(cli.Coin)
	if err != nil {
		return nil, err
	}

	// TODO: Verify that coin is correctly setup

	// TODO: Check to see if coin is running, if it is we are ok. If not and
	//       we do not have the --start option set, then we kick off the
	// 		 coin's daemon and wait for it to get going.

	return &Monitor{
		Coin: c,
		CLI:  cli,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

func (m *Monitor) Start() error {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("Monitoring coin %s\n", m.Coin.GetName())
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
