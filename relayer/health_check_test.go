package relayer

import (
	"context"
	"testing"

	"github.com/VolumeFi/whoops"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "relayer")
}

var _ = Describe("health check", func() {
	ctx := context.Background()
	var r *Relayer
	var m *mocks.PalomaClienter
	var pm *chainmocks.Processor
	var fm *mocks.EvmFactorier

	retErr := whoops.String("oh no")

	cfg := &config.Config{}

	var val *stakingtypes.Validator
	BeforeEach(func() {
		val = &stakingtypes.Validator{}
	})

	BeforeEach(func() {
		m = mocks.NewPalomaClienter(GinkgoT())
	})

	BeforeEach(func() {
		cfg.EVM = map[string]config.EVM{
			"test": {},
		}
	})

	JustBeforeEach(func() {
		r = &Relayer{
			palomaClient: m,
			cfg:          cfg,
			evmFactory:   fm,
		}
	})

	Context("paloma returned errors", func() {
		When("query get evm chain infos returns an error", func() {
			It("returns it back", func() {
				m.On("QueryGetEVMChainInfos", mock.Anything).Return(nil, retErr)
				err := r.HealthCheck(ctx)
				Expect(err).To(MatchError(retErr))
			})
		})
		When("get validator return an error", func() {
			It("returns it back", func() {
				m.On("QueryGetEVMChainInfos", mock.Anything).Return([]*evmtypes.ChainInfo{}, nil)
				m.On("GetValidator", mock.Anything).Return(nil, retErr)
				err := r.HealthCheck(ctx)
				Expect(err).To(MatchError(retErr))
			})
		})
	})

	Context("when validator is not staking", func() {
		BeforeEach(func() {
			val.Jailed = true
		})

		BeforeEach(func() {
			pm = chainmocks.NewProcessor(GinkgoT())
			fm = mocks.NewEvmFactorier(GinkgoT())
		})

		When("validator is bonded", func() {
			BeforeEach(func() {
				val.Status = stakingtypes.Bonded
			})

			It("doesnt return any error back", func() {
				m.On("QueryGetEVMChainInfos", mock.Anything).Return([]*evmtypes.ChainInfo{
					{
						ChainReferenceID:  "test",
						MinOnChainBalance: "100000",
					},
				}, nil)
				fm.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pm, nil)
				m.On("GetValidator", mock.Anything).Return(val, nil)

				pm.On("HealthCheck", mock.Anything).Return(retErr)

				err := r.HealthCheck(ctx)
				Expect(err).To(BeNil())
			})
		})
	})
	Context("when validator is not staking", func() {
		BeforeEach(func() {
			val.Jailed = true
		})

		BeforeEach(func() {
			pm = chainmocks.NewProcessor(GinkgoT())
			fm = mocks.NewEvmFactorier(GinkgoT())
		})

		It("doesnt return any error back", func() {
			m.On("QueryGetEVMChainInfos", mock.Anything).Return([]*evmtypes.ChainInfo{
				{
					ChainReferenceID:  "test",
					MinOnChainBalance: "100000",
				},
			}, nil)
			fm.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pm, nil)
			m.On("GetValidator", mock.Anything).Return(val, nil)

			pm.On("HealthCheck", mock.Anything).Return(retErr)

			err := r.HealthCheck(ctx)
			Expect(err).To(BeNil())
		})
	})
	Context("when validator is staking", func() {
		BeforeEach(func() {
			val.Jailed = false
		})

		BeforeEach(func() {
			pm = chainmocks.NewProcessor(GinkgoT())
			fm = mocks.NewEvmFactorier(GinkgoT())
		})

		When("validator is bonded", func() {
			BeforeEach(func() {
				val.Status = stakingtypes.Bonded
			})

			It("doesnt return any error back", func() {
				m.On("QueryGetEVMChainInfos", mock.Anything).Return([]*evmtypes.ChainInfo{
					{
						ChainReferenceID:  "test",
						MinOnChainBalance: "100000",
					},
				}, nil)
				fm.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pm, nil)
				m.On("GetValidator", mock.Anything).Return(val, nil)

				pm.On("HealthCheck", mock.Anything).Return(retErr)

				err := r.HealthCheck(ctx)
				Expect(err).To(MatchError(retErr))
			})
		})
	})
})
