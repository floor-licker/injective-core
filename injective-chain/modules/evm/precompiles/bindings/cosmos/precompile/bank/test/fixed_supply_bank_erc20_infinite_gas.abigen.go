// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bank

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FixedSupplyBankERC20InfiniteGasMetaData contains all meta data concerning the FixedSupplyBankERC20InfiniteGas contract.
var FixedSupplyBankERC20InfiniteGasMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"initial_supply_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Failure\",\"inputs\":[{\"name\":\"message\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x6080604052606460055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015610050575f5ffd5b50604051612402380380612402833981810160405281019061007291906107bd565b8383838383838360405180602001604052805f81525060405180602001604052805f81525081600390816100a69190610a60565b5080600490816100b69190610a60565b5050505f835111806100c857505f8251115b806100d557505f8160ff16115b1561017a5760055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166337d2c2f48484846040518463ffffffff1660e01b815260040161013893929190610b86565b6020604051808303815f875af1158015610154573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101789190610bfe565b505b5050505f8111156101965761019533826101a360201b60201c565b5b5050505050505050610eaf565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610213575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161020a9190610c68565b60405180910390fd5b6102245f838361022860201b60201c565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361040d5760055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f1983836040518363ffffffff1660e01b81526004016102b7929190610c90565b6020604051808303815f875af19250505080156102f257506040513d601f19601f820116820180604052508101906102ef9190610bfe565b60015b610407576102fe610cc3565b806308c379a0036103795750610312610ce2565b8061031d575061037b565b8060405160200161032e9190610dd1565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103709190610df6565b60405180910390fd5b505b3d805f81146103a5576040519150601f19603f3d011682016040523d82523d5f602084013e6103aa565b606091505b50806040516020016103bc9190610e3c565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103fe9190610df6565b60405180910390fd5b50610583565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036104e25760055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639dc29fac84836040518363ffffffff1660e01b815260040161049c929190610c90565b6020604051808303815f875af11580156104b8573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104dc9190610bfe565b50610582565b60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663beabacc88484846040518463ffffffff1660e01b815260040161054093929190610e61565b6020604051808303815f875af115801561055c573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105809190610bfe565b505b5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105e09190610e96565b60405180910390a3505050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61064c82610606565b810181811067ffffffffffffffff8211171561066b5761066a610616565b5b80604052505050565b5f61067d6105ed565b90506106898282610643565b919050565b5f67ffffffffffffffff8211156106a8576106a7610616565b5b6106b182610606565b9050602081019050919050565b5f5b838110156106db5780820151818401526020810190506106c0565b5f8484015250505050565b5f6106f86106f38461068e565b610674565b90508281526020810184848401111561071457610713610602565b5b61071f8482856106be565b509392505050565b5f82601f83011261073b5761073a6105fe565b5b815161074b8482602086016106e6565b91505092915050565b5f60ff82169050919050565b61076981610754565b8114610773575f5ffd5b50565b5f8151905061078481610760565b92915050565b5f819050919050565b61079c8161078a565b81146107a6575f5ffd5b50565b5f815190506107b781610793565b92915050565b5f5f5f5f608085870312156107d5576107d46105f6565b5b5f85015167ffffffffffffffff8111156107f2576107f16105fa565b5b6107fe87828801610727565b945050602085015167ffffffffffffffff81111561081f5761081e6105fa565b5b61082b87828801610727565b935050604061083c87828801610776565b925050606061084d878288016107a9565b91505092959194509250565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806108a757607f821691505b6020821081036108ba576108b9610863565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261091c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826108e1565b61092686836108e1565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61096161095c6109578461078a565b61093e565b61078a565b9050919050565b5f819050919050565b61097a83610947565b61098e61098682610968565b8484546108ed565b825550505050565b5f5f905090565b6109a5610996565b6109b0818484610971565b505050565b5b818110156109d3576109c85f8261099d565b6001810190506109b6565b5050565b601f821115610a18576109e9816108c0565b6109f2846108d2565b81016020851015610a01578190505b610a15610a0d856108d2565b8301826109b5565b50505b505050565b5f82821c905092915050565b5f610a385f1984600802610a1d565b1980831691505092915050565b5f610a508383610a29565b9150826002028217905092915050565b610a6982610859565b67ffffffffffffffff811115610a8257610a81610616565b5b610a8c8254610890565b610a978282856109d7565b5f60209050601f831160018114610ac8575f8415610ab6578287015190505b610ac08582610a45565b865550610b27565b601f198416610ad6866108c0565b5f5b82811015610afd57848901518255600182019150602085019450602081019050610ad8565b86831015610b1a5784890151610b16601f891682610a29565b8355505b6001600288020188555050505b505050505050565b5f82825260208201905092915050565b5f610b4982610859565b610b538185610b2f565b9350610b638185602086016106be565b610b6c81610606565b840191505092915050565b610b8081610754565b82525050565b5f6060820190508181035f830152610b9e8186610b3f565b90508181036020830152610bb28185610b3f565b9050610bc16040830184610b77565b949350505050565b5f8115159050919050565b610bdd81610bc9565b8114610be7575f5ffd5b50565b5f81519050610bf881610bd4565b92915050565b5f60208284031215610c1357610c126105f6565b5b5f610c2084828501610bea565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610c5282610c29565b9050919050565b610c6281610c48565b82525050565b5f602082019050610c7b5f830184610c59565b92915050565b610c8a8161078a565b82525050565b5f604082019050610ca35f830185610c59565b610cb06020830184610c81565b9392505050565b5f8160e01c9050919050565b5f60033d1115610cdf5760045f5f3e610cdc5f51610cb7565b90505b90565b5f60443d10610d6e57610cf36105ed565b60043d036004823e80513d602482011167ffffffffffffffff82111715610d1b575050610d6e565b808201805167ffffffffffffffff811115610d395750505050610d6e565b80602083010160043d038501811115610d56575050505050610d6e565b610d6582602001850186610643565b82955050505050505b90565b7f6661696c656420746f206d696e743a2000000000000000000000000000000000815250565b5f81905092915050565b5f610dab82610859565b610db58185610d97565b9350610dc58185602086016106be565b80840191505092915050565b5f610ddb82610d71565b601082019150610deb8284610da1565b915081905092915050565b5f6020820190508181035f830152610e0e8184610b3f565b905092915050565b7f6661696c656420746f206d696e743a20756e6b6e6f776e206572726f723a2000815250565b5f610e4682610e16565b601f82019150610e568284610da1565b915081905092915050565b5f606082019050610e745f830186610c59565b610e816020830185610c59565b610e8e6040830184610c81565b949350505050565b5f602082019050610ea95f830184610c81565b92915050565b61154680610ebc5f395ff3fe608060405234801561000f575f5ffd5b5060043610610091575f3560e01c8063313ce56711610064578063313ce5671461013157806370a082311461014f57806395d89b411461017f578063a9059cbb1461019d578063dd62ed3e146101cd57610091565b806306fdde0314610095578063095ea7b3146100b357806318160ddd146100e357806323b872dd14610101575b5f5ffd5b61009d6101fd565b6040516100aa9190610dfe565b60405180910390f35b6100cd60048036038101906100c89190610ebc565b6102aa565b6040516100da9190610f14565b60405180910390f35b6100eb6102cc565b6040516100f89190610f3c565b60405180910390f35b61011b60048036038101906101169190610f55565b61036b565b6040516101289190610f14565b60405180910390f35b610139610399565b6040516101469190610fc0565b60405180910390f35b61016960048036038101906101649190610fd9565b610447565b6040516101769190610f3c565b60405180910390f35b6101876104ea565b6040516101949190610dfe565b60405180910390f35b6101b760048036038101906101b29190610ebc565b6105a0565b6040516101c49190610f14565b60405180910390f35b6101e760048036038101906101e29190611004565b6105c2565b6040516101f49190610f3c565b60405180910390f35b60608060055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632ba21572306040518263ffffffff1660e01b815260040161025a9190611051565b5f60405180830381865afa158015610274573d5f5f3e3d5ffd5b505050506040513d5f823e3d601f19601f8201168201806040525081019061029c91906111b2565b905050809150508091505090565b5f5f6102b4610644565b90506102c181858561064b565b600191505092915050565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e4dc2aa4306040518263ffffffff1660e01b81526004016103279190611051565b602060405180830381865afa158015610342573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610366919061124e565b905090565b5f5f610375610644565b905061038285828561065d565b61038d8585856106f0565b60019150509392505050565b5f5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632ba21572306040518263ffffffff1660e01b81526004016103f59190611051565b5f60405180830381865afa15801561040f573d5f5f3e3d5ffd5b505050506040513d5f823e3d601f19601f8201168201806040525081019061043791906111b2565b9091509050809150508091505090565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f7888aec30846040518363ffffffff1660e01b81526004016104a4929190611279565b602060405180830381865afa1580156104bf573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104e3919061124e565b9050919050565b60605b60016104ed57606060055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632ba21572306040518263ffffffff1660e01b815260040161054f9190611051565b5f60405180830381865afa158015610569573d5f5f3e3d5ffd5b505050506040513d5f823e3d601f19601f8201168201806040525081019061059191906111b2565b90915050809150508091505090565b5f5f6105aa610644565b90506105b78185856106f0565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f33905090565b61065883838360016107e0565b505050565b5f61066884846105c2565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8110156106ea57818110156106db578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016106d2939291906112a0565b60405180910390fd5b6106e984848484035f6107e0565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610760575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016107579190611051565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036107d0575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016107c79190611051565b60405180910390fd5b6107db8383836109af565b505050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610850575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016108479190611051565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036108c0575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016108b79190611051565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555080156109a9578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516109a09190610f3c565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610b945760055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f1983836040518363ffffffff1660e01b8152600401610a3e9291906112d5565b6020604051808303815f875af1925050508015610a7957506040513d601f19601f82011682018060405250810190610a769190611326565b60015b610b8e57610a8561135d565b806308c379a003610b005750610a9961137c565b80610aa45750610b02565b80604051602001610ab5919061146b565b6040516020818303038152906040526040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610af79190610dfe565b60405180910390fd5b505b3d805f8114610b2c576040519150601f19603f3d011682016040523d82523d5f602084013e610b31565b606091505b5080604051602001610b4391906114b6565b6040516020818303038152906040526040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b859190610dfe565b60405180910390fd5b50610d0a565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c695760055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639dc29fac84836040518363ffffffff1660e01b8152600401610c239291906112d5565b6020604051808303815f875af1158015610c3f573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610c639190611326565b50610d09565b60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663beabacc88484846040518463ffffffff1660e01b8152600401610cc7939291906114db565b6020604051808303815f875af1158015610ce3573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d079190611326565b505b5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610d679190610f3c565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015610dab578082015181840152602081019050610d90565b5f8484015250505050565b5f601f19601f8301169050919050565b5f610dd082610d74565b610dda8185610d7e565b9350610dea818560208601610d8e565b610df381610db6565b840191505092915050565b5f6020820190508181035f830152610e168184610dc6565b905092915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610e5882610e2f565b9050919050565b610e6881610e4e565b8114610e72575f5ffd5b50565b5f81359050610e8381610e5f565b92915050565b5f819050919050565b610e9b81610e89565b8114610ea5575f5ffd5b50565b5f81359050610eb681610e92565b92915050565b5f5f60408385031215610ed257610ed1610e27565b5b5f610edf85828601610e75565b9250506020610ef085828601610ea8565b9150509250929050565b5f8115159050919050565b610f0e81610efa565b82525050565b5f602082019050610f275f830184610f05565b92915050565b610f3681610e89565b82525050565b5f602082019050610f4f5f830184610f2d565b92915050565b5f5f5f60608486031215610f6c57610f6b610e27565b5b5f610f7986828701610e75565b9350506020610f8a86828701610e75565b9250506040610f9b86828701610ea8565b9150509250925092565b5f60ff82169050919050565b610fba81610fa5565b82525050565b5f602082019050610fd35f830184610fb1565b92915050565b5f60208284031215610fee57610fed610e27565b5b5f610ffb84828501610e75565b91505092915050565b5f5f6040838503121561101a57611019610e27565b5b5f61102785828601610e75565b925050602061103885828601610e75565b9150509250929050565b61104b81610e4e565b82525050565b5f6020820190506110645f830184611042565b92915050565b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6110a882610db6565b810181811067ffffffffffffffff821117156110c7576110c6611072565b5b80604052505050565b5f6110d9610e1e565b90506110e5828261109f565b919050565b5f67ffffffffffffffff82111561110457611103611072565b5b61110d82610db6565b9050602081019050919050565b5f61112c611127846110ea565b6110d0565b9050828152602081018484840111156111485761114761106e565b5b611153848285610d8e565b509392505050565b5f82601f83011261116f5761116e61106a565b5b815161117f84826020860161111a565b91505092915050565b61119181610fa5565b811461119b575f5ffd5b50565b5f815190506111ac81611188565b92915050565b5f5f5f606084860312156111c9576111c8610e27565b5b5f84015167ffffffffffffffff8111156111e6576111e5610e2b565b5b6111f28682870161115b565b935050602084015167ffffffffffffffff81111561121357611212610e2b565b5b61121f8682870161115b565b92505060406112308682870161119e565b9150509250925092565b5f8151905061124881610e92565b92915050565b5f6020828403121561126357611262610e27565b5b5f6112708482850161123a565b91505092915050565b5f60408201905061128c5f830185611042565b6112996020830184611042565b9392505050565b5f6060820190506112b35f830186611042565b6112c06020830185610f2d565b6112cd6040830184610f2d565b949350505050565b5f6040820190506112e85f830185611042565b6112f56020830184610f2d565b9392505050565b61130581610efa565b811461130f575f5ffd5b50565b5f81519050611320816112fc565b92915050565b5f6020828403121561133b5761133a610e27565b5b5f61134884828501611312565b91505092915050565b5f8160e01c9050919050565b5f60033d11156113795760045f5f3e6113765f51611351565b90505b90565b5f60443d106114085761138d610e1e565b60043d036004823e80513d602482011167ffffffffffffffff821117156113b5575050611408565b808201805167ffffffffffffffff8111156113d35750505050611408565b80602083010160043d0385018111156113f0575050505050611408565b6113ff8260200185018661109f565b82955050505050505b90565b7f6661696c656420746f206d696e743a2000000000000000000000000000000000815250565b5f81905092915050565b5f61144582610d74565b61144f8185611431565b935061145f818560208601610d8e565b80840191505092915050565b5f6114758261140b565b601082019150611485828461143b565b915081905092915050565b7f6661696c656420746f206d696e743a20756e6b6e6f776e206572726f723a2000815250565b5f6114c082611490565b601f820191506114d0828461143b565b915081905092915050565b5f6060820190506114ee5f830186611042565b6114fb6020830185611042565b6115086040830184610f2d565b94935050505056fea26469706673582212209e6d4497f44b6ce516feeba9e9907f7aaf4d74132be01f9488d4559f11c8946b64736f6c634300081e0033",
}

