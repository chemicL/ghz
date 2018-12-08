package ghz

import (
	"os"
	"os/signal"
	"time"

	"github.com/bojand/ghz/protodesc"

	"github.com/jhump/protoreflect/desc"
)

// Run executes the test
func Run(call, host string, options ...Option) (*Report, error) {
	c, err := newConfig(call, host, options...)

	if err != nil {
		return nil, err
	}

	var mtd *desc.MethodDescriptor
	if c.proto != "" {
		mtd, err = protodesc.GetMethodDescFromProto(call, c.proto, c.importPaths)
	} else {
		mtd, err = protodesc.GetMethodDescFromProtoSet(call, c.protoset)
	}

	if err != nil {
		return nil, err
	}

	reqr, err := newRequester(mtd, c)

	if err != nil {
		return nil, err
	}

	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt)
	go func() {
		<-cancel
		reqr.Stop(ReasonCancel)
	}()

	if c.z > 0 {
		go func() {
			time.Sleep(c.z)
			reqr.Stop(ReasonTimeout)
		}()
	}

	return reqr.Run()
}