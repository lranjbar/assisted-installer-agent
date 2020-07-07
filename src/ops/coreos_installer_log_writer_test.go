package ops

import (
	"io/ioutil"
	"testing"

	"github.com/eranco74/assisted-installer/src/inventory_client"
	"github.com/golang/mock/gomock"

	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestCoreosInstallerLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "installer_test")
}

var _ = Describe("Verify CoreosInstallerLogger", func() {
	var (
		l = logrus.New()
	)
	l.SetOutput(ioutil.Discard)
	Context("test coreosInstlalerLogger", func() {
		var (
			cilogger     *CoreosInstallerLogWriter
			hook         *test.Hook
			logger       *logrus.Logger
			mockbmclient *inventory_client.MockInventoryClient
		)
		udpateStatusSuccess := func(statuses []string) {
			for _, status := range statuses {
				mockbmclient.EXPECT().UpdateHostStatus(status, "hostID").Return(nil).Times(1)

			}
		}

		BeforeEach(func() {
			logger, hook = test.NewNullLogger()
			ctrl := gomock.NewController(GinkgoT())
			mockbmclient = inventory_client.NewMockInventoryClient(ctrl)
			cilogger = NewCoreosInstallerLogWriter(logger, mockbmclient, "hostID")
		})
		It("test log with new line", func() {
			_, err := cilogger.Write([]byte("some log with a new line \n"))
			Expect(err).Should(BeNil())
			Expect(len(hook.Entries)).Should(Equal(1))

		})
		It("test full progress line", func() {
			udpateStatusSuccess([]string{"Writing image to disk - 56%"})
			_, err := cilogger.Write([]byte("> Read disk 473.8 MiB/844.7 MiB (56%)   \r"))
			Expect(err).Should(BeNil())
			Expect(len(hook.Entries)).Should(Equal(1))
		})
		It("test partial line", func() {
			_, err := cilogger.Write([]byte("844.7 MiB"))
			Expect(err).Should(BeNil())
			Expect(len(hook.Entries)).Should(Equal(0))

		})
		It("test partial line - should log", func() {
			udpateStatusSuccess([]string{"Writing image to disk - 58%"})
			testLogs := []string{"> Read ",
				"disk",
				" ",
				"473.6 MiB",
				"/",
				"844.7 MiB",
				" (",
				"58",
				"%)   \r",
			}
			for i := range testLogs {
				_, err := cilogger.Write([]byte(testLogs[i]))
				Expect(err).Should(BeNil())
			}
			Expect(len(hook.Entries)).Should(Equal(1))
		})
		It("test multiple lines", func() {
			udpateStatusSuccess([]string{"Writing image to disk - 55%",
				"Writing image to disk - 60%",
				"Writing image to disk - 66%",
			})
			testLogs := []string{"> Read disk 471.2 MiB/844.7 MiB (55%)   \r",
				"> Read ",
				"disk",
				" ",
				"471.6 MiB",
				"/",
				"844.7 MiB",
				" (",
				"55",
				"%)   \r",
				"> Read ",
				"disk",
				" ",
				"472.1 MiB",
				"/",
				"844.7 MiB",
				" (",
				"55",
				"%)   \r",
				"> Read disk 472.6 MiB/844.7 MiB (55%)   \r",
				"> Read disk 472.8 MiB/844.7 MiB (55%)   \r",
				"> Read disk 472.9 MiB/844.7 MiB (55%)   \r",
				"> Read disk 473.0 MiB/844.7 MiB (60%)   \r",
				"> Read disk 473.3 MiB/844.7 MiB (62%)   \r",
				"> Read ",
				"disk",
				" ",
				"473.6 MiB",
				"/",
				"844.7 MiB",
				" (",
				"56",
				"%)   \r",
				"> Read disk 473.8 MiB/844.7 MiB (66%)   \r"}
			for i := range testLogs {
				_, err := cilogger.Write([]byte(testLogs[i]))
				Expect(err).Should(BeNil())
			}
			Expect(len(hook.Entries)).Should(Equal(10))

		})
	})
})