// FixedSupplyBankERC20InfiniteGasABI is the input ABI used to generate the binding from.
// Deprecated: Use FixedSupplyBankERC20InfiniteGasMetaData.ABI instead.
var FixedSupplyBankERC20InfiniteGasABI = FixedSupplyBankERC20InfiniteGasMetaData.ABI

// FixedSupplyBankERC20InfiniteGasBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FixedSupplyBankERC20InfiniteGasMetaData.Bin instead.
var FixedSupplyBankERC20InfiniteGasBin = FixedSupplyBankERC20InfiniteGasMetaData.Bin

// DeployFixedSupplyBankERC20InfiniteGas deploys a new Ethereum contract, binding an instance of FixedSupplyBankERC20InfiniteGas to it.
func DeployFixedSupplyBankERC20InfiniteGas(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, decimals_ uint8, initial_supply_ *big.Int) (common.Address, *types.Transaction, *FixedSupplyBankERC20InfiniteGas, error) {
	parsed, err := FixedSupplyBankERC20InfiniteGasMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FixedSupplyBankERC20InfiniteGasBin), backend, name_, symbol_, decimals_, initial_supply_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FixedSupplyBankERC20InfiniteGas{FixedSupplyBankERC20InfiniteGasCaller: FixedSupplyBankERC20InfiniteGasCaller{contract: contract}, FixedSupplyBankERC20InfiniteGasTransactor: FixedSupplyBankERC20InfiniteGasTransactor{contract: contract}, FixedSupplyBankERC20InfiniteGasFilterer: FixedSupplyBankERC20InfiniteGasFilterer{contract: contract}}, nil
}

