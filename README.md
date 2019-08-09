# api
Contains the API definitions used by [OLM](olm) and [Marketplace](marketplace).

## Marketplace APIs

* `CatalogSourceConfig`: CatalogSourceConfigs are used to enable an operator present in the OperatorSource to your cluster. Behind the scenes, it will configure an OLM CatalogSource so that the operator can then be managed by OLM.
* `OperatorSource`: OperatorSources are used to define the external datastore we are using to store operator bundles.

### Generating deepcopy functions
The generate deepcopy functions can be updated after changing Marketpalce APIs by running the following command with version v0.10.0 of the [Operator-SDK](operator-sdk):
```bash
$ operator-sdk generate k8s
```

[operator-sdk]: https://github.com/operator-framework/operator-sdk/

[marketplace]: https://github.com/operator-framework/operator-marketplace/

[olm]: https://github.com/operator-framework/operator-lifecycle-manager/