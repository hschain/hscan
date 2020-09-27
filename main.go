package main

import (
	"log"
	"os"

	"hscan/client"
	"hscan/config"
	"hscan/db"
	"hscan/scanner"
	"hscan/schema"
	"hscan/server"

	//tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	sdk "github.com/hschain/hschain/types"

	codec "github.com/hschain/hschain/codec"
	"github.com/hschain/hschain/types/module"
	"github.com/hschain/hschain/x/auth"
	"github.com/hschain/hschain/x/bank"
	"github.com/hschain/hschain/x/crisis"
	distr "github.com/hschain/hschain/x/distribution"
	"github.com/hschain/hschain/x/genaccounts"
	"github.com/hschain/hschain/x/genutil"
	"github.com/hschain/hschain/x/gov"
	"github.com/hschain/hschain/x/mint"
	"github.com/hschain/hschain/x/params"
	paramsclient "github.com/hschain/hschain/x/params/client"
	"github.com/hschain/hschain/x/slashing"
	"github.com/hschain/hschain/x/staking"
	"github.com/hschain/hschain/x/supply"
)

var (
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
	)
)

func newCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

func main() {
	l := log.New(os.Stdout, "Hschain scanner", log.Lshortfile|log.LstdFlags)

	cfg := config.ParseConfig()
	l.Printf("config is %+v", *cfg)

	client := client.NewClient(
		cfg.Node,
	)

	db := db.NewDB(cfg.Mysql)
	db.LogMode(true)
	db.AutoMigrate(&schema.Block{}, &schema.Transaction{}, &schema.NodeInfo{})

	cdc := newCodec()

	scanner := scanner.NewScanner(l, client, db, cdc)

	scanner.Start()

	server := server.NewServer("127.0.0.1:"+cfg.Web.Port, l, db, cdc, client)

	server.Start()

}
