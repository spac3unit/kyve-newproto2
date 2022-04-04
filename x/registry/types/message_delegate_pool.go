package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDelegatePool = "delegate_pool"

var _ sdk.Msg = &MsgDelegatePool{}

func NewMsgDelegatePool(creator string, id uint64, staker string, amount uint64) *MsgDelegatePool {
	return &MsgDelegatePool{
		Creator: creator,
		Id:      id,
		Staker:  staker,
		Amount:  amount,
	}
}

func (msg *MsgDelegatePool) Route() string {
	return RouterKey
}

func (msg *MsgDelegatePool) Type() string {
	return TypeMsgDelegatePool
}

func (msg *MsgDelegatePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDelegatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDelegatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Staker)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid staker address (%s)", err)
	}
	return nil
}
