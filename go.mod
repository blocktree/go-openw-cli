module github.com/blocktree/go-openw-cli

go 1.12

require (
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-openw-sdk v1.5.0
	github.com/blocktree/go-owcdrivers v1.2.0
	github.com/blocktree/openwallet v1.7.0
	github.com/bndr/gotabulate v1.1.2
	github.com/google/uuid v1.1.1
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	go.etcd.io/bbolt v1.3.3
	gopkg.in/urfave/cli.v1 v1.20.0
)

//replace github.com/blocktree/go-openw-sdk => ../go-openw-sdk

//replace github.com/blocktree/openwallet => ../openwallet

//replace github.com/blocktree/go-owcdrivers => ../go-owcdrivers
