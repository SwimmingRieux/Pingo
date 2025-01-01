package abstraction

import "context"

type NetworkLogRecorder interface {
	Record(context context.Context)
}
