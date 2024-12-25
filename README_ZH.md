# 使用 Elasticsearch 和 HuggingFace Model Embeddings 实现语义搜索  

本项目结合 Elasticsearch 和 HuggingFace 模型生成的嵌入，实现高效的语义搜索。  

## 主要功能  
- **嵌入生成：** 轻松集成 HuggingFace 的开源模型，为文本数据生成高质量嵌入。  
- **Elasticsearch 支持：** 利用 Elasticsearch 强大的索引和查询功能，实现可扩展的语义搜索。  
- **模型自定义：** 可根据应用需求选择和配置嵌入模型。  
- **快速语义检索：** 基于相似度的搜索，提供相关且有意义的搜索结果，提升用户体验。  

该仓库非常适合用于构建文档搜索引擎、推荐系统和知识管理工具等应用。  

立即开始，实现适合您需求的前沿语义搜索！

## 依赖
- go ~> 1.22
- python ~> 3.12
- protoc ~> 25.3

## 使用方法
- 安装依赖：
```bash
make deps
```

- 启动语义搜索应用
```bash
# ./bin/searcher -e http://127.0.0.1:9200 -m ./local_models/paraphrase-multilingual-MiniLM-L12-v2
# -e 用于指定 Elasticsearch 地址。
# -m 用于指定模型。没有相对路径意味着从 HuggingFace 加载。
make run
```

- 模型微调
```bash
# python ./modelx/fine_tuning.py --dataset ./dataset/paraphrase-multilingual-minilm-l12-v2_dataset.csv --model ./output/local_models/paraphrase-multilingual-MiniLM-L12-v2 --version v1
# --dataset 指定微调数据集.
# --model 指定模型.
make ft
```

```text
我是一只快乐的程序猿

                          ___
                      .-'`     `'.
               __    /  .-. .-.   \
            .'`__`'.| /  ()|  ()\  \
           / /`   `\\ |_ .-.-. _|  ;  __
           ||     .-'`  (/`|`\) `-./'`__`'.
           \ \. .'                 `.`  `\ \
            `-./  _______            \    ||
               | |\      ''''---.__   |_./ /
               ' \ `'---..________/|  /.-'`
                `.`._            _/  /
                  `-._'-._____.-' _.`
                   _,-''.__...--'`
               _.-'_.    ,-. _ `'-._
            .-' ,-' /   /   \\`'-._ `'.
          <`  ,'   /   /     \\    / /
           `.  \  ;   ;       ;'  / /_
     __   (`\`. \ |   |       ||.' // )
  .'`_ `\(`'.`.\_\|   |    o  |/_,'/.' )
 / .' `; |`-._ ` /;    \     / \   _.-'
 | |  (_/  (_..-' _\    `'--' | `-.._)
 ; \        _.'_.' / /'.___.; \
  \ '-.__.-'_.'   ; '        \ \
   `-.,__.-'      | ;         ; '
                  | |         | |
                  | |         / /
                .-' '.      ,' `-._
              /`    _ `.   /  _    `.
             '-/ / / `\_) (_/` \  .`,)
              | || |            | | |
              `-'\_'            (_/-'
```                                 
