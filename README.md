# go-mods

来源于日常开发中积累的go代码、抽出作为模块、方便以后开发使用

1. 优先尽可能引用 golang 标准库
2. 基础功能工具
   - 可引用 [Goframe](https://github.com/gogf/gf) 框架内模块化包
   - 或者 [lo](https://github.com/samber/lo) 泛型支持工具包
3. 最后是使用三方库

### TODO

后续需整理下组件类型
有些组件需要实例存在的，比如captcha、qrcode等（类似作为goframe中的service层）
有些组件类似于工具函数、不需要实例存在的，比如utils

http-req组件目前没太好想法、目前既有工具函数、也有实例方法
