package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakePool = "stake_pool"

var _ sdk.Msg = &MsgStakePool{}

func NewMsgStakePool(creator string, id uint64, amount uint64) *MsgStakePool {
	return &MsgStakePool{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgStakePool) Route() string {
	return RouterKey
}

func (msg *MsgStakePool) Type() string {
	return TypeMsgStakePool
}

func (msg *MsgStakePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStakePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStakePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
