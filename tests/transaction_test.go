package tests

import (
	"path/filepath"
	"testing"

	"github.com/Friends-Of-Noso/NosoData-Go/legacy"
	"gotest.tools/v3/assert"
)

func TestTransactionReadFromFile(t *testing.T) {
	t.Parallel()

	block := &legacy.LegacyBlock{}

	block100_000File := filepath.Join(dataFolder, "100000.blk")

	err := block.ReadFromFile(block100_000File)
	assert.NilError(t, err)

	transaction := block.Transactions[0]

	assert.Equal(t, int32(100_000), transaction.Block)
}

const (
	transactionAsJSON = `{
  "block": 100000,
  "order-id": "OR10dn3p32kljebwume4lk1ec3b61xpnndza3wwd9yuk9clcgcap",
  "orders-count": 1,
  "order-type": "TRFR",
  "id": "tRJ1PhKvn9p4qUQgdNPJC1XQjSJo798zqEMgcwk4ExrnRwNc",
  "timestamp": 1678108260,
  "reference": "PoolPay_GoneFishing",
  "transfer-index": 1,
  "sender": "N3ESwXxCAR4jw3GVHgmKiX9zx1ojWEf",
  "address": "N3ESwXxCAR4jw3GVHgmKiX9zx1ojWEf",
  "receiver": "N3WBskvLhDVoc56kEwLThnz9GxeqwGM",
  "fee": 10064,
  "amount": 100638791,
  "signature": "MEQCIBIUVSFyYbZuxDwm+GNrrrk0WxIGLHGwc+QltWAKKZTYAiAsGs+itbuDuH91HwKLhWooC9vRZtwznUuX+AAhf2Yv7Q=="
}`
)

func TestTransactionAsJSON(t *testing.T) {
	t.Parallel()

	block := &legacy.LegacyBlock{}

	block100_000File := filepath.Join(dataFolder, "100000.blk")

	err := block.ReadFromFile(block100_000File)
	assert.NilError(t, err)

	transaction := block.Transactions[0]

	assert.Equal(t, transactionAsJSON, transaction.AsJSON())
}
