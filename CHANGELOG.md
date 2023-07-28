# Changelog

## [v1.9.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.9.0) (2023-07-28)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.8.0...v1.9.0)

**Merged pull requests:**

- Release/1.9.0 [\#120](https://github.com/nifcloud/terraform-provider-nifcloud/pull/120) ([fuku2014](https://github.com/fuku2014))
- Feature/ravgw [\#119](https://github.com/nifcloud/terraform-provider-nifcloud/pull/119) ([fuku2014](https://github.com/fuku2014))
- Feature/paas close [\#118](https://github.com/nifcloud/terraform-provider-nifcloud/pull/118) ([fuku2014](https://github.com/fuku2014))
- Fix acceptance test parallelism: 7 -\> 4 [\#117](https://github.com/nifcloud/terraform-provider-nifcloud/pull/117) ([aokumasan](https://github.com/aokumasan))
- Fix/lint errors [\#116](https://github.com/nifcloud/terraform-provider-nifcloud/pull/116) ([aokumasan](https://github.com/aokumasan))

## [v1.8.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.8.0) (2023-04-20)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.7.0...v1.8.0)

**Closed issues:**

- nifcloud\_dns\_record: Can not recognized @ record [\#111](https://github.com/nifcloud/terraform-provider-nifcloud/issues/111)

**Merged pull requests:**

- Feature/update nifcloud [\#115](https://github.com/nifcloud/terraform-provider-nifcloud/pull/115) ([fuku2014](https://github.com/fuku2014))
- Fix doc of load\_balancer\#import [\#114](https://github.com/nifcloud/terraform-provider-nifcloud/pull/114) ([ystkfujii](https://github.com/ystkfujii))
- Support dns record name shorthand [\#113](https://github.com/nifcloud/terraform-provider-nifcloud/pull/113) ([SogoKato](https://github.com/SogoKato))
- Change behaviour of dns record flattener [\#112](https://github.com/nifcloud/terraform-provider-nifcloud/pull/112) ([SogoKato](https://github.com/SogoKato))
- Update examples for change Ubuntu verison [\#110](https://github.com/nifcloud/terraform-provider-nifcloud/pull/110) ([fuku2014](https://github.com/fuku2014))
- Fix doc typo nifcloud [\#109](https://github.com/nifcloud/terraform-provider-nifcloud/pull/109) ([tanopanta](https://github.com/tanopanta))

## [v1.7.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.7.0) (2022-09-01)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.6.0...v1.7.0)

**Merged pull requests:**

- Update RDB storage types [\#108](https://github.com/nifcloud/terraform-provider-nifcloud/pull/108) ([tily](https://github.com/tily))
- Feature/add storage bucket [\#107](https://github.com/nifcloud/terraform-provider-nifcloud/pull/107) ([aokumasan](https://github.com/aokumasan))

## [v1.6.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.6.0) (2022-03-30)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.5.1...v1.6.0)

**Merged pull requests:**

- doc\) Add subcategory ESS [\#106](https://github.com/nifcloud/terraform-provider-nifcloud/pull/106) ([fuku2014](https://github.com/fuku2014))
- add ess [\#105](https://github.com/nifcloud/terraform-provider-nifcloud/pull/105) ([fuku2014](https://github.com/fuku2014))
- Feature/update sdk [\#104](https://github.com/nifcloud/terraform-provider-nifcloud/pull/104) ([fuku2014](https://github.com/fuku2014))
- Add nifcloud\_dns\_record resource [\#103](https://github.com/nifcloud/terraform-provider-nifcloud/pull/103) ([tunakyonn](https://github.com/tunakyonn))
- resource/nifcloud\_route\_table Fix to exclude EnableVgwRoutePropagation [\#102](https://github.com/nifcloud/terraform-provider-nifcloud/pull/102) ([fuku2014](https://github.com/fuku2014))
- Add nifcloud\_dns\_zone resource [\#101](https://github.com/nifcloud/terraform-provider-nifcloud/pull/101) ([tunakyonn](https://github.com/tunakyonn))
- Fix document of db instance [\#100](https://github.com/nifcloud/terraform-provider-nifcloud/pull/100) ([SogoKato](https://github.com/SogoKato))

## [v1.5.1](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.5.1) (2021-10-21)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.5.0...v1.5.1)

**Merged pull requests:**

- resource/nifcloud\_security\_group\_rule: Change import specification [\#98](https://github.com/nifcloud/terraform-provider-nifcloud/pull/98) ([fuku2014](https://github.com/fuku2014))

## [v1.5.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.5.0) (2021-09-17)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.4.0...v1.5.0)

**Implemented enhancements:**

- resource/nifcloud\_hatoba\_cluster: Add support kube\_config attribute [\#95](https://github.com/nifcloud/terraform-provider-nifcloud/issues/95)

**Closed issues:**

- Router creation fails with v.1.4.0 [\#88](https://github.com/nifcloud/terraform-provider-nifcloud/issues/88)
- There was a network creation error [\#85](https://github.com/nifcloud/terraform-provider-nifcloud/issues/85)
- "This NetworkName is processing" error [\#74](https://github.com/nifcloud/terraform-provider-nifcloud/issues/74)
- RDBの作成でステータスが「エラー」になった場合の即時復帰について [\#72](https://github.com/nifcloud/terraform-provider-nifcloud/issues/72)

**Merged pull requests:**

- Feature/hatoba kube config [\#97](https://github.com/nifcloud/terraform-provider-nifcloud/pull/97) ([fuku2014](https://github.com/fuku2014))
- Fix acceptance workflow branch name [\#96](https://github.com/nifcloud/terraform-provider-nifcloud/pull/96) ([tunakyonn](https://github.com/tunakyonn))
- Fix/network mutex [\#93](https://github.com/nifcloud/terraform-provider-nifcloud/pull/93) ([fuku2014](https://github.com/fuku2014))
- resource/nifcloud\_elastic\_ip: Fix bug of can not import private ip [\#92](https://github.com/nifcloud/terraform-provider-nifcloud/pull/92) ([fuku2014](https://github.com/fuku2014))
- Update to go 1.17 version [\#91](https://github.com/nifcloud/terraform-provider-nifcloud/pull/91) ([fuku2014](https://github.com/fuku2014))
- Add debug logger to dump the API request and response [\#90](https://github.com/nifcloud/terraform-provider-nifcloud/pull/90) ([aokumasan](https://github.com/aokumasan))
- resource/nifcloud\_load\_balancer: Fix bug of flatten options [\#89](https://github.com/nifcloud/terraform-provider-nifcloud/pull/89) ([fuku2014](https://github.com/fuku2014))
- Fix/sdk version [\#94](https://github.com/nifcloud/terraform-provider-nifcloud/pull/94) ([fuku2014](https://github.com/fuku2014))
- Fix/import state verify ignore [\#87](https://github.com/nifcloud/terraform-provider-nifcloud/pull/87) ([fuku2014](https://github.com/fuku2014))
- Feature/separate instance rule [\#79](https://github.com/nifcloud/terraform-provider-nifcloud/pull/79) ([matsuoka-k-git](https://github.com/matsuoka-k-git))

## [v1.4.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.4.0) (2021-07-06)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.3.2...v1.4.0)

**Implemented enhancements:**

- resource/nifcloud\_volume: Add support detach/attach [\#52](https://github.com/nifcloud/terraform-provider-nifcloud/issues/52)

**Closed issues:**

- Specifying an existing ID for network\_id of nifcloud\_vpn\_gateway results in an error [\#84](https://github.com/nifcloud/terraform-provider-nifcloud/issues/84)
- Add a VM after creating VPN resources results in error [\#69](https://github.com/nifcloud/terraform-provider-nifcloud/issues/69)

**Merged pull requests:**

- resource/nifcloud\_vpn\_gateway:  Remove network\_id validation [\#86](https://github.com/nifcloud/terraform-provider-nifcloud/pull/86) ([fuku2014](https://github.com/fuku2014))
- resource/nifcloud\_db\_security\_group: Fix lint error [\#83](https://github.com/nifcloud/terraform-provider-nifcloud/pull/83) ([fuku2014](https://github.com/fuku2014))
- Fix nat table for terrafrom v1.0.0 [\#82](https://github.com/nifcloud/terraform-provider-nifcloud/pull/82) ([tunakyonn](https://github.com/tunakyonn))
- resource/nifcloud\_volume: Fix save instance\_id to state [\#81](https://github.com/nifcloud/terraform-provider-nifcloud/pull/81) ([fuku2014](https://github.com/fuku2014))
- Update nifclous-sdk-go to v1.7.0 [\#80](https://github.com/nifcloud/terraform-provider-nifcloud/pull/80) ([fuku2014](https://github.com/fuku2014))
- Fix/hatoba cluster bug [\#78](https://github.com/nifcloud/terraform-provider-nifcloud/pull/78) ([aokumasan](https://github.com/aokumasan))
- Improve the AD server provisioning script for NAS instance test [\#77](https://github.com/nifcloud/terraform-provider-nifcloud/pull/77) ([aokumasan](https://github.com/aokumasan))
- Fix DB security group rule waiter: wait until target rule is revoked [\#76](https://github.com/nifcloud/terraform-provider-nifcloud/pull/76) ([aokumasan](https://github.com/aokumasan))
- Add dhcp default value for nifcloud\_router [\#75](https://github.com/nifcloud/terraform-provider-nifcloud/pull/75) ([aokumasan](https://github.com/aokumasan))
- Fix/nifcloud nas instance bug [\#73](https://github.com/nifcloud/terraform-provider-nifcloud/pull/73) ([aokumasan](https://github.com/aokumasan))
- Feature/add support terraform 1.0.0 [\#71](https://github.com/nifcloud/terraform-provider-nifcloud/pull/71) ([fuku2014](https://github.com/fuku2014))
- Fix/vpn connection [\#70](https://github.com/nifcloud/terraform-provider-nifcloud/pull/70) ([fuku2014](https://github.com/fuku2014))
- Feature/hatoba cluster [\#68](https://github.com/nifcloud/terraform-provider-nifcloud/pull/68) ([aokumasan](https://github.com/aokumasan))
- Fix nifcloud\_security\_group name [\#67](https://github.com/nifcloud/terraform-provider-nifcloud/pull/67) ([matsuoka-k-git](https://github.com/matsuoka-k-git))
- Feature/hatoba firewall group [\#66](https://github.com/nifcloud/terraform-provider-nifcloud/pull/66) ([aokumasan](https://github.com/aokumasan))
- Feature/nas instance [\#65](https://github.com/nifcloud/terraform-provider-nifcloud/pull/65) ([aokumasan](https://github.com/aokumasan))
- resource/nifcloud\_instance: Fix buf of flatten network\_interface [\#64](https://github.com/nifcloud/terraform-provider-nifcloud/pull/64) ([fuku2014](https://github.com/fuku2014))
- Feature/nas security group [\#63](https://github.com/nifcloud/terraform-provider-nifcloud/pull/63) ([aokumasan](https://github.com/aokumasan))
- Fix to use latest nifcloud-sdk-go [\#62](https://github.com/nifcloud/terraform-provider-nifcloud/pull/62) ([aokumasan](https://github.com/aokumasan))

## [v1.3.2](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.3.2) (2021-04-20)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.3.1...v1.3.2)

**Merged pull requests:**

- Add only detach attach operation for volume [\#61](https://github.com/nifcloud/terraform-provider-nifcloud/pull/61) ([tunakyonn](https://github.com/tunakyonn))

## [v1.3.1](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.3.1) (2021-04-19)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.3.0...v1.3.1)

**Implemented enhancements:**

- nifcloud\_vpn\_gateway: Add attribute public ip address [\#51](https://github.com/nifcloud/terraform-provider-nifcloud/issues/51)

**Fixed bugs:**

- resource/nifcloud\_load\_balancer: Fix bug for import error message and ip\_version not saved in state  [\#54](https://github.com/nifcloud/terraform-provider-nifcloud/issues/54)
- resource/nifcloud\_instance: Fix bug for can't remove the FW  [\#53](https://github.com/nifcloud/terraform-provider-nifcloud/issues/53)

**Merged pull requests:**

- Add support detach attach for volume [\#60](https://github.com/nifcloud/terraform-provider-nifcloud/pull/60) ([tunakyonn](https://github.com/tunakyonn))
- Fix testSweepDbSecurityGroup [\#59](https://github.com/nifcloud/terraform-provider-nifcloud/pull/59) ([tunakyonn](https://github.com/tunakyonn))
- Fix extend volume size waiter [\#58](https://github.com/nifcloud/terraform-provider-nifcloud/pull/58) ([tunakyonn](https://github.com/tunakyonn))
- Add public\_ip\_address in nifcloud\_vpn\_gateway [\#57](https://github.com/nifcloud/terraform-provider-nifcloud/pull/57) ([fuku2014](https://github.com/fuku2014))
- Fix bugs of nifcloud\_instance [\#56](https://github.com/nifcloud/terraform-provider-nifcloud/pull/56) ([fuku2014](https://github.com/fuku2014))
- Fix bugs of nifcloud\_load\_balancer [\#55](https://github.com/nifcloud/terraform-provider-nifcloud/pull/55) ([fuku2014](https://github.com/fuku2014))

## [v1.3.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.3.0) (2021-03-19)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.2.0...v1.3.0)

**Implemented enhancements:**

- nifcloud\_network\_interface: Add resource [\#43](https://github.com/nifcloud/terraform-provider-nifcloud/issues/43)
- nifcloud\_vpn\_connection: Add resource [\#13](https://github.com/nifcloud/terraform-provider-nifcloud/issues/13)

**Merged pull requests:**

- Fixed a bug that load balancer creatation cause an error in multi account [\#49](https://github.com/nifcloud/terraform-provider-nifcloud/pull/49) ([fuku2014](https://github.com/fuku2014))
- Change revoke\_rules\_on\_delete default to true [\#48](https://github.com/nifcloud/terraform-provider-nifcloud/pull/48) ([fuku2014](https://github.com/fuku2014))
- Fix/add description to security group rule example [\#47](https://github.com/nifcloud/terraform-provider-nifcloud/pull/47) ([tily](https://github.com/tily))
- refactor directory structure [\#46](https://github.com/nifcloud/terraform-provider-nifcloud/pull/46) ([fuku2014](https://github.com/fuku2014))
- Feature/network interface [\#45](https://github.com/nifcloud/terraform-provider-nifcloud/pull/45) ([fuku2014](https://github.com/fuku2014))
- resource/nifcloud\_vpn\_connection: Add resource [\#44](https://github.com/nifcloud/terraform-provider-nifcloud/pull/44) ([tunakyonn](https://github.com/tunakyonn))
- Release/v1.3.0 [\#50](https://github.com/nifcloud/terraform-provider-nifcloud/pull/50) ([fuku2014](https://github.com/fuku2014))

## [v1.2.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.2.0) (2021-02-26)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.1.0...v1.2.0)

**Implemented enhancements:**

- nifcloud\_db\_instance: Add resource [\#16](https://github.com/nifcloud/terraform-provider-nifcloud/issues/16)
- nifcloud\_db\_security\_group: Add resource [\#15](https://github.com/nifcloud/terraform-provider-nifcloud/issues/15)
- nifcloud\_db\_parameter\_group: Add resource [\#14](https://github.com/nifcloud/terraform-provider-nifcloud/issues/14)
- nifcloud\_vpn\_gateway: Add resource [\#12](https://github.com/nifcloud/terraform-provider-nifcloud/issues/12)
- nifcloud\_lb: Add resource [\#3](https://github.com/nifcloud/terraform-provider-nifcloud/issues/3)

**Merged pull requests:**

- Release/v1.2.0 [\#42](https://github.com/nifcloud/terraform-provider-nifcloud/pull/42) ([tily](https://github.com/tily))
- Feature/add load balancer listener [\#41](https://github.com/nifcloud/terraform-provider-nifcloud/pull/41) ([o108minmin](https://github.com/o108minmin))
- Feature/add vpn gateway [\#39](https://github.com/nifcloud/terraform-provider-nifcloud/pull/39) ([nashik](https://github.com/nashik))
- add load balancer [\#38](https://github.com/nifcloud/terraform-provider-nifcloud/pull/38) ([o108minmin](https://github.com/o108minmin))
- Feature/db parameter group [\#36](https://github.com/nifcloud/terraform-provider-nifcloud/pull/36) ([aokumasan](https://github.com/aokumasan))
- resource/nifcloud\_db\_security\_group: Add resource [\#35](https://github.com/nifcloud/terraform-provider-nifcloud/pull/35) ([tunakyonn](https://github.com/tunakyonn))
- ci: Fix to execute acc only master and release branch [\#40](https://github.com/nifcloud/terraform-provider-nifcloud/pull/40) ([fuku2014](https://github.com/fuku2014))
- resource/nifcloud\_db\_instance: Add resource [\#37](https://github.com/nifcloud/terraform-provider-nifcloud/pull/37) ([fuku2014](https://github.com/fuku2014))

## [v1.1.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.1.0) (2021-01-29)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.0.3...v1.1.0)

**Implemented enhancements:**

- nifcloud\_elb: Add resource [\#11](https://github.com/nifcloud/terraform-provider-nifcloud/issues/11)
- nifcloud\_web\_proxy: Add resource [\#10](https://github.com/nifcloud/terraform-provider-nifcloud/issues/10)
- nifcloud\_router: Add resource [\#9](https://github.com/nifcloud/terraform-provider-nifcloud/issues/9)
- nifcloud\_dhcp\_config: Add resource [\#8](https://github.com/nifcloud/terraform-provider-nifcloud/issues/8)
- nifcloud\_dhcp\_option: Add resource [\#7](https://github.com/nifcloud/terraform-provider-nifcloud/issues/7)
- nifcloud\_nat\_table: Add resource [\#6](https://github.com/nifcloud/terraform-provider-nifcloud/issues/6)
- nifcloud\_route\_table: Add resource [\#5](https://github.com/nifcloud/terraform-provider-nifcloud/issues/5)
- nifcloud\_customer\_gateway: Add resource [\#4](https://github.com/nifcloud/terraform-provider-nifcloud/issues/4)

**Closed issues:**

- docs: Change subcategory to same as the control panel [\#19](https://github.com/nifcloud/terraform-provider-nifcloud/issues/19)

**Merged pull requests:**

- docs: Change subcategory to same as the control panel [\#33](https://github.com/nifcloud/terraform-provider-nifcloud/pull/33) ([fuku2014](https://github.com/fuku2014))
- Feature/elb listener [\#31](https://github.com/nifcloud/terraform-provider-nifcloud/pull/31) ([fuku2014](https://github.com/fuku2014))
- Feature/router [\#29](https://github.com/nifcloud/terraform-provider-nifcloud/pull/29) ([aokumasan](https://github.com/aokumasan))
- Release/v1.1.0 [\#34](https://github.com/nifcloud/terraform-provider-nifcloud/pull/34) ([fuku2014](https://github.com/fuku2014))
- resource/nifcloud\_web\_proxy: Add resource [\#32](https://github.com/nifcloud/terraform-provider-nifcloud/pull/32) ([fuku2014](https://github.com/fuku2014))
- ci: Fix to execute acc test at PR [\#30](https://github.com/nifcloud/terraform-provider-nifcloud/pull/30) ([fuku2014](https://github.com/fuku2014))

## [v1.0.3](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.0.3) (2021-01-14)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.0.2...v1.0.3)

**Merged pull requests:**

- Fix/ssl certificate without ca [\#25](https://github.com/nifcloud/terraform-provider-nifcloud/pull/25) ([aokumasan](https://github.com/aokumasan))
- Fix a typo in nifcloud\_ssl\_certificate doc [\#24](https://github.com/nifcloud/terraform-provider-nifcloud/pull/24) ([aokumasan](https://github.com/aokumasan))

## [v1.0.2](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.0.2) (2020-12-14)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.0.1...v1.0.2)

**Merged pull requests:**

- resource/nifcloud\_security\_group\_rule: Fix bug of flatten [\#18](https://github.com/nifcloud/terraform-provider-nifcloud/pull/18) ([fuku2014](https://github.com/fuku2014))

## [v1.0.1](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.0.1) (2020-12-02)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/v1.0.0...v1.0.1)

**Merged pull requests:**

- resource/nifcloud\_security\_group\_rule: Fix to use mutexkv [\#2](https://github.com/nifcloud/terraform-provider-nifcloud/pull/2) ([fuku2014](https://github.com/fuku2014))
- fix nifcloud\_private\_lan descriptions [\#1](https://github.com/nifcloud/terraform-provider-nifcloud/pull/1) ([o108minmin](https://github.com/o108minmin))

## [v1.0.0](https://github.com/nifcloud/terraform-provider-nifcloud/tree/v1.0.0) (2020-11-30)

[Full Changelog](https://github.com/nifcloud/terraform-provider-nifcloud/compare/7827149ca2ee39c56f4493521329b1f0bb962e61...v1.0.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
