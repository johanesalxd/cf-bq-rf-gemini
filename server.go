package bqrfgemini

import (
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var (
	clientPool *sync.Pool
	initOnce   sync.Once
	initError  error
)

// init initializes the HTTP function handler for BQRFGemini
func init() {
	functions.HTTP("BQRFGemini", BQRFGemini)
	initOnce.Do(initializePool)
}
