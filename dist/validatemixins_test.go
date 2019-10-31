package dist

import (
	"testing"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	h "github.com/buildpack/pack/testhelpers"
)

func TestMixinValidation(t *testing.T) {
	color.Disable(true)
	defer func() { color.Disable(false) }()
	spec.Run(t, "testMixinValidation", testMixinValidation, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testMixinValidation(t *testing.T, when spec.G, it spec.S) {
	when("#ValidateBuildpackMixins", func() {
		when("not validating against run image mixins", func() {
			it("ignores run-only mixins", func() {
				bp := BuildpackDescriptor{
					Info: BuildpackInfo{
						ID:      "some.buildpack.id",
						Version: "some.buildpack.version",
					},
					Stacks: []Stack{{
						ID:     "some.stack.id",
						Mixins: []string{"mixinA", "build:mixinB", "run:mixinD"},
					}},
				}

				providedMixins := []string{"mixinA", "build:mixinB", "mixinC"}
				h.AssertNil(t, ValidateBuildpackMixins(bp, "some.stack.id", providedMixins, true))
			})

			it("returns an error with any missing (and non-ignored) mixins", func() {
				bp := BuildpackDescriptor{
					Info: BuildpackInfo{
						ID:      "some.buildpack.id",
						Version: "some.buildpack.version",
					},
					Stacks: []Stack{{
						ID:     "some.stack.id",
						Mixins: []string{"mixinX", "mixinY", "run:mixinZ"},
					}},
				}

				providedMixins := []string{"mixinA", "mixinB"}
				err := ValidateBuildpackMixins(bp, "some.stack.id", providedMixins, true)

				h.AssertError(t, err, "buildpack 'some.buildpack.id@some.buildpack.version' requires missing mixin(s): mixinX, mixinY")
			})
		})

		when("validating against run image mixins", func() {
			it("requires run-only mixins", func() {
				bp := BuildpackDescriptor{
					Info: BuildpackInfo{
						ID:      "some.buildpack.id",
						Version: "some.buildpack.version",
					},
					Stacks: []Stack{{
						ID:     "some.stack.id",
						Mixins: []string{"mixinA", "build:mixinB", "run:mixinD"},
					}},
				}

				providedMixins := []string{"mixinA", "build:mixinB", "mixinC", "run:mixinD"}

				h.AssertNil(t, ValidateBuildpackMixins(bp, "some.stack.id", providedMixins, false))
			})

			it("returns an error with any missing mixins", func() {
				bp := BuildpackDescriptor{
					Info: BuildpackInfo{
						ID:      "some.buildpack.id",
						Version: "some.buildpack.version",
					},
					Stacks: []Stack{{
						ID:     "some.stack.id",
						Mixins: []string{"mixinX", "mixinY", "run:mixinZ"},
					}},
				}

				providedMixins := []string{"mixinA", "mixinB"}

				err := ValidateBuildpackMixins(bp, "some.stack.id", providedMixins, false)

				h.AssertError(t, err, "buildpack 'some.buildpack.id@some.buildpack.version' requires missing mixin(s): mixinX, mixinY, run:mixinZ")
			})
		})

		it("returns an error when buildpack does not support stack", func() {
			bp := BuildpackDescriptor{
				Info: BuildpackInfo{
					ID:      "some.buildpack.id",
					Version: "some.buildpack.version",
				},
				Stacks: []Stack{{
					ID:     "some.stack.id",
					Mixins: []string{"mixinX", "mixinY"},
				}},
			}

			err := ValidateBuildpackMixins(bp, "some.nonexistent.stack.id", []string{"mixinA"}, false)

			h.AssertError(t, err, "buildpack 'some.buildpack.id@some.buildpack.version' does not support stack 'some.nonexistent.stack.id")
		})

		it("skips validating order buildpack", func() {
			bp := BuildpackDescriptor{
				Info: BuildpackInfo{
					ID:      "some.buildpack.id",
					Version: "some.buildpack.version",
				},
				Stacks: []Stack{},
			}

			h.AssertNil(t, ValidateBuildpackMixins(bp, "some.stack.id", []string{"mixinA"}, false))
		})
	})

	when("#ValidateStackMixins", func() {
		it("ignores stage-specific mixins", func() {
			buildMixins := []string{"mixinA", "build:mixinB"}
			runMixins := []string{"mixinA", "run:mixinC"}

			h.AssertNil(t, ValidateStackMixins("some/build", buildMixins, "some/run", runMixins))
		})

		it("allows extraneous run image mixins", func() {
			buildMixins := []string{"mixinA"}
			runMixins := []string{"mixinA", "mixinB"}

			h.AssertNil(t, ValidateStackMixins("some/build", buildMixins, "some/run", runMixins))
		})

		it("returns an error with any missing run image mixins", func() {
			buildMixins := []string{"mixinA", "mixinB"}
			runMixins := []string{}

			err := ValidateStackMixins("some/build", buildMixins, "some/run", runMixins)

			h.AssertError(t, err, "'some/run' missing required mixin(s): mixinA, mixinB")
		})

		it("returns an error with any invalid build image mixins", func() {
			buildMixins := []string{"run:mixinA", "run:mixinB"}
			runMixins := []string{}

			err := ValidateStackMixins("some/build", buildMixins, "some/run", runMixins)

			h.AssertError(t, err, "'some/build' contains run-only mixin(s): run:mixinA, run:mixinB")
		})

		it("returns an error with any invalid run image mixins", func() {
			buildMixins := []string{}
			runMixins := []string{"build:mixinA", "build:mixinB"}

			err := ValidateStackMixins("some/build", buildMixins, "some/run", runMixins)

			h.AssertError(t, err, "'some/run' contains build-only mixin(s): build:mixinA, build:mixinB")
		})
	})
}
