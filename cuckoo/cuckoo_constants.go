package cuckoo

// implementation specific constants
const (
	FINGERPRINT_BITS = 27
	SHIFT_LEEWAY     = 5 // distance to next multiple of 8 from FINGERPRINT_BITS
	FINGERPRINT_SIZE = 4 // ceiling of (FINGERPRINT_BITS / 8)
	BUCKET_SIZE      = 4
	MAX_EVICTIONS    = 500
)

// test constants
const (
	TEST_NUM_BUCKETS      = 32
	TEST_ZERO_NUM_BUCKETS = 0
	TEST_BAD_NUM_BUCKETS  = 31
)
