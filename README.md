# GO-SVC-TPL

go-svc-tpl 是为了快速构建 golang web 项目开发的项目框架，以及相关的一系列工具链。



## 框架使用

1. 未完成，请 clone 下来，在其基础上现改。
2. run `go generate ./api` in to update docs

## 设计思路

见文档 [design](./docs/design.md) (未更新)

Design-Drive 的后端开发模式，先进行 API 和 model 设计，后进行后端开发。

## 开发计划

`*` 为优先级较低的项目





- 使用 apifox 作为 api 文档编辑，导出 OpenAPI 文档并生成代码
- 寻找 go 中使用 `@` 做代码注解的实现





- [x] 完成一个最简示例版本
- [x] 解析 ts 生成代码的工具 (应该是使用node)
- [ ] 分析文档生成 api 的工具，或许还有模板代码的生成工具 (go开发)
- [ ] 支持 mongo





### ~~API-swag~~

- [ ] 支持 3.0
- [ ] 支持通过 model 中定义的 form/head tag 生成相应接口的参数
- [ ] model 解析格式向 json 规则对齐
- [ ] 支持更多校验规则
- [ ] 通过分析代码生成接口文档，而不是注释



### utils

封装一些常用功能，方便开发

或者更新一些作者已经明确标识不会维护的repo

- [ ] **fork stacktrace 并改进 current/code 等功能**
- [ ] *cache (generic interface)
