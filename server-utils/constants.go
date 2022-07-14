package server_utils

const (
	DONGLE = 0
	BEACON = 1
)

const (
	RISK = 0
	EPI  = 1
)

const MAX_CHUNK_EPHID_COUNT = 512

const MINUTES_IN_14_DAYS = (14*24*60)
const EXPONENT_TOO_LARGE = 24

const BROADCAST_SERVICE_UUID = 0x2222
const BEACON_SERVICE_ID_MASK = 0x0000ffff
const MAX_BEACON_IDS = BEACON_SERVICE_ID_MASK + 1
