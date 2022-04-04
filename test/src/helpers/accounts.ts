import { KyveSDK, KyveWallet } from "@kyve/sdk-test";

export const ADDRESS_ALICE = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd";
export const MNEMONIC_ALICE =
  "worry grief loyal smoke pencil arrow trap focus high pioneer tomato hedgehog essence purchase dove pond knee custom phone gentle sunset addict mother fabric";
export const alice = new KyveSDK(new KyveWallet("local", MNEMONIC_ALICE));

export const ADDRESS_BOB = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq";
export const MNEMONIC_BOB =
  "crash sick toilet stumble join cash erode glory door weird diagram away lizard solid segment apple urge joy annual able tank define candy demise";
export const bob = new KyveSDK(new KyveWallet("local", MNEMONIC_BOB));

export const ADDRESS_CHARLIE = "kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s";
export const MNEMONIC_CHARLIE =
  "shoot inject fragile width trend satisfy army enact volcano crowd message strike true divorce search rich office shoulder sport relax rhythm symbol gadget size";
export const charlie = new KyveSDK(new KyveWallet("local", MNEMONIC_CHARLIE));
