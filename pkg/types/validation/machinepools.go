package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift-metalkube/kni-installer/pkg/types"
	"github.com/openshift-metalkube/kni-installer/pkg/types/aws"
	awsvalidation "github.com/openshift-metalkube/kni-installer/pkg/types/aws/validation"
	"github.com/openshift-metalkube/kni-installer/pkg/types/azure"
	azurevalidation "github.com/openshift-metalkube/kni-installer/pkg/types/azure/validation"
	"github.com/openshift-metalkube/kni-installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift-metalkube/kni-installer/pkg/types/baremetal/validation"
	"github.com/openshift-metalkube/kni-installer/pkg/types/libvirt"
	libvirtvalidation "github.com/openshift-metalkube/kni-installer/pkg/types/libvirt/validation"
	"github.com/openshift-metalkube/kni-installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift-metalkube/kni-installer/pkg/types/openstack/validation"
)

var (
	validHyperthreadingModes = map[types.HyperthreadingMode]bool{
		types.HyperthreadingDisabled: true,
		types.HyperthreadingEnabled:  true,
	}

	validHyperthreadingModeValues = func() []string {
		v := make([]string, 0, len(validHyperthreadingModes))
		for m := range validHyperthreadingModes {
			v = append(v, string(m))
		}
		return v
	}()
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *types.Platform, p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Replicas != nil {
		if *p.Replicas < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), p.Replicas, "number of replicas must not be negative"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("replicas"), "replicas is required"))
	}
	if !validHyperthreadingModes[p.Hyperthreading] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("hyperthreading"), p.Hyperthreading, validHyperthreadingModeValues))
	}
	allErrs = append(allErrs, validateMachinePoolPlatform(platform, &p.Platform, fldPath.Child("platform"))...)
	return allErrs
}

func validateMachinePoolPlatform(platform *types.Platform, p *types.MachinePoolPlatform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	platformName := platform.Name()
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		f := fldPath.Child(n)
		if platformName == n {
			allErrs = append(allErrs, validation(f)...)
		} else {
			allErrs = append(allErrs, field.Invalid(f, value, fmt.Sprintf("cannot specify %q for machine pool when cluster is using %q", n, platformName)))
		}
	}
	if p.AWS != nil {
		validate(aws.Name, p.AWS, func(f *field.Path) field.ErrorList { return awsvalidation.ValidateMachinePool(platform.AWS, p.AWS, f) })
	}
	if p.Azure != nil {
		validate(azure.Name, p.Azure, func(f *field.Path) field.ErrorList { return azurevalidation.ValidateMachinePool(p.Azure, f) })
	}
	if p.Libvirt != nil {
		validate(libvirt.Name, p.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidateMachinePool(p.Libvirt, f) })
	}
	if p.OpenStack != nil {
		validate(openstack.Name, p.OpenStack, func(f *field.Path) field.ErrorList { return openstackvalidation.ValidateMachinePool(p.OpenStack, f) })
	}
	if p.BareMetal != nil {
		validate(baremetal.Name, p.BareMetal, func(f *field.Path) field.ErrorList { return baremetalvalidation.ValidateMachinePool(p.BareMetal, f) })
	}
	return allErrs
}
