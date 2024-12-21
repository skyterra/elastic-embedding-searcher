package emb

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/skyterra/elastic-embedding-searcher/pb"
	"time"
)

var _ = Describe("modelx runner", func() {
	Context("ModelX", func() {
		It("start ModelX", func() {
			err := StartModelX(3, "../../modelx/service.py", "paraphrase-multilingual-MiniLM-L12-v2")
			Expect(err).Should(Succeed())
		})

		It("ping", func() {
			reply, err := GetClient().Ping(context.Background(), &pb.PingRequest{})
			Expect(err).Should(Succeed())
			Expect(reply.Code == 0).Should(BeTrue())
		})

		It("kill modelx", func() {
			err := ins.kill()
			Expect(err).Should(Succeed())
		})

		It("ping", func() {
			var err error
			for i := 0; i < 10; i++ {
				if err = ins.ping(ins.client); err == nil {
					return
				}
				time.Sleep(time.Second)
			}

			Expect(err).Should(Succeed())
		})

		It("stop modelX", func() {
			err := StopModelX()
			Expect(err).Should(Succeed())
		})
	})
})
