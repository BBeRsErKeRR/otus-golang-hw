package httprouter

import "context"

type Application interface {
	CreateEvent(context.Context, string, string) error
}
