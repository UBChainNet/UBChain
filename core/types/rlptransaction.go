package types

import (
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/core/types/functionbody/tokenhub_func"
)

type RlpTransaction struct {
	TxHead *TransactionHead
	TxBody []byte
}

type RlpContract struct {
	TxHead *TransactionHead
	TxBody RlpContractBody
}

type RlpContractBody struct {
	Contract     hasharry.Address
	Type         contractv2.ContractType
	FunctionType contractv2.FunctionType
	Function     []byte
	State        ContractState
	Message      []byte
}

func (rt *RlpTransaction) TranslateToTransaction() *Transaction {
	switch rt.TxHead.TxType {
	case Transfer_:
		var nt *TransferBody
		rlp.DecodeBytes(rt.TxBody, &nt)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: nt,
		}
	case Contract_:
		var ct *ContractBody
		rlp.DecodeBytes(rt.TxBody, &ct)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: ct,
		}
	case ContractV2_:
		var ct = &TxContractV2Body{}
		var rlpCt *RlpContractBody
		rlp.DecodeBytes(rt.TxBody, &rlpCt)
		switch rlpCt.FunctionType {
		case contractv2.Exchange_Init:
			var init *exchange_func.ExchangeInitBody
			rlp.DecodeBytes(rlpCt.Function, &init)
			ct.Function = init
		case contractv2.Exchange_SetAdmin:
			var set *exchange_func.ExchangeAdmin
			rlp.DecodeBytes(rlpCt.Function, &set)
			ct.Function = set
		case contractv2.Exchange_SetFeeTo:
			var set *exchange_func.ExchangeFeeTo
			rlp.DecodeBytes(rlpCt.Function, &set)
			ct.Function = set
		case contractv2.Exchange_ExactIn:
			var in *exchange_func.ExactIn
			rlp.DecodeBytes(rlpCt.Function, &in)
			ct.Function = in
		case contractv2.Exchange_ExactOut:
			var out *exchange_func.ExactOut
			rlp.DecodeBytes(rlpCt.Function, &out)
			ct.Function = out
		case contractv2.Pair_AddLiquidity:
			var create *exchange_func.ExchangeAddLiquidity
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pair_RemoveLiquidity:
			var create *exchange_func.ExchangeRemoveLiquidity
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_Init:
			var create *exchange_func.PledgeInitBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_Start:
			var create *exchange_func.PledgeStartBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_AddPool:
			var create *exchange_func.PledgeAddPoolBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_RemovePool:
			var create *exchange_func.PledgeRemovePoolBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_Add:
			var create *exchange_func.PledgeAddBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_Remove:
			var create *exchange_func.PledgeRemoveBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_RemoveReward:
			var create *exchange_func.PledgeRewardRemoveBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.Pledge_Update:
			var create *exchange_func.PledgeUpdateBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.TokenHub_init:
			var create *tokenhub_func.TokenHubInitBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.TokenHub_Ack:
			var create *tokenhub_func.TokenHubAckBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.TokenHub_TransferOut:
			var create *tokenhub_func.TokenHubTransferOutBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.TokenHub_TransferIn:
			var create *tokenhub_func.TokenHubTransferInBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		case contractv2.TokenHub_FinishAcross:
			var create *tokenhub_func.TokenHubFinishAcrossBody
			rlp.DecodeBytes(rlpCt.Function, &create)
			ct.Function = create
		}
		rlp.DecodeBytes(rt.TxBody, &ct)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: ct,
		}
	case LoginCandidate_:
		var nt *LoginTransactionBody
		rlp.DecodeBytes(rt.TxBody, &nt)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: nt,
		}
		/*case LogoutCandidate:
			return &Transaction{
				TxHead: rt.TxHead,
				TxBody: &LogoutTransactionBody{},
			}
		case VoteToCandidate:
			var nt *VoteTransactionBody
			rlp.DecodeBytes(rt.TxBody, &nt)
			return &Transaction{
				TxHead: rt.TxHead,
				TxBody: nt,
			}*/
	}
	return nil
}
