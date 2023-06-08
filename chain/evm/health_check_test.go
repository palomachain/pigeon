package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/palomachain/pigeon/chain"
	"github.com/stretchr/testify/mock"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Health checks")
}

var _ = Describe("health check", func() {
	var ethClientConn *mockEthClientConn
	var client *Client
	var p Processor

	minOnChainBalance := big.NewInt(50000)
	ctx := context.Background()

	var addr common.Address

	BeforeEach(func() {
		ethClientConn = newMockEthClientConn(GinkgoT())
	})

	JustBeforeEach(func() {
		client = &Client{
			conn: ethClientConn,
			addr: addr,
		}

		p = Processor{
			evmClient:         client,
			minOnChainBalance: minOnChainBalance,
		}
	})

	Context("missing addr", func() {
		It("returns an error that the account is missing", func() {
			err := p.HealthCheck(ctx)
			Expect(err).To(MatchError(chain.ErrMissingAccount))
		})
	})

	Context("valid addr", func() {
		BeforeEach(func() {
			addr = common.HexToAddress("0x1234")
		})

		Context("balance returns an error", func() {
			It("returns the same error back", func() {
				retErr := whoops.String("oh no")
				ethClientConn.On("BalanceAt", mock.Anything, addr, (*big.Int)(nil)).Return(nil, retErr)
				err := p.HealthCheck(ctx)
				Expect(err).To(MatchError(retErr))
			})
		})

		Context("balance is set", func() {
			When("balance is lower", func() {
				It("returns the error", func() {
					ethClientConn.On("BalanceAt", mock.Anything, addr, (*big.Int)(nil)).Return(big.NewInt(1), nil)
					err := p.HealthCheck(ctx)
					Expect(err).To(MatchError(chain.ErrAccountBalanceLow))
				})
			})
			When("balance is zero", func() {
				It("returns the error", func() {
					ethClientConn.On("BalanceAt", mock.Anything, addr, (*big.Int)(nil)).Return(big.NewInt(0), nil)
					err := p.HealthCheck(ctx)
					Expect(err).To(MatchError(chain.ErrAccountBalanceLow))
				})
			})
			When("balance is greater than the min", func() {
				It("returns no error", func() {
					ethClientConn.On("BalanceAt", mock.Anything, addr, (*big.Int)(nil)).Return(big.NewInt(9999999), nil)
					err := p.HealthCheck(ctx)
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
