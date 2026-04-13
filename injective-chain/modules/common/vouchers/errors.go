package vouchers

import "cosmossdk.io/errors"

var ErrVoucherNotFound = errors.Register("vouchers", 1, "voucher not found")
