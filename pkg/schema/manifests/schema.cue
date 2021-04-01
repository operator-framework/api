package schema

#icon: {
    base64data: !=""
    mediatype: !=""
}

// Generic gvk struct
#gvk: {
    group: !=""
    version: !=""
    kind: !=""
}

#olmgvk: #property & {
  type: "olm.gvk"
  value: #gvk
}

#olmgvkrequired: #property & {
  type: "olm.gvk.required"
  value: #gvk
}

// Generic package struct
#package: {
  packageName: !=""
  version: !=""
}

#olmpackageproperty: #property & {
  type: "olm.package"
  value: #package
}

#olmpackagerequired: #property & {
  type: "olm.package.required"
  value: #package
}

// Generic channel struct
#channel: {
  name: !=""
  replaces?: !=""
}

#olmchannel: #property & {
  type: "olm.channel"
  value: #channel
}

#olmskips: #property & {
  type: "olm.skips"
  value: !=""
}

#olmskipRange: #property & {
  type: "olm.skipRange"
  value: !=""
}

#bundleref: {
  ref: !=""
}

#olmbundleobject: #property & {
  type: "olm.bundle.object"
  value: #bundleref
}

// Generic property struct
#property: {
  type: !=""
  ...
}

#relatedImages: {
  name: !=""
  image: !=""
}

// schema: "olm.package"
#olmpackage: #item & {
  schema: "olm.package"
  name: !=""
  defaultChannel: !=""
  icon: #icon
  description: !=""
  ...
}

// schema: "olm.bundle"
#olmbundle: #item & {
  schema: "olm.bundle"
  name: !=""
  image: !=""
  relatedImages: [...#relatedImages]
  ...
}

#item: {
  schema: !=""
  properties: [...#property]
  ...
}
