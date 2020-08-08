package main

import (
	"log"
	"os"

	"hscan/client"
	"hscan/config"
	"hscan/db"
	"hscan/scanner"
	"hscan/schema"

	//tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
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
	db.AutoMigrate(&schema.Block{})

	cdc := newCodec()

	scanner := scanner.NewScanner(l, client, db, cdc)

	scanner.Start()

	// if b, err := client.GetBlock(1); err != nil {
	// 	l.Printf("err is %s", err)
	// } else {
	// 	l.Printf("block is %+v", *b)
	// }

}
