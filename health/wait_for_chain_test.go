package health

import (
	"context"
	"testing"
	"time"

	"github.com/VolumeFi/whoops"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/palomachain/pigeon/health/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Health checks")
}

var _ = Describe("waiting for paloma", func() {
	var m *mocks.PalomaStatuser
	var ctx context.Context

	BeforeEach(func() {
		m = mocks.NewPalomaStatuser(GinkgoT())
		ctx = context.Background()
	})

	When("paloma is online", func() {
		It("returns nil", func() {
			m.On("PalomaStatus", mock.Anything).Return(nil)
			err := WaitForPaloma(ctx, m)

			Expect(err).To(BeNil())
		})
	})

	When("paloma is not online", func() {
		Context("context returns done", func() {
			It("returns context error", func() {
				ctx, cancel := context.WithCancel(ctx)
				m.On("PalomaStatus", mock.Anything).Return(whoops.String("doesnt matter")).Run(func(_ mock.Arguments) {
					cancel()
				})
				err := WaitForPaloma(ctx, m)

				Expect(err).To(MatchError(context.Canceled))
			})
		})
		Context("paloma is online after second loop", func() {
			It("returns context error", func() {
				m.On("PalomaStatus", mock.Anything).Return(whoops.String("doesnt matter")).Times(1)
				m.On("PalomaStatus", mock.Anything).Return(nil).Times(1)
				err := WaitForPaloma(ctx, m)

				Expect(err).To(BeNil())
			})
		})
	})
})

var _ = Describe("canceling context if paloma is down", func() {
	var m *mocks.PalomaStatuser
	var ctx context.Context
	var cancelCtx func()

	BeforeEach(func() {
		m = mocks.NewPalomaStatuser(GinkgoT())
		ctx = context.Background()
		ctx, cancelCtx = context.WithCancel(ctx)
	})

	AfterEach(func() {
		cancelCtx()
	})

	Context("paloma is online", func() {
		It("returns nil", func() {
			m.On("PalomaStatus", mock.Anything).Return(nil).Times(1)
			ctx := CancelContextIfPalomaIsDown(ctx, m)
			time.Sleep(1 * time.Second)
			Expect(ctx.Err()).To(BeNil())
		})

		When("paloma goes offline", func() {
			It("returns ctx error", func() {
				m.On("PalomaStatus", mock.Anything).Return(nil).Times(1)
				m.On("PalomaStatus", mock.Anything).Return(whoops.String("oh no")).Times(1)
				ctx := CancelContextIfPalomaIsDown(ctx, m)
				Expect(ctx.Err()).To(BeNil())
				time.Sleep(6 * time.Second)
				Expect(ctx.Err()).To(MatchError(context.Canceled))
			})
		})
	})
})
