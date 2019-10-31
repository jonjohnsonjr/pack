package dist

import (
	"fmt"
	"sort"
	"strings"

	"github.com/buildpack/pack/style"
)

// TODO: Move to SupportsStack()
func ValidateBuildpackMixins(bp BuildpackDescriptor, builderStackID string, providedMixins []string, ignoreRunOnly bool) error {
	avail := map[string]interface{}{}
	for _, m := range providedMixins {
		avail[m] = nil
	}

	if len(bp.Stacks) == 0 {
		return nil // Order buildpack, no validation required
	}

	bpMixins, err := findBuildpackMixins(bp, builderStackID)
	if err != nil {
		return err
	}

	var missing []string
	for _, m := range bpMixins {
		ignored := ignoreRunOnly && strings.HasPrefix(m, "run:")
		if _, ok := avail[m]; !ignored && !ok {
			missing = append(missing, m)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return fmt.Errorf("buildpack %s requires missing mixin(s): %s", style.Symbol(bp.String()), strings.Join(missing, ", "))
	}
	return nil
}

func findBuildpackMixins(bp BuildpackDescriptor, stackID string) ([]string, error) {
	for _, s := range bp.Stacks {
		if s.ID == stackID {
			return s.Mixins, nil
		}
	}
	return nil, fmt.Errorf("buildpack %s does not support stack %s", style.Symbol(bp.String()), style.Symbol(stackID))
}

func ValidateStackMixins(buildImageName string, buildImageMixins []string, runImageName string, runImageMixins []string) error {
	bMixins, err := mixinSet(buildImageMixins, buildImageName, false)
	if err != nil {
		return err
	}

	rMixins, err := mixinSet(runImageMixins, runImageName, true)
	if err != nil {
		return err
	}

	if err := validateMissing(rMixins, bMixins, runImageName); err != nil {
		return err
	}
	return nil
}

func mixinSet(mixins []string, imageName string, run bool) (map[string]interface{}, error) {
	set := map[string]interface{}{}
	invalidPrefix := "run"
	filterPrefix := "build"
	if run {
		invalidPrefix = "build"
		filterPrefix = "run"
	}

	var invalid []string
	for _, m := range mixins {
		if strings.HasPrefix(m, invalidPrefix+":") {
			invalid = append(invalid, m)
			continue
		}
		if strings.HasPrefix(m, filterPrefix+":") {
			continue
		}
		set[m] = nil
	}

	if len(invalid) > 0 {
		sort.Strings(invalid)
		return nil, fmt.Errorf("%s contains %s-only mixin(s): %s", style.Symbol(imageName), invalidPrefix, strings.Join(invalid, ", "))
	}
	return set, nil
}

func validateMissing(actual, required map[string]interface{}, actualImageName string) error {
	var missing []string
	for m := range required {
		if _, ok := actual[m]; !ok {
			missing = append(missing, m)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("%s missing required mixin(s): %s", style.Symbol(actualImageName), strings.Join(missing, ", "))
	}
	return nil
}