// FixedSupplyBankERC20InfiniteGas is an auto generated Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGas struct {
	FixedSupplyBankERC20InfiniteGasCaller     // Read-only binding to the contract
	FixedSupplyBankERC20InfiniteGasTransactor // Write-only binding to the contract
	FixedSupplyBankERC20InfiniteGasFilterer   // Log filterer for contract events
}

// FixedSupplyBankERC20InfiniteGasCaller is an auto generated read-only Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGasCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FixedSupplyBankERC20InfiniteGasTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGasTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FixedSupplyBankERC20InfiniteGasFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FixedSupplyBankERC20InfiniteGasFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FixedSupplyBankERC20InfiniteGasSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FixedSupplyBankERC20InfiniteGasSession struct {
	Contract     *FixedSupplyBankERC20InfiniteGas // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// FixedSupplyBankERC20InfiniteGasCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FixedSupplyBankERC20InfiniteGasCallerSession struct {
	Contract *FixedSupplyBankERC20InfiniteGasCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// FixedSupplyBankERC20InfiniteGasTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FixedSupplyBankERC20InfiniteGasTransactorSession struct {
	Contract     *FixedSupplyBankERC20InfiniteGasTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// FixedSupplyBankERC20InfiniteGasRaw is an auto generated low-level Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGasRaw struct {
	Contract *FixedSupplyBankERC20InfiniteGas // Generic contract binding to access the raw methods on
}

