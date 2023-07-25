# GO-SVC-TPL

go-svc-tpl 是为了快速构建 golang web 项目开发的项目框架，以及相关的一系列工具链。



## 框架使用

1. 未完成，请 clone 下来，在其基础上现改。
2. run `go generate ./apii` in `src` to update openapi docs

## 设计思路

见文档 [design](./docs/design.md)



## 开发计划

`*` 为优先级较低的项目

加粗 为较为优先的项目



- [ ] 完成一个最简示例版本
- [ ] 添加 generate 脚本/cmd



### API-swag

- [ ] 支持 3.0
- [ ] 支持通过 model 中定义的 form/head tag 生成相应接口的参数
- [ ] model 解析格式向 json 规则对齐
- [ ] 支持更多校验规则
- [ ] 通过分析代码生成接口文档，而不是注释



### utils

- [ ] **fork stacktrace 并改进 current/code 等功能**
- [ ] *cache (generic interface)
- [ ] fork [guregu/null](https://github.com/guregu/null/tree/master) 并确认在 gorm 中没有预期外行为，

