package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawPool = "withdraw_pool"

var _ sdk.Msg = &MsgWithdrawPool{}

func NewMsgWithdrawPool(creator string, id uint64, staker string) *MsgWithdrawPool {
	return &MsgWithdrawPool{
		Creator: creator,
		Id:      id,
		Staker:  staker,
	}
}

func (msg *MsgWithdrawPool) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawPool) Type() string {
	return TypeMsgWithdrawPool
}

func (msg *MsgWithdrawPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawPool) ValidateBasic() error {
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