// FixedSupplyBankERC20InfiniteGasCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGasCallerRaw struct {
	Contract *FixedSupplyBankERC20InfiniteGasCaller // Generic read-only contract binding to access the raw methods on
}

// FixedSupplyBankERC20InfiniteGasTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FixedSupplyBankERC20InfiniteGasTransactorRaw struct {
	Contract *FixedSupplyBankERC20InfiniteGasTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFixedSupplyBankERC20InfiniteGas creates a new instance of FixedSupplyBankERC20InfiniteGas, bound to a specific deployed contract.
func NewFixedSupplyBankERC20InfiniteGas(address common.Address, backend bind.ContractBackend) (*FixedSupplyBankERC20InfiniteGas, error) {
	contract, err := bindFixedSupplyBankERC20InfiniteGas(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGas{FixedSupplyBankERC20InfiniteGasCaller: FixedSupplyBankERC20InfiniteGasCaller{contract: contract}, FixedSupplyBankERC20InfiniteGasTransactor: FixedSupplyBankERC20InfiniteGasTransactor{contract: contract}, FixedSupplyBankERC20InfiniteGasFilterer: FixedSupplyBankERC20InfiniteGasFilterer{contract: contract}}, nil
}

// NewFixedSupplyBankERC20InfiniteGasCaller creates a new read-only instance of FixedSupplyBankERC20InfiniteGas, bound to a specific deployed contract.
func NewFixedSupplyBankERC20InfiniteGasCaller(address common.Address, caller bind.ContractCaller) (*FixedSupplyBankERC20InfiniteGasCaller, error) {
	contract, err := bindFixedSupplyBankERC20InfiniteGas(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasCaller{contract: contract}, nil
}

// NewFixedSupplyBankERC20InfiniteGasTransactor creates a new write-only instance of FixedSupplyBankERC20InfiniteGas, bound to a specific deployed contract.
func NewFixedSupplyBankERC20InfiniteGasTransactor(address common.Address, transactor bind.ContractTransactor) (*FixedSupplyBankERC20InfiniteGasTransactor, error) {
	contract, err := bindFixedSupplyBankERC20InfiniteGas(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasTransactor{contract: contract}, nil
}

// NewFixedSupplyBankERC20InfiniteGasFilterer creates a new log filterer instance of FixedSupplyBankERC20InfiniteGas, bound to a specific deployed contract.
func NewFixedSupplyBankERC20InfiniteGasFilterer(address common.Address, filterer bind.ContractFilterer) (*FixedSupplyBankERC20InfiniteGasFilterer, error) {
	contract, err := bindFixedSupplyBankERC20InfiniteGas(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasFilterer{contract: contract}, nil
}

// bindFixedSupplyBankERC20InfiniteGas binds a generic wrapper to an already deployed contract.
func bindFixedSupplyBankERC20InfiniteGas(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FixedSupplyBankERC20InfiniteGasMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FixedSupplyBankERC20InfiniteGas.Contract.FixedSupplyBankERC20InfiniteGasCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.FixedSupplyBankERC20InfiniteGasTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.FixedSupplyBankERC20InfiniteGasTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FixedSupplyBankERC20InfiniteGas.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Allowance(&_FixedSupplyBankERC20InfiniteGas.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Allowance(&_FixedSupplyBankERC20InfiniteGas.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.BalanceOf(&_FixedSupplyBankERC20InfiniteGas.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.BalanceOf(&_FixedSupplyBankERC20InfiniteGas.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Decimals() (uint8, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Decimals(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) Decimals() (uint8, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Decimals(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Name() (string, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Name(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) Name() (string, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Name(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Symbol() (string, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Symbol(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) Symbol() (string, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Symbol(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FixedSupplyBankERC20InfiniteGas.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) TotalSupply() (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.TotalSupply(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasCallerSession) TotalSupply() (*big.Int, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.TotalSupply(&_FixedSupplyBankERC20InfiniteGas.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Approve(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Approve(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Transfer(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.Transfer(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.TransferFrom(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _FixedSupplyBankERC20InfiniteGas.Contract.TransferFrom(&_FixedSupplyBankERC20InfiniteGas.TransactOpts, from, to, value)
}

// FixedSupplyBankERC20InfiniteGasApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasApprovalIterator struct {
	Event *FixedSupplyBankERC20InfiniteGasApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FixedSupplyBankERC20InfiniteGasApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FixedSupplyBankERC20InfiniteGasApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FixedSupplyBankERC20InfiniteGasApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FixedSupplyBankERC20InfiniteGasApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FixedSupplyBankERC20InfiniteGasApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FixedSupplyBankERC20InfiniteGasApproval represents a Approval event raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FixedSupplyBankERC20InfiniteGasApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasApprovalIterator{contract: _FixedSupplyBankERC20InfiniteGas.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FixedSupplyBankERC20InfiniteGasApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FixedSupplyBankERC20InfiniteGasApproval)
				if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) ParseApproval(log types.Log) (*FixedSupplyBankERC20InfiniteGasApproval, error) {
	event := new(FixedSupplyBankERC20InfiniteGasApproval)
	if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FixedSupplyBankERC20InfiniteGasFailureIterator is returned from FilterFailure and is used to iterate over the raw logs and unpacked data for Failure events raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasFailureIterator struct {
	Event *FixedSupplyBankERC20InfiniteGasFailure // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FixedSupplyBankERC20InfiniteGasFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FixedSupplyBankERC20InfiniteGasFailure)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FixedSupplyBankERC20InfiniteGasFailure)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FixedSupplyBankERC20InfiniteGasFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FixedSupplyBankERC20InfiniteGasFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FixedSupplyBankERC20InfiniteGasFailure represents a Failure event raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasFailure struct {
	Message string
	Data    []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFailure is a free log retrieval operation binding the contract event 0x66c9257b5635d9c11609ab746e0972276ff2412ab2085de9630ecb2300a019a6.
//
// Solidity: event Failure(string message, bytes data)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) FilterFailure(opts *bind.FilterOpts) (*FixedSupplyBankERC20InfiniteGasFailureIterator, error) {

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.FilterLogs(opts, "Failure")
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasFailureIterator{contract: _FixedSupplyBankERC20InfiniteGas.contract, event: "Failure", logs: logs, sub: sub}, nil
}

// WatchFailure is a free log subscription operation binding the contract event 0x66c9257b5635d9c11609ab746e0972276ff2412ab2085de9630ecb2300a019a6.
//
// Solidity: event Failure(string message, bytes data)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) WatchFailure(opts *bind.WatchOpts, sink chan<- *FixedSupplyBankERC20InfiniteGasFailure) (event.Subscription, error) {

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.WatchLogs(opts, "Failure")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FixedSupplyBankERC20InfiniteGasFailure)
				if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Failure", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFailure is a log parse operation binding the contract event 0x66c9257b5635d9c11609ab746e0972276ff2412ab2085de9630ecb2300a019a6.
//
// Solidity: event Failure(string message, bytes data)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) ParseFailure(log types.Log) (*FixedSupplyBankERC20InfiniteGasFailure, error) {
	event := new(FixedSupplyBankERC20InfiniteGasFailure)
	if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Failure", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FixedSupplyBankERC20InfiniteGasTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasTransferIterator struct {
	Event *FixedSupplyBankERC20InfiniteGasTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FixedSupplyBankERC20InfiniteGasTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FixedSupplyBankERC20InfiniteGasTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FixedSupplyBankERC20InfiniteGasTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FixedSupplyBankERC20InfiniteGasTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FixedSupplyBankERC20InfiniteGasTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FixedSupplyBankERC20InfiniteGasTransfer represents a Transfer event raised by the FixedSupplyBankERC20InfiniteGas contract.
type FixedSupplyBankERC20InfiniteGasTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FixedSupplyBankERC20InfiniteGasTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FixedSupplyBankERC20InfiniteGasTransferIterator{contract: _FixedSupplyBankERC20InfiniteGas.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FixedSupplyBankERC20InfiniteGasTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FixedSupplyBankERC20InfiniteGas.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FixedSupplyBankERC20InfiniteGasTransfer)
				if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FixedSupplyBankERC20InfiniteGas *FixedSupplyBankERC20InfiniteGasFilterer) ParseTransfer(log types.Log) (*FixedSupplyBankERC20InfiniteGasTransfer, error) {
	event := new(FixedSupplyBankERC20InfiniteGasTransfer)
	if err := _FixedSupplyBankERC20InfiniteGas.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
