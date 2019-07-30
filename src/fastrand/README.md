As you can see, fastrand.Uint32n scales on multiple CPUs, while rand.Int31n doesn't scale. Their performance is comparable on GOMAXPROCS=1, but fastrand.Uint32n runs 3x faster than rand.Int31n on GOMAXPROCS=2 and 10x faster than rand.Int31n on GOMAXPROCS=4

GOMAXPROCS=4 go test -bench=.