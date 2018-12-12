package openwcli

import (
	"github.com/asdine/storm"
	"github.com/blocktree/OpenWallet/openwallet"
	"time"
)

func getKeyPair() {
	db, err = openwallet.OpenStormDB(
		wm.DBFile(appID),
		storm.Batch(),
		storm.BoltOptions(0600, &bolt.Options{Timeout: 3 * time.Second}),
	)
	log.Debug("open storm db appID:", appID)
	if err != nil {
		return nil, err
	}
}