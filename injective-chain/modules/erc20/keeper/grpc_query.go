package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/erc20/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{Keeper: k}
}

func (q queryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{Params: q.GetParams(sdk.UnwrapSDKContext(c))}, nil
}

func (q queryServer) AllTokenPairs(c context.Context, req *types.QueryAllTokenPairsRequest) (*types.QueryAllTokenPairsResponse, error) {
	if req == nil {
		return nil, errors.Wrap(types.ErrInvalidQueryRequest, "no request provided")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{}
	}
	if req.Pagination.Limit == 0 {
		req.Pagination.Limit = query.DefaultLimit
	}
	if req.Pagination.Reverse {
		return nil, errors.Wrap(types.ErrInvalidQueryRequest, "pagination.reverse is unsupported for AllTokenPairs")
	}
	if len(req.Pagination.Key) > 0 {
		return nil, errors.Wrap(types.ErrInvalidQueryRequest, "pagination.key is unsupported for AllTokenPairs")
	}
	ctx := sdk.UnwrapSDKContext(c)

	tokenPairs, count, err := q.getPaginatedTokenPairs(ctx, req.Pagination)
	if err != nil {
		return nil, errors.Wrap(err, "can't get TokenPairs")
	}

	switch {
	case len(tokenPairs) >= int(req.Pagination.Limit):
		return &types.QueryAllTokenPairsResponse{
			TokenPairs: tokenPairs,
		}, nil
	case len(tokenPairs) > 0:
		req.Pagination.Offset = 0
		req.Pagination.Limit -= uint64(len(tokenPairs))
	default: // len(tokenPairs) == 0
		req.Pagination.Offset -= count
	}

	// also include erc20: denoms
	erc20Pairs, foundAddresses, err := q.getPaginatedERC20Denoms(ctx, req.Pagination)
	if err != nil {
		return nil, errors.Wrap(err, "can't iterate erc20 denoms")
	}

	pairs := append(append([]*types.TokenPair{}, tokenPairs...), erc20Pairs...)

	switch {
	case len(erc20Pairs) >= int(req.Pagination.Limit):
		return &types.QueryAllTokenPairsResponse{
			TokenPairs: pairs,
		}, nil
	case len(erc20Pairs) > 0:
		req.Pagination.Offset = 0
		req.Pagination.Limit -= uint64(len(erc20Pairs))
	default: // len(erc20Pairs) == 0
		req.Pagination.Offset -= uint64(len(foundAddresses))
	}

	// also include erc20 metadata
	erc20MetadataPairs, err := q.getPaginatedERC20DenomMetadata(ctx, req.Pagination, foundAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "can't iterate erc20 denoms with metadata")
	}

	pairs = append(pairs, erc20MetadataPairs...)

	return &types.QueryAllTokenPairsResponse{
		TokenPairs: pairs,
	}, nil
}

func (q queryServer) getPaginatedTokenPairs(ctx sdk.Context, pageReq *query.PageRequest) ([]*types.TokenPair, uint64, error) {
	store := q.getTokenPairsStoreByBankDenom(ctx)
	pairs := make([]*types.TokenPair, 0)

	pageReq.CountTotal = true

	pageRes, err := query.Paginate(store, pageReq, func(key, value []byte) error {
		pair := &types.TokenPair{
			BankDenom:    string(key),
			Erc20Address: common.BytesToAddress(value).String(),
		}
		pairs = append(pairs, pair)
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return pairs, pageRes.Total, nil
}

func (q queryServer) getPaginatedERC20Denoms(ctx sdk.Context, pageReq *query.PageRequest) ([]*types.TokenPair, map[string]struct{}, error) {
	pairs := make([]*types.TokenPair, 0)
	foundAddresses := make(map[string]struct{})
	ranger := (&collections.Range[string]{}).Prefix(types.DenomPrefix)

	err := q.bankKeeper.IterateDenoms(ctx, ranger, func(bankDenom string) bool {
		if len(pairs) >= int(pageReq.Limit) {
			return true
		}

		erc20Addr, _ := erc20AddressFromBankDenomName(bankDenom)
		foundAddresses[string(erc20Addr.Bytes())] = struct{}{}

		if uint64(len(foundAddresses))-1 < pageReq.Offset {
			return false
		}

		pair := &types.TokenPair{
			BankDenom:    bankDenom,
			Erc20Address: erc20Addr.String(),
		}
		pairs = append(pairs, pair)

		return false
	})
	if err != nil {
		return nil, nil, err
	}

	return pairs, foundAddresses, nil
}

func (q queryServer) getPaginatedERC20DenomMetadata(ctx sdk.Context, pageReq *query.PageRequest, foundAddresses map[string]struct{}) ([]*types.TokenPair, error) {
	count := uint64(0)
	pairs := make([]*types.TokenPair, 0)
	ranger := (&collections.Range[string]{}).Prefix(types.DenomPrefix)

	err := q.bankKeeper.IterateDenomsWithMetaData(ctx, ranger, func(bankDenom string) bool {
		if len(pairs) >= int(pageReq.Limit) {
			return true
		}

		erc20Addr, _ := erc20AddressFromBankDenomName(bankDenom)

		if _, exists := foundAddresses[string(erc20Addr.Bytes())]; exists {
			return false
		}

		if count < pageReq.Offset {
			count++
			return false
		}

		pair := &types.TokenPair{
			BankDenom:    bankDenom,
			Erc20Address: erc20Addr.String(),
		}
		pairs = append(pairs, pair)

		return false
	})
	if err != nil {
		return nil, err
	}

	return pairs, nil
}

func (q queryServer) TokenPairByDenom(c context.Context, req *types.QueryTokenPairByDenomRequest) (*types.QueryTokenPairByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pair, err := q.GetTokenPairForDenom(ctx, req.BankDenom)
	if err != nil {
		return nil, err
	}

	if pair != nil {
		return &types.QueryTokenPairByDenomResponse{TokenPair: pair}, nil
	}

	// no stored TokenPair, check if the bank denom is of "erc20:" type and has some supply or metadata
	erc20Addr, isERC20 := erc20AddressFromBankDenomName(req.BankDenom)
	if isERC20 {
		erc20Denom := types.DenomPrefix + erc20Addr.String()
		if q.HasBankDenomOrMetadata(ctx, erc20Denom) {
			pair := &types.TokenPair{
				BankDenom:    types.DenomPrefix + erc20Addr.String(),
				Erc20Address: erc20Addr.String(),
			}
			return &types.QueryTokenPairByDenomResponse{TokenPair: pair}, nil
		}
	}

	return &types.QueryTokenPairByDenomResponse{TokenPair: pair}, nil
}

func (q queryServer) TokenPairByERC20Address(c context.Context, req *types.QueryTokenPairByERC20AddressRequest) (*types.QueryTokenPairByERC20AddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	erc20Address := common.HexToAddress(req.Erc20Address)

	pair, err := q.GetTokenPairForERC20(ctx, erc20Address)
	if err != nil {
		return nil, err
	}

	if pair != nil {
		return &types.QueryTokenPairByERC20AddressResponse{TokenPair: pair}, nil
	}

	erc20Denom := types.DenomPrefix + erc20Address.String()
	if q.HasBankDenomOrMetadata(ctx, erc20Denom) {
		pair := &types.TokenPair{
			BankDenom:    erc20Denom,
			Erc20Address: erc20Address.String(),
		}
		return &types.QueryTokenPairByERC20AddressResponse{TokenPair: pair}, nil
	}

	return &types.QueryTokenPairByERC20AddressResponse{TokenPair: pair}, nil
}
