# Changelog

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
