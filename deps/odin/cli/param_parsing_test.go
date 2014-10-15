package cli_test

import (
	. "github.com/Sam-Izdat/pogo/deps/odin/cli"
	"github.com/Sam-Izdat/pogo/deps/odin/cli/values"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Param Parsing", func() {

	var cli *CLI
	var cmd Command
	var didRun bool

	BeforeEach(func() {
		didRun = false
		runFn := func(c Command) {
			cmd = c
			didRun = true
		}
		cli = New("v1.0.0", "sample description", runFn)
		cli.ErrorHandling = PanicOnError
		cli.Mute()

		cli.DefineParams("paramA", "paramB")
	})

	Context("missing params", func() {
		It("should panic on a single missing param", func() {
			Expect(func() { cli.Start("cmd", "a") }).Should(Panic())
		})

		It("should panic on a multiple missing params", func() {
			Expect(func() { cli.Start("cmd") }).Should(Panic())
		})
	})

	It("should set the parameters by position", func() {
		cli.Start("cmd", "foo", "bar")
		Expect(cmd.Param("paramA").Get()).To(Equal("foo"))
		Expect(cmd.Param("paramB").Get()).To(Equal("bar"))
		Expect(cmd.Params()).To(
			Equal(
				values.Map{"paramA": cmd.Param("paramA"), "paramB": cmd.Param("paramB")},
			),
		)
	})

	Context("when a paramter is mising", func() {
		It("should raise an error", func() {
			Expect(func() { cli.Start("cmd") }).Should(Panic())
		})
	})

	It("Should be a value map", func() {
		cli.Start("cmd", "foo", "bar")
		Expect(cmd.Params().Keys()).To(ContainElement("paramA"))
		Expect(cmd.Params().Keys()).To(ContainElement("paramB"))
		Expect(cmd.Params().Values().GetAll()).To(ContainElement("foo"))
		Expect(cmd.Params().Values().GetAll()).To(ContainElement("bar"))
		Expect(cmd.Params().Values().Strings()).To(ContainElement("foo"))
		Expect(cmd.Params().Values().Strings()).To(ContainElement("bar"))
	})

})
