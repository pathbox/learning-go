package internal

const (
	L0_CompactionTrigger     = 4
	L0_SlowdownWritesTrigger = 8
	Write_buffer_size        = 4 << 20
	NumLevels                = 7
	MaxOpenFiles             = 1000
	NumNonTableCacheFiles    = 10
	MaxMemCompactLevel       = 2
	MaxFileSize              = 2 << 20
)
