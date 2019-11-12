// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RqPodAffinity) DeepCopyInto(out *RqPodAffinity) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RqPodAffinity.
func (in *RqPodAffinity) DeepCopy() *RqPodAffinity {
	if in == nil {
		return nil
	}
	out := new(RqPodAffinity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rqcluster) DeepCopyInto(out *Rqcluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rqcluster.
func (in *Rqcluster) DeepCopy() *Rqcluster {
	if in == nil {
		return nil
	}
	out := new(Rqcluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Rqcluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RqclusterList) DeepCopyInto(out *RqclusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Rqcluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RqclusterList.
func (in *RqclusterList) DeepCopy() *RqclusterList {
	if in == nil {
		return nil
	}
	out := new(RqclusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RqclusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RqclusterSpec) DeepCopyInto(out *RqclusterSpec) {
	*out = *in
	in.PodAffinity.DeepCopyInto(&out.PodAffinity)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RqclusterSpec.
func (in *RqclusterSpec) DeepCopy() *RqclusterSpec {
	if in == nil {
		return nil
	}
	out := new(RqclusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RqclusterStatus) DeepCopyInto(out *RqclusterStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RqclusterStatus.
func (in *RqclusterStatus) DeepCopy() *RqclusterStatus {
	if in == nil {
		return nil
	}
	out := new(RqclusterStatus)
	in.DeepCopyInto(out)
	return out
}
