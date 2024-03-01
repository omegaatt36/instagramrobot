package health

import (
	"math"
)

// InstanceIDNotFound means tenant.signal-ctrl may not find the instance(signal/strategy producer).
const InstanceIDNotFound int32 = math.MaxInt32
