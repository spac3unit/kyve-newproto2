package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUndelegatePool = "undelegate_pool"

var _ sdk.Msg = &MsgUndelegatePool{}

func NewMsgUndelegatePool(creator string, id uint64, staker string, amount uint64) *MsgUndelegatePool {
	return &MsgUndelegatePool{
		Creator: creator,
		Id:      id,
		Staker:  staker,
		Amount:  amount,
	}
}

func (msg *MsgUndelegatePool) Route() string {
	return RouterKey
}

func (msg *MsgUndelegatePool) Type() string {
	return TypeMsgUndelegatePool
}

func (msg *MsgUndelegatePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUndelegatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUndelegatePool) ValidateBasic() error {
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
