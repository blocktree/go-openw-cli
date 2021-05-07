module github.com/blocktree/go-openw-cli/v2

go 1.12

require (
	github.com/DisposaBoy/JsonConfigReader v0.0.0-20171218180944-5ea4d0ddac55
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-openw-sdk/v2 v2.1.6
	github.com/blocktree/go-owcdrivers v1.2.22 // indirect
	github.com/blocktree/go-owcrypt v1.1.7
	github.com/blocktree/openwallet/v2 v2.1.0
	github.com/bndr/gotabulate v1.1.2
	github.com/google/uuid v1.1.1
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	go.etcd.io/bbolt v1.3.3
	gopkg.in/urfave/cli.v1 v1.20.0
)

//replace github.com/blocktree/go-openw-sdk/v2 => ../go-openw-sdk

//replace github.com/blocktree/openwallet => ../openwallet

//replace github.com/blocktree/go-owcdrivers => ../go-owcdrivers
