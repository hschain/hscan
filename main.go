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
	sdk "github.com/zxs-paryada/hschain/types"

	codec "github.com/zxs-paryada/hschain/codec"
	"github.com/zxs-paryada/hschain/types/module"
	"github.com/zxs-paryada/hschain/x/auth"
	"github.com/zxs-paryada/hschain/x/bank"
	"github.com/zxs-paryada/hschain/x/crisis"
	distr "github.com/zxs-paryada/hschain/x/distribution"
	"github.com/zxs-paryada/hschain/x/genaccounts"
	"github.com/zxs-paryada/hschain/x/genutil"
	"github.com/zxs-paryada/hschain/x/gov"
	"github.com/zxs-paryada/hschain/x/mint"
	"github.com/zxs-paryada/hschain/x/params"
	paramsclient "github.com/zxs-paryada/hschain/x/params/client"
	"github.com/zxs-paryada/hschain/x/slashing"
	"github.com/zxs-paryada/hschain/x/staking"
	"github.com/zxs-paryada/hschain/x/supply"
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
	db.AutoMigrate(&schema.Block{}, &schema.Transaction{})

	cdc := newCodec()

	scanner := scanner.NewScanner(l, client, db, cdc)

	scanner.Start()

	server := server.NewServer("127.0.0.1:"+cfg.Web.Port, l, db, cdc, client)

	server.Start()

}
