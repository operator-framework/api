package types

import (
	"sort"
	"strings"
)

// AnnotationsFile holds annotation information about a bundle
type AnnotationsFile struct {
	// annotations is a list of annotations for a given bundle
	Annotations Annotations `json:"annotations" yaml:"annotations"`
}

// Annotations is a list of annotations for a given bundle
type Annotations struct {
	// PackageName is the name of the overall package, ala `etcd`.
	PackageName string `json:"operators.operatorframework.io.bundle.package.v1" yaml:"operators.operatorframework.io.bundle.package.v1"`

	// Channels are a comma separated list of the declared channels for the bundle, ala `stable` or `alpha`.
	Channels string `json:"operators.operatorframework.io.bundle.channels.v1" yaml:"operators.operatorframework.io.bundle.channels.v1"`

	// DefaultChannelName is, if specified, the name of the default channel for the package. The
	// default channel will be installed if no other channel is explicitly given. If the package
	// has a single channel, then that channel is implicitly the default.
	DefaultChannelName string `json:"operators.operatorframework.io.bundle.channel.default.v1" yaml:"operators.operatorframework.io.bundle.channel.default.v1"`
}


// GetName returns the package name of the bundle
func (a *AnnotationsFile) GetName() string {
	return a.Annotations.PackageName
}

// GetChannels returns the channels that this bundle should be added to
func (a *AnnotationsFile) GetChannels() []string {
	if a.Annotations.Channels != "" {
		return strings.Split(a.Annotations.Channels, ",")
	}
	return []string{}
}

// GetDefaultChannelName returns the name of the default channel
func (a *AnnotationsFile) GetDefaultChannelName() string {
	return a.Annotations.DefaultChannelName
}

// SelectDefaultChannel returns the first item in channel list that is sorted
// in lexicographic order.
func (a *AnnotationsFile) SelectDefaultChannel() string {
	if a.Annotations.Channels != "" {
		channels := strings.Split(a.Annotations.Channels, ",")
		sort.Strings(channels)
		return channels[0]
	}

	return ""
}
