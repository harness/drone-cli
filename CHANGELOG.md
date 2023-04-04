# Changelog

## [1.7.0](https://github.com/harness/drone-cli/tree/1.7.0) (2023-01-24)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.6.2...1.7.0)

**Implemented enhancements:**

- Add support for jpath in jsonnet [\#224](https://github.com/harness/drone-cli/pull/224) ([rhiaxion](https://github.com/rhiaxion))

## [v1.6.2](https://github.com/harness/drone-cli/tree/v1.6.2) (2022-11-24)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.6.1...v1.6.2)

**Fixed bugs:**

- fix: use right parameter name for secrets-file [\#226](https://github.com/harness/drone-cli/pull/226) ([kameshsampath](https://github.com/kameshsampath))

**Merged pull requests:**

- \(maint\) prep v1.6.2 [\#227](https://github.com/harness/drone-cli/pull/227) ([tphoney](https://github.com/tphoney))

## [v1.6.1](https://github.com/harness/drone-cli/tree/v1.6.1) (2022-10-21)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.6.0...v1.6.1)

**Fixed bugs:**

- \(fix\) add secret file to compiler in exec [\#222](https://github.com/harness/drone-cli/pull/222) ([tphoney](https://github.com/tphoney))

**Merged pull requests:**

- \(maint\) release prep 1.6.1 & go tidy [\#223](https://github.com/harness/drone-cli/pull/223) ([tphoney](https://github.com/tphoney))

## [v1.6.0](https://github.com/harness/drone-cli/tree/v1.6.0) (2022-10-19)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.5.0...v1.6.0)

**Implemented enhancements:**

- \(drone-372\) use the modern docker runner for exec, upgrade deps [\#216](https://github.com/harness/drone-cli/pull/216) ([tphoney](https://github.com/tphoney))
- \(feat\) start migration from drone-yaml dependency. use docker compiler for lint [\#210](https://github.com/harness/drone-cli/pull/210) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- fix: use .drone.yml as default pipeline file [\#219](https://github.com/harness/drone-cli/pull/219) ([kameshsampath](https://github.com/kameshsampath))
- \(fix\): add labels for tooling to query containers [\#218](https://github.com/harness/drone-cli/pull/218) ([kameshsampath](https://github.com/kameshsampath))

**Merged pull requests:**

- \(maint\) v1.6.0 release prep [\#221](https://github.com/harness/drone-cli/pull/221) ([tphoney](https://github.com/tphoney))
- add community information [\#213](https://github.com/harness/drone-cli/pull/213) ([mrsantons](https://github.com/mrsantons))

## [v1.5.0](https://github.com/harness/drone-cli/tree/v1.5.0) (2022-01-04)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.4.0...v1.5.0)

**Implemented enhancements:**

- Update README.md [\#205](https://github.com/harness/drone-cli/pull/205) ([tphoney](https://github.com/tphoney))
- \(dron-124\) Add new command for build incomplete V2 [\#204](https://github.com/harness/drone-cli/pull/204) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- update libaries with vulnerabilies [\#208](https://github.com/harness/drone-cli/pull/208) ([eoinmcafee00](https://github.com/eoinmcafee00))
- fixes issue where template info command wasn't working [\#203](https://github.com/harness/drone-cli/pull/203) ([eoinmcafee00](https://github.com/eoinmcafee00))
- Update go-jsonnet to version v0.17.0 [\#202](https://github.com/harness/drone-cli/pull/202) ([hjkatz](https://github.com/hjkatz))

**Merged pull requests:**

- release prep v1.5 [\#211](https://github.com/harness/drone-cli/pull/211) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.4.0](https://github.com/harness/drone-cli/tree/v1.4.0) (2021-09-08)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.3.3...v1.4.0)

**Implemented enhancements:**

- provides ability to update auto-cancel-running flag on repo [\#198](https://github.com/harness/drone-cli/pull/198) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Fixed bugs:**

- \(fix\) fix multi-docs/lists in starlark [\#200](https://github.com/harness/drone-cli/pull/200) ([tphoney](https://github.com/tphoney))

**Merged pull requests:**

- \(maint\) release-v1.4.0 prep [\#201](https://github.com/harness/drone-cli/pull/201) ([tphoney](https://github.com/tphoney))

## [v1.3.3](https://github.com/harness/drone-cli/tree/v1.3.3) (2021-08-26)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.3.2...v1.3.3)

**Fixed bugs:**

- Fix `--stream` combined with `--format` for `jsonnet` [\#195](https://github.com/harness/drone-cli/pull/195) ([julienduchesne](https://github.com/julienduchesne))

**Merged pull requests:**

- \(maint\) release prep v1.3.3 [\#197](https://github.com/harness/drone-cli/pull/197) ([tphoney](https://github.com/tphoney))

## [v1.3.2](https://github.com/harness/drone-cli/tree/v1.3.2) (2021-08-25)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.3.1...v1.3.2)

**Fixed bugs:**

- \(fix\) re-enable format option for jsonnet command [\#193](https://github.com/harness/drone-cli/pull/193) ([tphoney](https://github.com/tphoney))

**Merged pull requests:**

- \(maint\) release prep for v1.3.2 [\#194](https://github.com/harness/drone-cli/pull/194) ([tphoney](https://github.com/tphoney))

## [v1.3.1](https://github.com/harness/drone-cli/tree/v1.3.1) (2021-08-20)

[Full Changelog](https://github.com/harness/drone-cli/compare/v1.3.0...v1.3.1)

**Fixed bugs:**

- Defect/permission template [\#191](https://github.com/harness/drone-cli/pull/191) ([eoinmcafee00](https://github.com/eoinmcafee00))
- \(DRON-113\) use ghodss/yaml for yaml printing [\#189](https://github.com/harness/drone-cli/pull/189) ([tphoney](https://github.com/tphoney))
- fixes issue were cli required an additional parameter in order to com… [\#188](https://github.com/harness/drone-cli/pull/188) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Merged pull requests:**

- release prep v1.3.1 [\#192](https://github.com/harness/drone-cli/pull/192) ([eoinmcafee00](https://github.com/eoinmcafee00))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
