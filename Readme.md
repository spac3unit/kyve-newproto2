# Chain-Node
###### v0.0.1


The Chain-Nodes are the backbone of KYVE. The chain layer is a
completely sovereign [Proof of Stake](https://en.wikipedia.org/wiki/Proof_of_stake)
blockchain build with [Starport](https://starport.com/). This
blockchain is run by independent nodes we call _Chain-Nodes_
since they're running on the chain level. The native currency
of the KYVE chain is $[KYVE](https://docs.kyve.network/basics/kyve.html), it secures the chain
and allows chain nodes to stake and other users to delegate into them.

---

## Building from source

Currently, building from source is only supported with [Starport](https://starport.com/)

```bash
starport chain build --release --release.prefix chain
```
The output can be found in ./release

If you need to build for different architectures, use the `-t` flag, e.g. `-t linux:amd64,linux:arm64`

## Running a chain node

The full introductions for setting up a node are provided at:
[https://docs.kyve.network/intro/chain-node.html](https://docs.kyve.network/intro/chain-node.html)
