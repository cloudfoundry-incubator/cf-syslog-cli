package command_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cf-drain-cli/internal/command"
)

var _ = Describe("PushSpaceDrain", func() {
	var (
		logger *stubLogger
		cli    *stubCliConnection
	)

	BeforeEach(func() {
		logger = &stubLogger{}
		cli = newStubCliConnection()
		cli.currentSpaceGuid = "space-guid"
		cli.apiEndpoint = "https://api.something.com"
	})

	It("pushes app from the given space-drain directory", func() {
		command.PushSpaceDrain(
			cli,
			[]string{
				"--path", "some-temp-dir",
				"--drain-name", "some-drain",
				"--drain-url", "https://some-drain",
				"--username", "some-user",
				"--password", "some-password",
				"--skip-ssl-validation", "true",
			},
			logger,
		)

		Expect(cli.cliCommandArgs).To(HaveLen(2))
		Expect(cli.cliCommandArgs[0]).To(Equal(
			[]string{
				"push", "space-drain",
				"-p", "some-temp-dir",
				"-b", "binary_buildpack",
				"-c", "./space_manager",
				"--health-check-type", "process",
				"--no-start",
			},
		))

		Expect(cli.cliCommandWithoutTerminalOutputArgs).To(ConsistOf(
			[]string{"set-env", "space-drain", "SPACE_ID", "space-guid"},
			[]string{"set-env", "space-drain", "DRAIN_NAME", "some-drain"},
			[]string{"set-env", "space-drain", "DRAIN_URL", "https://some-drain"},
			[]string{"set-env", "space-drain", "API_ADDR", "https://api.something.com"},
			[]string{"set-env", "space-drain", "UAA_ADDR", "https://uaa.something.com"},
			[]string{"set-env", "space-drain", "CLIENT_ID", "cf"},
			[]string{"set-env", "space-drain", "USERNAME", "some-user"},
			[]string{"set-env", "space-drain", "PASSWORD", "some-password"},
			[]string{"set-env", "space-drain", "SKIP_CERT_VERIFY", "true"},
		))

		Expect(cli.cliCommandArgs[1]).To(Equal(
			[]string{
				"start", "space-drain",
			},
		))
	})

	DescribeTable("fatally logs if setting env variables fails", func(env string) {
		cli.setEnvErrors[env] = errors.New("some-error")

		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
					"--skip-ssl-validation", "true",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("some-error"))
	},
		Entry("SPACE_ID", "SPACE_ID"),
		Entry("DRAIN_NAME", "DRAIN_NAME"),
		Entry("DRAIN_URL", "DRAIN_URL"),
		Entry("API_ADDR", "API_ADDR"),
		Entry("UAA_ADDR", "UAA_ADDR"),
		Entry("CLIENT_ID", "CLIENT_ID"),
		Entry("USERNAME", "USERNAME"),
		Entry("PASSWORD", "PASSWORD"),
		Entry("SKIP_CERT_VERIFY", "SKIP_CERT_VERIFY"),
	)

	It("fatally logs if fetching the space fails", func() {
		cli.currentSpaceError = errors.New("some-error")
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
					"--skip-ssl-validation", "true",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("some-error"))
	})

	It("fatally logs if fetching the api endpoint fails", func() {
		cli.apiEndpointErr = errors.New("some-error")
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
					"--skip-ssl-validation", "true",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("some-error"))
	})

	It("fatally logs if the push fails", func() {
		cli.pushAppError = errors.New("failed to push")
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
					"--skip-ssl-validation", "true",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("failed to push"))
	})

	It("fatally logs if the space-drain path is not provided", func() {
		Expect(func() {
			command.PushSpaceDrain(cli,
				[]string{
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("required flag --path missing"))
	})

	It("fatally logs if the space-drain drain-name is not provided", func() {
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
					"--password", "some-password",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("required flag --drain-name missing"))
	})

	It("fatally logs if the space-drain drain-url is not provided", func() {
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--username", "some-user",
					"--password", "some-password",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("required flag --drain-url missing"))
	})

	It("fatally logs if the space-drain username is not provided", func() {
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--password", "some-password",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("required flag --username missing"))
	})

	It("fatally logs if the space-drain password is not provided", func() {
		Expect(func() {
			command.PushSpaceDrain(
				cli,
				[]string{
					"--path", "some-temp-dir",
					"--drain-name", "some-drain",
					"--drain-url", "https://some-drain",
					"--username", "some-user",
				},
				logger,
			)
		}).To(Panic())
		Expect(logger.fatalfMessage).To(Equal("required flag --password missing"))
	})
})