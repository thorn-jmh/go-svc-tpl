# GO-SVC-TPL

go-svc-tpl 是为了快速构建 golang web 项目开发的项目框架，以及相关的一系列工具链。



## 框架使用

1. 未完成，请 clone 下来，在其基础上现改。
2. run `go generate ./api` in to update docs

## 设计思路

见文档 [design](./docs/design.md)

目前设想有两种开发模式选择：
- 其一是预先定义好 API 文档，从而生成相应的代码，只需填写必要的逻辑部分，这应该能适应大部分项目的合理开发流程。
- 另一种是先写代码，再通过代码生成 API 文档，这种方式不可避免要写很多重复，低逻辑的代码，这是 web 后端本身的特性所决定的，但是这种方式可能可以适应一些前后端开发周期不一致的项目，或是一些以 db 为中心的项目。


## 开发计划

`*` 为优先级较低的项目



#### 2023.07.26

- 决定使用 ts 作为 API/model 描述语言，废弃使用 openAPI
- 寻找 go 中使用 `@` 做代码注解的实现





- [x] 完成一个最简示例版本
- [ ] 解析 ts 生成代码的工具 (应该是使用node)
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
- [ ] fork [guregu/null](https://github.com/guregu/null/tree/master) 以确保在 gorm 中没有预期外行为，

