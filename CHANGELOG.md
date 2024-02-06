# Changelog

## [4.0.0](https://github.com/joshuar/go-hass-anything/compare/v3.2.0...v4.0.0) (2024-02-06)


### ⚠ BREAKING CHANGES

* config -> preferences rewrite

### Bug Fixes

* **examples:** :bug: fix exampleapp New func ([fd8200a](https://github.com/joshuar/go-hass-anything/commit/fd8200a400879e439ac61765a884d584a6da9072))
* **tools:** :bug: fix import path in generator template ([cab23b2](https://github.com/joshuar/go-hass-anything/commit/cab23b28785425fd3dd83c99608099c5da1e0f82))


### Code Refactoring

* config -&gt; preferences rewrite ([9816c02](https://github.com/joshuar/go-hass-anything/commit/9816c026d472a27fd16b7401db694fc8b527c400))

## [3.2.0](https://github.com/joshuar/go-hass-anything/compare/v3.1.0...v3.2.0) (2024-02-01)


### Features

* **config:** :fire: remove deprecated and unused viperconfig wrapper ([889563e](https://github.com/joshuar/go-hass-anything/commit/889563e3ec511f522bd78a79f7fcacb17060a039))
* **config:** :sparkles: assume a default path/file for config, provide methods for overriding ([13d190c](https://github.com/joshuar/go-hass-anything/commit/13d190c61c44c3477a5985c32386974c9e1ea324))

## [3.1.0](https://github.com/joshuar/go-hass-anything/compare/v3.0.1...v3.1.0) (2024-01-26)


### Features

* **config:** registering an app is now independent of any app config; a path to the registry can be specified for customisation ([d3f16db](https://github.com/joshuar/go-hass-anything/commit/d3f16dbd583c7c422333be8d6b5d9cc90b0d479b))

## [3.0.1](https://github.com/joshuar/go-hass-anything/compare/v3.0.0...v3.0.1) (2024-01-26)


### Bug Fixes

* **cmd:** :bug: fix missing parameter to mqtt.NewMQTTClient ([aa247a4](https://github.com/joshuar/go-hass-anything/commit/aa247a49c177ad49be38181da191afe645779bc1))
* **config:** :bug: make sure Save/Load use default preferences where appropriate ([9ecc19d](https://github.com/joshuar/go-hass-anything/commit/9ecc19d9053b7d0036ee3473aafaf87249c9bcec))
* **mqtt,config:** :bug: prefs propagation fixes ([8b720f0](https://github.com/joshuar/go-hass-anything/commit/8b720f089b1777b6d9abb9e873940f1f5cc0ab2c))

## [3.0.0](https://github.com/joshuar/go-hass-anything/compare/v2.0.3...v3.0.0) (2024-01-25)


### ⚠ BREAKING CHANGES

* **all:** update import path for breaking change
* **config:** rename exported struct AppPreferences -> Preferences
* **config,agent,mqtt:** allow specifying a path to MQTT config file

### Features

* **all:** update import path for breaking change ([866e894](https://github.com/joshuar/go-hass-anything/commit/866e8943d9eeda8bf6321d8cf865068d99ac34e8))
* **config,agent,mqtt:** allow specifying a path to MQTT config file ([48eb067](https://github.com/joshuar/go-hass-anything/commit/48eb0674520f5d54db09e36076232e3b65884485))


### Code Refactoring

* **config:** rename exported struct AppPreferences -&gt; Preferences ([b849504](https://github.com/joshuar/go-hass-anything/commit/b84950424402e7e255eb8b4c60ea7fad272b4fb1))

## [2.0.3](https://github.com/joshuar/go-hass-anything/compare/v2.0.2...v2.0.3) (2024-01-25)


### Bug Fixes

* **all:** :ambulance: fix module path after new major version ([00506af](https://github.com/joshuar/go-hass-anything/commit/00506afc14cb55763efaa610a749bfbc1cd693a3))

## [2.0.2](https://github.com/joshuar/go-hass-anything/compare/v2.0.1...v2.0.2) (2024-01-24)


### Miscellaneous Chores

* release 2.0.2 ([6a6343a](https://github.com/joshuar/go-hass-anything/commit/6a6343a94a35ea66ba1f0ff1b84fb0c545c59eb0))

## [2.0.1](https://github.com/joshuar/go-hass-anything/compare/v2.0.0...v2.0.1) (2024-01-24)


### Bug Fixes

* **config:** :recycle: ensure tomlconfig has own package name following conventions ([8514f5f](https://github.com/joshuar/go-hass-anything/commit/8514f5f6127d89bce787f34114689e07b5624974))

## [2.0.0](https://github.com/joshuar/go-hass-anything/compare/v1.4.0...v2.0.0) (2024-01-24)


### ⚠ BREAKING CHANGES

* **all:** major refactor

### Features

* **all:** major refactor ([1fa84be](https://github.com/joshuar/go-hass-anything/commit/1fa84bed3633fcb1ab9cbc68b5c44bb069286403))

## [1.4.0](https://github.com/joshuar/go-hass-anything/compare/v1.3.0...v1.4.0) (2024-01-24)


### Features

* **tools:** :sparkles: appgenerator path adjustment ([7387f7b](https://github.com/joshuar/go-hass-anything/commit/7387f7ba4d6f4929f60c3a985c78ef2b0a0e9567))
* **web:** :sparkles: add retries to web requests ([1f29917](https://github.com/joshuar/go-hass-anything/commit/1f29917ad9b9b1ddf2f39cb6e64f823fe7901e9c))

## [1.3.0](https://github.com/joshuar/go-hass-anything/compare/v1.2.0...v1.3.0) (2024-01-09)


### Features

* **config:** add a new simple toml-based config package ([8afda51](https://github.com/joshuar/go-hass-anything/commit/8afda51cf97bd5af73ff7e2c303d4c04124538cc))
* **web:** simplify web request handling ([3f5ee46](https://github.com/joshuar/go-hass-anything/commit/3f5ee46e2e761b24f0f6cf9f5224bb4a0222574d))


### Bug Fixes

* **container:** easier inclusion of own apps into container image ([4fee795](https://github.com/joshuar/go-hass-anything/commit/4fee795ba268baf354323d0f1cd6da395d0bf3eb))
* **container:** fix Dockerfile, adjust README ([b7b9e51](https://github.com/joshuar/go-hass-anything/commit/b7b9e5106862afd3986fa4d122010e131bfba892))
* **container:** fully qualify base image ([6f79fcd](https://github.com/joshuar/go-hass-anything/commit/6f79fcdfa4b582bcfacc0c6d171a0f3027a030fe))
* **container:** ignore failure if apps dir is not a symlink ([342da68](https://github.com/joshuar/go-hass-anything/commit/342da684568b864663336c7ec30d043654d598e3))

## [1.2.0](https://github.com/joshuar/go-hass-anything/compare/v1.1.0...v1.2.0) (2023-10-28)


### Features

* **mqtt:** use a retry-backoff for initial mqtt connection ([541c64a](https://github.com/joshuar/go-hass-anything/commit/541c64a167694afb7b46a1a032f3f879e81f020d))


### Bug Fixes

* **config:** config start logic ([ad18869](https://github.com/joshuar/go-hass-anything/commit/ad18869074a8d59bc555cd5f4eb0dacd8d0fff22))

## [1.1.0](https://github.com/joshuar/go-hass-anything/compare/v1.0.1...v1.1.0) (2023-10-11)


### Features

* **config:** embed the app version for use within the code ([1f1d643](https://github.com/joshuar/go-hass-anything/commit/1f1d643992f1c86081250bedd6f1507d7ec61d66))
* **config:** re-do config ([9c43341](https://github.com/joshuar/go-hass-anything/commit/9c43341cbd07fa27d30833a2f3874b67f4e105ce))
* **hass:** add Register and UnRegister functions ([126f753](https://github.com/joshuar/go-hass-anything/commit/126f753fb5ed649986de58e88af569226a29576a))
* **hass:** improve entity configs ([7aefbee](https://github.com/joshuar/go-hass-anything/commit/7aefbeef42b6eed149307616f81345c6ea2498d0))
* **init:** add a prestart command for systemd service ([063c585](https://github.com/joshuar/go-hass-anything/commit/063c585dab37a9225866214dc21281953f7767c5))


### Bug Fixes

* **appgenerator:** better app detection ([e6f1b3a](https://github.com/joshuar/go-hass-anything/commit/e6f1b3afbbc68ddf7189dd062b5cf2c0a83912d5))
* **config:** moq interfaces ([659b7fd](https://github.com/joshuar/go-hass-anything/commit/659b7fdfba4d11c21c7d64457f31549f877cc53f))
* **config:** registration methods do not need any inputs ([fc48486](https://github.com/joshuar/go-hass-anything/commit/fc484866976ebfc25c32d79588380dcdeaa76392))
* **exampleapp:** simplify Clear function ([1ceebf8](https://github.com/joshuar/go-hass-anything/commit/1ceebf82e4118df89c5c4f2dec27f324cee2db70))


### Reverts

* **init:** don't do prestart ([acdee55](https://github.com/joshuar/go-hass-anything/commit/acdee55b91750813ab15cbe7667b9ec27912d218))

## [1.0.1](https://github.com/joshuar/go-hass-anything/compare/v1.0.0...v1.0.1) (2023-10-09)


### Bug Fixes

* **apps:** move helpers to external pkg dir ([0358ab0](https://github.com/joshuar/go-hass-anything/commit/0358ab070b4516746bd3f5883017aadbb432dd61))

## 1.0.0 (2023-10-04)


### Features

* **agent:** autogenerate list of apps ([5bb2f7f](https://github.com/joshuar/go-hass-anything/commit/5bb2f7fda6b2e38878e662ae4d1900376a1342fd))
* **all:** add a bunch of helpers for entities, use in Example App ([8e53c5d](https://github.com/joshuar/go-hass-anything/commit/8e53c5d0c09accae819ca6750555f1b3c07af7b5))
* **all:** add a systemd service file ([74d8e54](https://github.com/joshuar/go-hass-anything/commit/74d8e54a630aeff6f0cc852e6d4b23730f37a04f))
* **all:** add an example app ([f1a09c8](https://github.com/joshuar/go-hass-anything/commit/f1a09c8892d8f16198c331d898ea9ada84582bcf))
* **all:** initial commit ([5450059](https://github.com/joshuar/go-hass-anything/commit/5450059fddaddc4d82abe767f1316ef2dc3aaddf))
* **web:** add ability to customise request timeout ([c79c61e](https://github.com/joshuar/go-hass-anything/commit/c79c61eca72f127d3016592f6b6cb36c3569ee24))


### Bug Fixes

* **apps:** actually add example app ([723e826](https://github.com/joshuar/go-hass-anything/commit/723e82682755d535384523dc64536edbf905d4d4))
