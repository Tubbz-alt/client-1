// Auto-generated by avdl-compiler v1.3.22 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/stellar1/remote.avdl

package stellar1

import (
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type AccountID string

func (o AccountID) DeepCopy() AccountID {
	return o
}

type TransactionID string

func (o TransactionID) DeepCopy() TransactionID {
	return o
}

type KeybaseTransactionID []byte

func (o KeybaseTransactionID) DeepCopy() KeybaseTransactionID {
	return (func(x []byte) []byte {
		if x == nil {
			return nil
		}
		return append([]byte{}, x...)
	})(o)
}

type Asset struct {
	Type   string `codec:"type" json:"type"`
	Code   string `codec:"code" json:"code"`
	Issuer string `codec:"issuer" json:"issuer"`
}

func (o Asset) DeepCopy() Asset {
	return Asset{
		Type:   o.Type,
		Code:   o.Code,
		Issuer: o.Issuer,
	}
}

type Balance struct {
	Asset  Asset  `codec:"asset" json:"asset"`
	Amount string `codec:"amount" json:"amount"`
	Limit  string `codec:"limit" json:"limit"`
}

func (o Balance) DeepCopy() Balance {
	return Balance{
		Asset:  o.Asset.DeepCopy(),
		Amount: o.Amount,
		Limit:  o.Limit,
	}
}

type EncryptedNote struct {
	V   int          `codec:"v" json:"v"`
	E   []byte       `codec:"e" json:"e"`
	N   []byte       `codec:"n" json:"n"`
	KID keybase1.KID `codec:"KID" json:"KID"`
}

func (o EncryptedNote) DeepCopy() EncryptedNote {
	return EncryptedNote{
		V: o.V,
		E: (func(x []byte) []byte {
			if x == nil {
				return nil
			}
			return append([]byte{}, x...)
		})(o.E),
		N: (func(x []byte) []byte {
			if x == nil {
				return nil
			}
			return append([]byte{}, x...)
		})(o.N),
		KID: o.KID.DeepCopy(),
	}
}

type TransactionSummary struct {
	StellarID   TransactionID        `codec:"stellarID" json:"stellarID"`
	KeybaseID   KeybaseTransactionID `codec:"keybaseID" json:"keybaseID"`
	Note        EncryptedNote        `codec:"note" json:"note"`
	Asset       Asset                `codec:"asset" json:"asset"`
	Amount      string               `codec:"amount" json:"amount"`
	StellarFrom AccountID            `codec:"stellarFrom" json:"stellarFrom"`
	StellarTo   AccountID            `codec:"stellarTo" json:"stellarTo"`
	KeybaseFrom keybase1.UID         `codec:"keybaseFrom" json:"keybaseFrom"`
	KeybaseTo   keybase1.UID         `codec:"keybaseTo" json:"keybaseTo"`
}

func (o TransactionSummary) DeepCopy() TransactionSummary {
	return TransactionSummary{
		StellarID:   o.StellarID.DeepCopy(),
		KeybaseID:   o.KeybaseID.DeepCopy(),
		Note:        o.Note.DeepCopy(),
		Asset:       o.Asset.DeepCopy(),
		Amount:      o.Amount,
		StellarFrom: o.StellarFrom.DeepCopy(),
		StellarTo:   o.StellarTo.DeepCopy(),
		KeybaseFrom: o.KeybaseFrom.DeepCopy(),
		KeybaseTo:   o.KeybaseTo.DeepCopy(),
	}
}

type Operation struct {
	ID              string       `codec:"ID" json:"ID"`
	OpType          string       `codec:"opType" json:"opType"`
	CreatedAt       int          `codec:"createdAt" json:"createdAt"`
	TransactionHash string       `codec:"TransactionHash" json:"TransactionHash"`
	Asset           Asset        `codec:"asset" json:"asset"`
	Amount          string       `codec:"amount" json:"amount"`
	StellarFrom     AccountID    `codec:"stellarFrom" json:"stellarFrom"`
	StellarTo       AccountID    `codec:"stellarTo" json:"stellarTo"`
	KeybaseFrom     keybase1.UID `codec:"keybaseFrom" json:"keybaseFrom"`
	KeybaseTo       keybase1.UID `codec:"keybaseTo" json:"keybaseTo"`
}

func (o Operation) DeepCopy() Operation {
	return Operation{
		ID:              o.ID,
		OpType:          o.OpType,
		CreatedAt:       o.CreatedAt,
		TransactionHash: o.TransactionHash,
		Asset:           o.Asset.DeepCopy(),
		Amount:          o.Amount,
		StellarFrom:     o.StellarFrom.DeepCopy(),
		StellarTo:       o.StellarTo.DeepCopy(),
		KeybaseFrom:     o.KeybaseFrom.DeepCopy(),
		KeybaseTo:       o.KeybaseTo.DeepCopy(),
	}
}

type TransactionDetails struct {
	StellarID             TransactionID        `codec:"stellarID" json:"stellarID"`
	KeybaseID             KeybaseTransactionID `codec:"keybaseID" json:"keybaseID"`
	Hash                  string               `codec:"Hash" json:"Hash"`
	Ledger                int                  `codec:"ledger" json:"ledger"`
	LedgerCloseTime       int                  `codec:"ledgerCloseTime" json:"ledgerCloseTime"`
	SourceAccount         AccountID            `codec:"sourceAccount" json:"sourceAccount"`
	SourceAccountSequence string               `codec:"sourceAccountSequence" json:"sourceAccountSequence"`
	FeePaid               int                  `codec:"feePaid" json:"feePaid"`
	Note                  EncryptedNote        `codec:"note" json:"note"`
	Signatures            []string             `codec:"signatures" json:"signatures"`
	Operations            []Operation          `codec:"operations" json:"operations"`
}

func (o TransactionDetails) DeepCopy() TransactionDetails {
	return TransactionDetails{
		StellarID:             o.StellarID.DeepCopy(),
		KeybaseID:             o.KeybaseID.DeepCopy(),
		Hash:                  o.Hash,
		Ledger:                o.Ledger,
		LedgerCloseTime:       o.LedgerCloseTime,
		SourceAccount:         o.SourceAccount.DeepCopy(),
		SourceAccountSequence: o.SourceAccountSequence,
		FeePaid:               o.FeePaid,
		Note:                  o.Note.DeepCopy(),
		Signatures: (func(x []string) []string {
			if x == nil {
				return nil
			}
			var ret []string
			for _, v := range x {
				vCopy := v
				ret = append(ret, vCopy)
			}
			return ret
		})(o.Signatures),
		Operations: (func(x []Operation) []Operation {
			if x == nil {
				return nil
			}
			var ret []Operation
			for _, v := range x {
				vCopy := v.DeepCopy()
				ret = append(ret, vCopy)
			}
			return ret
		})(o.Operations),
	}
}

type PaymentPost struct {
	From              AccountID     `codec:"from" json:"from"`
	To                AccountID     `codec:"to" json:"to"`
	SequenceNumber    uint64        `codec:"sequenceNumber" json:"sequenceNumber"`
	KeybaseFrom       string        `codec:"keybaseFrom" json:"keybaseFrom"`
	KeybaseTo         string        `codec:"keybaseTo" json:"keybaseTo"`
	TransactionType   string        `codec:"transactionType" json:"transactionType"`
	XlmAmount         string        `codec:"xlmAmount" json:"xlmAmount"`
	DisplayAmount     string        `codec:"displayAmount" json:"displayAmount"`
	DisplayCurrency   string        `codec:"displayCurrency" json:"displayCurrency"`
	Note              EncryptedNote `codec:"note" json:"note"`
	SignedTransaction string        `codec:"signedTransaction" json:"signedTransaction"`
}

func (o PaymentPost) DeepCopy() PaymentPost {
	return PaymentPost{
		From:              o.From.DeepCopy(),
		To:                o.To.DeepCopy(),
		SequenceNumber:    o.SequenceNumber,
		KeybaseFrom:       o.KeybaseFrom,
		KeybaseTo:         o.KeybaseTo,
		TransactionType:   o.TransactionType,
		XlmAmount:         o.XlmAmount,
		DisplayAmount:     o.DisplayAmount,
		DisplayCurrency:   o.DisplayCurrency,
		Note:              o.Note.DeepCopy(),
		SignedTransaction: o.SignedTransaction,
	}
}

type PaymentResult struct {
	StellarID TransactionID        `codec:"stellarID" json:"stellarID"`
	KeybaseID KeybaseTransactionID `codec:"keybaseID" json:"keybaseID"`
	Ledger    int                  `codec:"Ledger" json:"Ledger"`
}

func (o PaymentResult) DeepCopy() PaymentResult {
	return PaymentResult{
		StellarID: o.StellarID.DeepCopy(),
		KeybaseID: o.KeybaseID.DeepCopy(),
		Ledger:    o.Ledger,
	}
}

type BalancesArg struct {
	Uid       keybase1.UID `codec:"uid" json:"uid"`
	AccountID AccountID    `codec:"accountID" json:"accountID"`
}

type RecentTransactionsArg struct {
	Uid       keybase1.UID `codec:"uid" json:"uid"`
	AccountID AccountID    `codec:"accountID" json:"accountID"`
	Count     int          `codec:"count" json:"count"`
}

type TransactionArg struct {
	Uid keybase1.UID  `codec:"uid" json:"uid"`
	Id  TransactionID `codec:"id" json:"id"`
}

type SubmitPaymentArg struct {
	Uid     keybase1.UID `codec:"uid" json:"uid"`
	Payment PaymentPost  `codec:"payment" json:"payment"`
}

type WalletInitArg struct {
}

type RemoteInterface interface {
	Balances(context.Context, BalancesArg) ([]Balance, error)
	RecentTransactions(context.Context, RecentTransactionsArg) ([]TransactionSummary, error)
	Transaction(context.Context, TransactionArg) (TransactionDetails, error)
	SubmitPayment(context.Context, SubmitPaymentArg) (PaymentResult, error)
	WalletInit(context.Context) error
}

func RemoteProtocol(i RemoteInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "stellar.1.remote",
		Methods: map[string]rpc.ServeHandlerDescription{
			"balances": {
				MakeArg: func() interface{} {
					ret := make([]BalancesArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]BalancesArg)
					if !ok {
						err = rpc.NewTypeError((*[]BalancesArg)(nil), args)
						return
					}
					ret, err = i.Balances(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"recentTransactions": {
				MakeArg: func() interface{} {
					ret := make([]RecentTransactionsArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]RecentTransactionsArg)
					if !ok {
						err = rpc.NewTypeError((*[]RecentTransactionsArg)(nil), args)
						return
					}
					ret, err = i.RecentTransactions(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"transaction": {
				MakeArg: func() interface{} {
					ret := make([]TransactionArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]TransactionArg)
					if !ok {
						err = rpc.NewTypeError((*[]TransactionArg)(nil), args)
						return
					}
					ret, err = i.Transaction(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"submitPayment": {
				MakeArg: func() interface{} {
					ret := make([]SubmitPaymentArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]SubmitPaymentArg)
					if !ok {
						err = rpc.NewTypeError((*[]SubmitPaymentArg)(nil), args)
						return
					}
					ret, err = i.SubmitPayment(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"walletInit": {
				MakeArg: func() interface{} {
					ret := make([]WalletInitArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					err = i.WalletInit(ctx)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type RemoteClient struct {
	Cli rpc.GenericClient
}

func (c RemoteClient) Balances(ctx context.Context, __arg BalancesArg) (res []Balance, err error) {
	err = c.Cli.Call(ctx, "stellar.1.remote.balances", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) RecentTransactions(ctx context.Context, __arg RecentTransactionsArg) (res []TransactionSummary, err error) {
	err = c.Cli.Call(ctx, "stellar.1.remote.recentTransactions", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) Transaction(ctx context.Context, __arg TransactionArg) (res TransactionDetails, err error) {
	err = c.Cli.Call(ctx, "stellar.1.remote.transaction", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) SubmitPayment(ctx context.Context, __arg SubmitPaymentArg) (res PaymentResult, err error) {
	err = c.Cli.Call(ctx, "stellar.1.remote.submitPayment", []interface{}{__arg}, &res)
	return
}

func (c RemoteClient) WalletInit(ctx context.Context) (err error) {
	err = c.Cli.Call(ctx, "stellar.1.remote.walletInit", []interface{}{WalletInitArg{}}, nil)
	return
}
