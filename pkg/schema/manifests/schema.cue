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

#olmgvkprovided: #property & {
  type: "olm.gvk.provided"
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

#packageproperty: #property & {
  type: "olm.package"
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
