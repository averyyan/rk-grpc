// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-grpc/boot"
)

func main() {
	// Bootstrap basic entries from boot config.
	rkentry.RegisterInternalEntriesFromConfig("example/boot/secure/boot.yaml")

	// Bootstrap grpc entry from boot config
	res := rkgrpc.RegisterGrpcEntriesWithConfig("example/boot/secure/boot.yaml")

	// Bootstrap grpc entry
	res["greeter"].Bootstrap(context.Background())

	// Wait for shutdown signal
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Interrupt grpc entry
	res["greeter"].Interrupt(context.Background())
}