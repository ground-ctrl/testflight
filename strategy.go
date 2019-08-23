package main

// Status stores the state of the current simulation: the number of
// periods elapsed, the emitted signals and trades performed.
// The last two are returned by the Backtest to subsequently evaluate
// the strategy.
type Status struct {
	step    int
	signals []Signal
	trades  []Trade
}

func newStatus() *Status {
	return &Status{
		step:    0,
		signals: []Signal{},
		trades:  []Trade{},
	}
}

// A Signal is emitted by the signaler of the strategy. It can be
// one of (buy) "long", (buy) "short", or "sell".
type Signal struct {
	action string
	step   int
}

// A Trade is performed by the Entryer, Exiter and Stopper based on the signals
// emitted by the signaler. It can be one of (buy) "long", (buy) "short", or
// "sell". We trade full positions for now, money management will be added later.
type Trade struct {
	order string
	step  int
}

// A Strategy is a combination of a function that emits signals and functions that
// perform trades. The Entryer enters the market based on the received signals,
// the Stopper cuts losses if necessary, and the exiter exits the market based
// on signals or pre-defined rules.
type Strategy interface {
	Signaler([]float64) string
	Entryer([]float64) string
	Exiter([]float64) string
	Stopper([]float64) string
}

// A Backtest is the combination of a data feed (currently not implemented), a strategy
// to act uponf this data feed, and a rolling status.
type Backtest struct {
	strategy Strategy
	status   *Status
}

func newBacktest(strategy Strategy) *Backtest {
	status := newStatus()
	return &Backtest{
		strategy: strategy,
		status:   status,
	}
}

func (b *Backtest) Run() *Status {
	data := []float64{} // this is where we plug the data stream
	for step, _ := range data {
		if sig := b.strategy.Signaler(data[:step]); sig != "" {
			b.status.signals = append(b.status.signals, Signal{sig, step})
		}

		if entry := b.strategy.Entryer(data[:step]); entry != "" {
			b.status.trades = append(b.status.trades, Trade{entry, step})
		}

		if exit := b.strategy.Exiter(data[:step]); exit != "" {
			b.status.trades = append(b.status.trades, Trade{exit, step})
		}

		if stop := b.strategy.Stopper(data[:step]); stop != "" {
			b.status.trades = append(b.status.trades, Trade{stop, step})
		}
	}

	return b.status
}
