# Changelog

## [12.0.0](https://github.com/joshuar/go-hass-anything/compare/v11.1.0...v12.0.0) (2024-10-10)


### ⚠ BREAKING CHANGES

* The invocation for constructing an entity with builders has changed. Please see the example apps which have been updated to utilise the new methods for what needs to change.

### Bug Fixes

* **hass:** :bug: set camera topic before validating ([8b37801](https://github.com/joshuar/go-hass-anything/commit/8b378012de0a2cbef5a348e6ae0cbd02b93ef0d2))


### Performance Improvements

* **examples:** :zap: additional camera example app improvements ([7119a34](https://github.com/joshuar/go-hass-anything/commit/7119a349f24a562be40cd15083295faada1afaba))


### Code Refactoring

* :building_construction: better builder methods and validation of entities ([1370d95](https://github.com/joshuar/go-hass-anything/commit/1370d9575d1e100ce4f209547b1d7b436d199ffb))

## [11.1.0](https://github.com/joshuar/go-hass-anything/compare/v11.0.0...v11.1.0) (2024-08-10)


### Features

* **hass:** :sparkles: add image and camera entities ([c44d321](https://github.com/joshuar/go-hass-anything/commit/c44d321443b62d3361fec290b0c5a1972442067f))


### Bug Fixes

* **container:** :heavy_plus_sign: ensure linux-headers are installed for camera example app ([7efc46e](https://github.com/joshuar/go-hass-anything/commit/7efc46e1f26c1c145285cef595e30c04780e0273))
* **mqtt:** :bug: don't assume json encoding of message ([fd218cd](https://github.com/joshuar/go-hass-anything/commit/fd218cd83a9ef5d62a29d4249b8d5b15922014cc))

## [11.0.0](https://github.com/joshuar/go-hass-anything/compare/v10.1.0...v11.0.0) (2024-07-24)


### ⚠ BREAKING CHANGES

* migrate cobra->kong and zerolog->slog
* **agent:** improve app preference handling

### Features

* **agent:** improve app preference handling ([0c36708](https://github.com/joshuar/go-hass-anything/commit/0c3670812e03f68bc4536587d3326002ddf44642))
* **examples:** :recycle: split example into multiple apps to make it easier to follow and copy ([aede219](https://github.com/joshuar/go-hass-anything/commit/aede219589a678452dd7f95f53c1808443eac783))
* **hass:** :sparkles: add text entity type ([ddeb37e](https://github.com/joshuar/go-hass-anything/commit/ddeb37ec9758655e487b97d9c6d429d487dabbd3))
* migrate cobra-&gt;kong and zerolog->slog ([0926713](https://github.com/joshuar/go-hass-anything/commit/0926713685025fdb9bb7fb84e41e2d4a804d0398))


### Bug Fixes

* **agent:** :zap: gracefully handle sigterm ([21b9ad4](https://github.com/joshuar/go-hass-anything/commit/21b9ad484aa52e275041bc9ab989a951f6318028))


### Performance Improvements

* **agent:** :zap: improve handling of different types of apps and remove need for a waitgroup ([feb1991](https://github.com/joshuar/go-hass-anything/commit/feb1991dbcd990274bd7a849fdccca1aa4054b44))

## [10.1.0](https://github.com/joshuar/go-hass-anything/compare/v10.0.0...v10.1.0) (2024-07-06)


### Features

* **container:** :sparkles: switch to alpine base image for container ([676d6ba](https://github.com/joshuar/go-hass-anything/commit/676d6ba83e08f620de80d84a282d3922758d18d1))


### Bug Fixes

* **container:** :bug: container should run as (configurable) non-root user ([a7ba398](https://github.com/joshuar/go-hass-anything/commit/a7ba3988749d7e3ba62f670401c560e575543324))

## [10.0.0](https://github.com/joshuar/go-hass-anything/compare/v9.2.1...v10.0.0) (2024-06-25)


### ⚠ BREAKING CHANGES

* **mqtt:** require passing context when publishing
* rewrite preferences rewrite
* use koanf for app/agent preferences
* **web:** switch to resty for web requests
* **web:** switch to resty for web requests

### Features

* **container:** :sparkles: use mage for container build and support multiarch ([f10dfdd](https://github.com/joshuar/go-hass-anything/commit/f10dfdd8375e1642bdd54d8590e9056a4d831c5a))
* **mqtt:** require passing context when publishing ([b1305a6](https://github.com/joshuar/go-hass-anything/commit/b1305a614c020358d649ee2fda701dfb81ead076))
* **preferences:** :recycle: app preferences now use the same underlying structure as agent preferences ([89694e8](https://github.com/joshuar/go-hass-anything/commit/89694e844f5c41d9b8bc0edfc3876ea49a319d26))
* **preferences:** :sparkles: add a Keys method to get the list of preference keys ([eed4836](https://github.com/joshuar/go-hass-anything/commit/eed4836068ceafbd8483573f576a4ef6342043f6))
* rewrite preferences rewrite ([f8f32aa](https://github.com/joshuar/go-hass-anything/commit/f8f32aaac7b67b586aca32634d53d540e31446bf))
* **ui:** :recycle: ui management for both agent and apps ([9126f38](https://github.com/joshuar/go-hass-anything/commit/9126f3896a35935527161ff9857ef765d90ac964))
* use koanf for app/agent preferences ([d3b412a](https://github.com/joshuar/go-hass-anything/commit/d3b412a128dc0eb740cd7bdac9fdb8ed036c56b7))
* **web:** switch to resty for web requests ([1287ecf](https://github.com/joshuar/go-hass-anything/commit/1287ecf44dda4a9c471007ec6f6446b12326bf3e))
* **web:** switch to resty for web requests ([31c5326](https://github.com/joshuar/go-hass-anything/commit/31c5326b67c58c3e8871826dfe813772688513c0))


### Bug Fixes

* **container:** :bug: ensure container uses correct arch ([f02e096](https://github.com/joshuar/go-hass-anything/commit/f02e0967d20a4ee61583fa9a4fe88353cf5a7bb6))
* **container:** :bug: use TARGETARCH during build stage ([f1ce406](https://github.com/joshuar/go-hass-anything/commit/f1ce406d2fb18e2dfe542ee8c29f28f426f9ae39))
* **preferences:** :bug: fix loading preferences on initial run ([40832a6](https://github.com/joshuar/go-hass-anything/commit/40832a673dbb7b626fe15154b470440c67c6fcf3))
* **ui:** :lipstick: increase input field character limit ([c1644af](https://github.com/joshuar/go-hass-anything/commit/c1644af7c4e40cf90ad123ff83bad91d6818b322))

## [9.2.1](https://github.com/joshuar/go-hass-anything/compare/v9.2.0...v9.2.1) (2024-06-02)


### Bug Fixes

* **hass:** :children_crossing: set a default origin/device if none specified on entities ([708acc2](https://github.com/joshuar/go-hass-anything/commit/708acc243ac531936c189e1bec1cc0486839854d))

## [9.2.0](https://github.com/joshuar/go-hass-anything/compare/v9.1.0...v9.2.0) (2024-05-22)


### Features

* **mqtt:** :sparkles: publish app configs once MQTT connection is established ([74f59da](https://github.com/joshuar/go-hass-anything/commit/74f59da7f133c0a6298b25aa869bd01ae542d49b))


### Bug Fixes

* **mqtt:** :bug: re-add removed user/pass settings when connecting to MQTT ([c6097cd](https://github.com/joshuar/go-hass-anything/commit/c6097cd3955664e8c309f51fe335d36e7a6863e9))
* **mqtt:** :safety_vest: protect against potential nil panics ([7e20c1f](https://github.com/joshuar/go-hass-anything/commit/7e20c1f452bbe4221c320c165e0c0f50e0d78df7))

## [9.1.0](https://github.com/joshuar/go-hass-anything/compare/v9.0.0...v9.1.0) (2024-05-04)


### Features

* **mqtt:** :sparkles: Republish app configs after Home Assistant restarts ([2bce61a](https://github.com/joshuar/go-hass-anything/commit/2bce61ad280d007db70e53172f5934c5d1dada2d))

## [9.0.0](https://github.com/joshuar/go-hass-anything/compare/v8.0.0...v9.0.0) (2024-05-04)


### ⚠ BREAKING CHANGES

* the app interface the agent expects has been completely overhauled. Apps now don't need to access the MQTT client to publish messages directly. They can instead specify whether they require polling on an interval or are app driven (or even run once-only) and the agent will set up the necessary functionality to provide it. Additionally, the underlying MQTT client in use has switched to a newer v5 based one that should be more performant and feature-ful for future functionality.

### Features

* major internal rewrite for app interface and switch MQTT client library ([7c5c888](https://github.com/joshuar/go-hass-anything/commit/7c5c888596be33d7f760bf8c1fdae43c589b889e))

## [8.0.0](https://github.com/joshuar/go-hass-anything/compare/v7.2.1...v8.0.0) (2024-05-01)


### ⚠ BREAKING CHANGES

* use generics for entities

### Features

* use generics for entities ([c632384](https://github.com/joshuar/go-hass-anything/commit/c6323845f65ec2f42b80d3a1c6435fad5b3f354b))

## [7.2.1](https://github.com/joshuar/go-hass-anything/compare/v7.2.0...v7.2.1) (2024-05-01)


### Miscellaneous Chores

* release 7.2.1 ([4045e43](https://github.com/joshuar/go-hass-anything/commit/4045e437cce3cd408bdf4f4ef2cfe145113c5640))
* release 7.2.1 ([8ea2fa4](https://github.com/joshuar/go-hass-anything/commit/8ea2fa4b55c3cb259f1993f064b1583174eeadfd))

## [7.2.0](https://github.com/joshuar/go-hass-anything/compare/v7.1.0...v7.2.0) (2024-04-29)


### Features

* **examples:** :sparkles: update example app to demonstrate a number entity ([5f5a30b](https://github.com/joshuar/go-hass-anything/commit/5f5a30b7cdb17c5d049e3bb3627e9fe7809e9ade))
* **hass:** :sparkles: add support for number entities ([042681d](https://github.com/joshuar/go-hass-anything/commit/042681dd2b4de8a736f75e1f80ad02e2929112ef))

## [7.1.0](https://github.com/joshuar/go-hass-anything/compare/v7.0.0...v7.1.0) (2024-04-13)


### Features

* **mqtt:** :zap: use github.com/sourcegraph/conc/pool for sending messages to MQTT ([ee5da5c](https://github.com/joshuar/go-hass-anything/commit/ee5da5c8a7231e47661463f4dedaa4def9f919d2))

## [7.0.0](https://github.com/joshuar/go-hass-anything/compare/v6.0.0...v7.0.0) (2024-04-12)


### ⚠ BREAKING CHANGES

* this change updates the method for publishing app states and thus existing apps will need to make adjustments. See the docs for details.

### Features

* connection resilience ([0ca4ea6](https://github.com/joshuar/go-hass-anything/commit/0ca4ea654049c2318aabba2cda52891953d6eab2))


### Bug Fixes

* **examples:** :bug: correct client argument ([e63f8d7](https://github.com/joshuar/go-hass-anything/commit/e63f8d729c9d0ba49fc27b43c43babac51092384))

## [6.0.0](https://github.com/joshuar/go-hass-anything/compare/v5.0.1...v6.0.0) (2024-03-11)


### ⚠ BREAKING CHANGES

* Refactor to support MQTT brokers configured without persistence. The agent will now register all apps on startup and re-register   if Home Assistant is restarted. The code supports similar functionality for when imported as a package.

### Features

* zero persistence MQTT support ([96e5b2c](https://github.com/joshuar/go-hass-anything/commit/96e5b2c709af22b0f9ee6228e94ccb578cc4e126))


### Bug Fixes

* **tools:** :bug: fix missing version number for example app import path ([44b2937](https://github.com/joshuar/go-hass-anything/commit/44b293764db7491750d5d4e832dcaa7bc4de35ac))

## [5.0.1](https://github.com/joshuar/go-hass-anything/compare/v5.0.0...v5.0.1) (2024-02-26)


### Miscellaneous Chores

* release 5.0.1 ([8ab1115](https://github.com/joshuar/go-hass-anything/commit/8ab1115529e3760361205078b578e79c3b2530e0))

## [5.0.0](https://github.com/joshuar/go-hass-anything/compare/v4.0.0...v5.0.0) (2024-02-09)


### ⚠ BREAKING CHANGES

* Creating an MQTT client now requires passing a context. This supports cancellation to avoid a connection problem to the broker resulting in inifinite retries.

### Features

* context aware MQTT client creation ([63a125d](https://github.com/joshuar/go-hass-anything/commit/63a125d175e5728d16a0e15de6d70c66af575c2c))


### Bug Fixes

* **agent:** :bug: save new preferences when no existing preferences file ([774e09b](https://github.com/joshuar/go-hass-anything/commit/774e09b1a460df56a395fccfaff22a177f42beff))
* **examples:** :bug: fix api path in exampleapp ([71a9e35](https://github.com/joshuar/go-hass-anything/commit/71a9e3522d49bd53e0f0124779187c78a52847a2))
* **tools:** :bug: update appgenerator template for new api version ([4a82d47](https://github.com/joshuar/go-hass-anything/commit/4a82d47cfd49a18ebc5995b3e862895e1d3342ec))

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
