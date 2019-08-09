// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package shared

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectPhase) DeepCopyInto(out *ObjectPhase) {
	*out = *in
	out.Phase = in.Phase
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectPhase.
func (in *ObjectPhase) DeepCopy() *ObjectPhase {
	if in == nil {
		return nil
	}
	out := new(ObjectPhase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Phase) DeepCopyInto(out *Phase) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Phase.
func (in *Phase) DeepCopy() *Phase {
	if in == nil {
		return nil
	}
	out := new(Phase)
	in.DeepCopyInto(out)
	return out
}
