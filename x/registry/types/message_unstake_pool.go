package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnstakePool = "unstake_pool"

var _ sdk.Msg = &MsgUnstakePool{}

func NewMsgUnstakePool(creator string, id uint64, amount uint64) *MsgUnstakePool {
	return &MsgUnstakePool{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgUnstakePool) Route() string {
	return RouterKey
}

func (msg *MsgUnstakePool) Type() string {
	return TypeMsgUnstakePool
}

func (msg *MsgUnstakePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnstakePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnstakePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
