package broker

import "github.com/mdshahjahanmiah/banking-ledger/model"

type Producer interface {
	PublishTransaction(txn model.Transaction) error
}
